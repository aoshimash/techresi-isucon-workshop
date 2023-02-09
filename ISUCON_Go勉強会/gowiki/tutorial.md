# Writing Web Applications
https://go.dev/doc/articles/wiki/ の日本語訳 + 問題

## Introduction

このチュートリアルで扱う内容:

- ロードとセーブのメソッドを持つデータ構造の作成
- `net/http` パッケージを使用した Web アプリケーションの構築
- `html/template` パッケージを使用した HTML テンプレートの処理
- `regexp` パッケージを使用したユーザー入力の検証
クロージャの使用

想定される知識:
- プログラミングの経験
- 基本的なWeb技術（HTTP、HTML）の理解
- UNIX/DOSのコマンドラインに関する知識

## Getting started

現在、Goを実行するにはFreeBSD、Linux、macOS、Windowsのいずれかのマシンが必要です。ここでは、コマンドプロンプトを表すために `$` を使用します。

Go をインストールします（[インストール方法](https://go.dev/doc/install)を参照してください）。

環境変数 `GOPATH` に現在のディレクトリ `gowiki` を追加します (docker の場合は不要):
```
$ GOPATH=${GOPATH}:${PWD}
```

`wiki.go` という名前のファイルを作成し、好きなエディタで開いて、以下の行を追加します。
```
package main

import (
    "fmt"
    "os"
)
```

Go 標準ライブラリから `fmt` と `os` のパッケージをインポートしています。その後、追加機能を実装する際には、この `import` 宣言にさらにパッケージを追加していくことになります。

### Go basics
`go_basics.go` を開き中身を見てください。実行するには以下のコマンドを実行します:
```
$ docker run -v ${PWD}:/go -it --rm golang:1.19
# go build go_basics.go
# ./go_basics
```
ローカルに Go の実行環境があるなら、以下をそのまま実行します:
```
$ go build go_basics.go
$ ./go_basics
```

次のような出力が見れるはずです:
```
5 3.14 I'm a perfect human
0
1
2
3
4
5
6
7
8
9
I'm not a human
3
I'm Shiba, named Taro and 7 years old.
Taro is still young!
42
21
73
[2 3 5 7 11 13]
index 0, value 2
index 1, value 3
index 2, value 5
index 3, value 7
index 4, value 11
index 5, value 13
```

## Data structures

データ構造を定義することから始めましょう。wiki は相互に接続された一連のページで構成され、それぞれのページは Title と Body（ページのコンテンツ）を持っています。ここでは、Title と Body を表す2つのフィールドを持つ構造体として、`Page` を定義します。

```
type Page struct {
    Title string
    Body  []byte
}
```
型 `[]byte` は「byte の slice」を意味します。(`Body` 要素が文字列ではなく `[]byte` であるのは、後述するように、使用する `io` ライブラリが期待する型であるためです。

`Page` 構造体は、ページデータがメモリにどのように格納されるかを記述します。しかし、永続的な保存についてはどうでしょうか。それは、`Page` に `save` メソッドを作成することで対処できます。

```
func (p *Page) save() error {
    filename := p.Title + ".txt"
    return os.WriteFile(filename, p.Body, 0600)
}
```
このメソッドは、`Page` 構造体のメソッドとなり、`p.save()` と書くことで呼び出すことができます。返り値は `error` です。
これは `Page` の `Body` をテキストファイルに保存します。簡単のために、`Title` をファイル名として使用します。

`save` メソッドがエラー値を返すのは、それが `WriteFile`（バイトスライスをファイルに書き込む標準ライブラリ関数）の戻り値の型だからです。`save` メソッドはエラー値を返しますが、これはファイル書き込み中に何か問題が発生した場合に、アプリケーション側で対処できるようにするためです。うまくいけば、`Page.save()` は `nil`（他言語の null）を返します。

`WriteFile` の3番目のパラメータとして渡される8進数の整数リテラル`0600` は、ファイルが現在のユーザーのみの読み書き権限で作成されるべきであることを示します。(詳細については、Unix の man ページ `open(2)` を参照してください)。

ページを保存するだけでなく、ページをロードすることも必要になります。
```
func loadPage(title string) *Page {
    filename := title + ".txt"
    body, _ := os.ReadFile(filename)
    return &Page{Title: title, Body: body}
}
```
関数 `loadPage` は、`title` パラメータからファイル名を構築し、ファイルの内容を新しい変数 `body` に読み込み、適切な `title` と`body` の値で構築された `Page` リテラルへのポインタを返します。

関数は複数の値を返すことができます。標準ライブラリ関数である `os.ReadFile` は `[]byte` と `error` を返します。`loadPage` では、`error` はまだ処理されていません。アンダースコア（_）記号で表される「空白の識別子」を使って、エラーの戻り値を捨てています（要するに、その値を何も代入していない）。

しかし、`ReadFile` がエラーに遭遇した場合はどうなるのでしょうか？たとえば、ファイルが存在しないかもしれない。このようなエラーを無視してはいけません。この関数を修正して、`*Page` と `error` を返すようにしましょう。
```
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}
```

この変更によって、この関数は2番目のパラメータをチェックすることができます。もしそれがnilであれば、ページの読み込みに成功したことになります。もしそうでなければ，呼び出し側が処理できるエラーになります(詳細は言語仕様書を参照してください)。

この時点で，簡単なデータ構造と，ファイルへの保存とファイルからの読み込みができるようになりました．では、書いたものをテストするために`main` 関数を書いてみましょう。

```
func main() {
    p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
    p1.save()
    p2, _ := loadPage("TestPage")
    fmt.Println(string(p2.Body))
}
```

このコードをコンパイルして実行すると、`p1` の内容を含む`TestPage.txt` という名前のファイルが作成されます。このファイルを構造体 `p2` に読み込んで、その `Body` を画面に表示します。

このプログラムをコンパイルして実行するには、次のようにします。
```
$ go build wiki.go
$ ./wiki
```

## Introducing the net/http package (an interlude)

以下は Go で書かれた非常にシンプルな web サーバの一例です:
```
//go:build ignore

package main

import (
    "fmt"
    "log"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```
`main` 関数は、`http.HandleFunc` の呼び出しから始まり、Web ルート（"/"）へのすべてのリクエストをハンドラで処理するよう `http` パッケージに指示します。

次に `http.ListenAndServe` を呼び出し、任意のインターフェイス (":8080") 上のポート 8080 を listen するよう指定します。(この関数は、プログラムが終了するまでブロックされます。

`ListenAndServe` は常にエラーを返します。なぜなら、予期せぬエラーが発生したときのみ、エラーを返すからです。そのエラーを記録するために、関数呼び出しを `log.Fatal.Handler` でラップします。

`handler` メソッドは `http.HandlerFunc` 型です。これは、`http.ResponseWriter` と `http.Request` を引数として受け取ります。

`http.ResponseWriter` は、HTTPサーバーのレスポンスをまとめた値で、これに書き込むことでHTTPクライアントにデータを送信します。

`http.Request` は、クライアントのHTTPリクエストを表すデータ構造です。`r.URL.Path` は、リクエストURLのパスコンポーネントです。末尾の[1:]は、"1文字目から末尾までのPathのサブスライスを作成する "ことを意味します。これにより、パス名から先頭の"/"が削除されます。

このプログラムを実行し、以下の URL へアクセスしてみましょう:
```
http://localhost:8080/monkeys
```

次のようなページが表示されるはずです。
```
Hi there, I love monkeys!
```

## Using net/http to serve wiki pages
`net/http` パッケージを使うためには、まずインポートする必要があります:
```
import (
    "fmt"
    "os"
    "log"
    "net/http"
)
```

ユーザが wiki ページを見ることができるように `viewHandler` を作り、 "/view/" から始まる URL に紐づけてみましょう:
```
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}
```

再度アンダースコア (`_`) が登場していることに注意してください。ここでは `loadPage` メソッドのエラーが無視されています。簡単のために今はこう実装されていますが、一般にこれはバッドプラクティスです。これについては、後で説明します。

まず、この関数はリクエスト URL のパスコンポーネントである `r.URL.Path` からページタイトルを抽出します。Pathは`[len("/view/"):]`で再スライスされ、リクエストパスの先頭の "/view/" コンポーネントが削除されます。これは、パスが必ず「/view/」で始まるためで、これはページのタイトルの一部ではありません。

この関数は、ページのデータを読み込み、シンプルなHTMLの文字列でページをフォーマットし、`http.ResponseWriter` の `w` に書き出します。

このハンドラを使用するために、メイン関数を書き換えて、`viewHandler` を使用して http を初期化し、パス /view/ 以下のすべてのリクエストを処理するようにします。

```
func main() {
    http.HandleFunc("/view/", viewHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

いくつかのページデータ f(test.txt) を作成し、コードをコンパイルして、wiki ページを配信してみましょう。

test.txt をエディタで開き、"Hello world" という文字列 (引用符なし) をその中に保存します。

$ go build wiki.go
$ ./wiki
(Windowsを使用している場合、プログラムを実行するには、"./"なしで "wiki "と入力する必要があります)

このウェブサーバーが動いている状態で、http://localhost:8080/view/test にアクセスすると、"Hello world" という単語を含む "test" というタイトルのページが表示されるはずです。

## Editing pages
編集機能なしでは wiki とは言えません。次の新しい2つのハンドラを作ってみましょう:
- `editHandler`: edit フォームを提供するハンドラ
- `saveHandler`: フォームから入力されたデータを保存するハンドラ

まず、これら2つのハンドラを `main()` に追加します:
```
func main() {
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```
関数 `editHandler` は、ページをロードし（存在しない場合は空の`Page` 構造体を作成する）、HTMLフォームを表示するようにします:

```
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    fmt.Fprintf(w, "<h1>Editing %s</h1>"+
        "<form action=\"/save/%s\" method=\"POST\">"+
        "<textarea name=\"body\">%s</textarea><br>"+
        "<input type=\"submit\" value=\"Save\">"+
        "</form>",
        p.Title, p.Title, p.Body)
}
```
この機能はうまく動作しますが、ハードコードされたHTMLは醜いものです。もちろん、もっと良い方法があります。

## The html/template package
`html/template` パッケージは Go の標準ライブラリです。このライブラリを使うことによって、HTML を別のファイルに置くことができ、Go のコードを変えずに edit ページのレイアウトを変えることができます。

まず最初に、`html/template` を `import` のリストに入れることが必要です。また、もう `fmt` を使わないので同時に削除しておきましょう:

```
import (
    "html/template"
    "os"
    "net/http"
)
```

次に、HTML フォームを含むテンプレートを作りましょう。`edit.html` という名前のファイルを作り、以下の内容を追加します:

```
<h1>Editing {{.Title}}</h1>

<form action="/save/{{.Title}}" method="POST">
<div><textarea name="body" rows="20" cols="80">{{printf "%s" .Body}}</textarea></div>
<div><input type="submit" value="Save"></div>
</form>
```

テンプレートを使うように `editHandler` も変更します:

```
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    t, _ := template.ParseFiles("edit.html")
    t.Execute(w, p)
}
```
`template.ParseFiles` 関数は `edit.html` の内容を読み取り `*template.Template` を返します。

`t.Execute` はテンプレートを実行し、生成された HTML を `http.ResponseWritier` へ書き込みます。テンプレート中の `.Title` と `.Body` は `p.Title` と `p.Body` を参照します。

テンプレートの Go によって値が埋められる部分は、二重中括弧で括られます。 `printf "%s" .Body` は `fmt.Printf` と同様に、 `.Body` を `[]byte` ではなく `string` 型として出力する指示になります。`html/template` パッケージは安全かつ正しく見える HTML のみが生成されるように保証されています。例えば、`>` マークは常に `&gt;` に置き換えられ、ユーザの入力が form の HTML を破壊しないようになっています。

`viewHandler` のためのテンプレート `view.html` も作ってみましょう:
```
<h1>{{.Title}}</h1>

<p>[<a href="/edit/{{.Title}}">edit</a>]</p>

<div>{{printf "%s" .Body}}</div>
```

`viewHandler` も合わせて書き換えると、
```
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    t, _ := template.ParseFiles("view.html")
    t.Execute(w, p)
}
```

`viewHandler` も `editHandler` もテンプレートを実行する部分の処理は共通なので新しくメソッドにすると、
```
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, _ := template.ParseFiles(tmpl + ".html")
    t.Execute(w, p)
}
```

この関数を用いてハンドラを書き換えると、
```
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    renderTemplate(w, "view", p)
}
```

```
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}
```

`main` で未実装の保存ハンドラの登録をコメントアウトすれば、再びプログラムをビルドしてテストすることができます:
```
func main() {
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    // http.HandleFunc("/save/", saveHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Handling non-existent pages

もしあなたが `/view/APageThatDoesntExist` にアクセスしたらどうなるでしょうか？ HTML を含むページが表示されるでしょう。これは`loadPage` からのエラーの戻り値を無視し、データなしでテンプレートを埋めようとし続けるからです。代わりに、要求されたページが存在しない場合、コンテンツが作成されるようにクライアントを編集ページにリダイレクトする必要があります:

```
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}
```

`http.Redirect` 関数は、HTTP レスポンスに HTTP ステータスコードとして `http.StatusFound（302）`、`Location` ヘッダーを追加する関数です。

## Saving pages

関数 `saveHandler` は、編集ページに配置されたフォームの送信を処理します。`main` の関連行をアンコメントした後、ハンドラを実装してみましょう。

```
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    p.save()
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
```

URL で提供されるページのタイトルとフォームの唯一のフィールドである`Body` は、新しい `Page` に格納されます。その後、`save()` メソッドが呼び出されてデータがファイルに書き込まれ、クライアントは `/view/` ページにリダイレクトされます。

`FormValue` が返す値は `string` 型です。`Page` 構造体に収めるには、この値を `[]byte` に変換する必要があります。ここでは `[]byte(body)` を使って変換しています。

## Error handling

私たちのプログラムには、エラーを無視している箇所がいくつかあります。これは悪い習慣で、特にエラーが発生したときにプログラムが意図しない動作をすることになるからです。より良い解決策は、エラーを処理し、ユーザーにエラーメッセージを返すことです。そうすれば、何か問題が発生したときに、サーバーは私たちが望むように正確に機能し、ユーザーにはそのことが通知されます。

まず、renderTemplateでエラーを処理しましょう:

```
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, err := template.ParseFiles(tmpl + ".html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = t.Execute(w, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

`http.Error` 関数は、指定した `HTTP` レスポンスコード（ここでは "Internal Server Error"）とエラーメッセージを送信します。すでに、この関数を別の関数にしたことが功を奏しています。

では、`saveHandler` を修正しましょう:

```
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
```

`p.save()` 中に発生したエラーは、ユーザーに報告されます。

## (Optional) Template caching

このコードには非効率性があります。`renderTemplate` はページがレンダリングされるたびに `ParseFiles` を呼び出します。よりよい方法は、プログラムの初期化時に一度だけ `ParseFiles` を呼び出し、すべてのテンプレートを単一の `*Template` にパースすることです。それから、`ExecuteTemplate` メソッドを使って特定のテンプレートをレンダリングすることができます。

まず、`templates` というグローバル変数を作成し、それを `ParseFiles` で初期化します:

```
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
```

関数` template.Must` は、`nil` でないエラー値が渡されると panic を起こし、そうでない場合は `*Template` をそのまま返す便利なラッパーです。テンプレートを読み込むことができない場合、プログラムを終了することが唯一の賢明な方法です。

`ParseFiles` 関数はテンプレート・ファイルを識別する任意の数の文字列引数を取り、それらのファイルを基本ファイル名の後に名付けられたテンプレートにパースします。プログラムにさらにテンプレートを追加する場合は、その名前を `ParseFiles` 呼び出しの引数に追加します。

次に、適切なテンプレートの名前で `templates.ExecuteTemplate` メソッドを呼び出すように `renderTemplate` 関数を変更します:

```
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

なお、テンプレート名はテンプレートファイル名なので、`tmpl` の引数に `".html"` を追加する必要があります。

## (Optional) Validation

お気づきかもしれませんが、このプログラムには重大なセキュリティ上の欠陥があります。 ユーザがサーバー上で読み書きをするために、任意のパスを供給することができるのです。これを軽減するために、タイトルを正規表現で検証する関数を書くことができます。

まず、インポートリストに `"regexp"` を追加し、グローバル変数を作成してそこに検証式を格納します:

```
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
```

関数 `regexp.MustCompile` は正規表現を解析してコンパイルし、`regexp.Reexp` を返します。`MustCompile` は `Compile` と異なり、コンパイルに失敗すると panic を引き起こしますが、`Compile` は第2引数としてエラーを返します。

では、validPath式を使って、パスを検証し、ページのタイトルを取り出す関数を書いてみましょう:

```
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("invalid Page Title")
    }
    return m[2], nil // The title is the second subexpression.
}
```

タイトルが有効な場合は、`nil` のエラー値とともに返されます。タイトルが無効な場合は、HTTP 接続に "404 Not Found" エラーを書き込み、ハンドラにエラーを返します。新しいエラーを作成するには、`errors` パッケージをインポートする必要があります。

各ハンドラーに getTitle の呼び出しを記述してみましょう:

```
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}
```

```
func editHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}
```

```
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err = p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
```

## (Optional) Introducing Function Literals and Closures

各ハンドラでエラーをキャッチすると、多くのコードが繰り返されることになります。各ハンドラを、この検証やエラーチェックを行う関数でラップできたらどうでしょう？ Go の関数リテラルは、機能を抽象化する強力な手段であり、ここで役に立ちます。

まず、各ハンドラの定義をタイトル文字列を受け入れるように書き直します:

```
func viewHandler(w http.ResponseWriter, r *http.Request, title string)
func editHandler(w http.ResponseWriter, r *http.Request, title string)
func saveHandler(w http.ResponseWriter, r *http.Request, title string)
```

次に、上記の型の*関数*を受け取り、`http.HandlerFunc` 型の関数を返すラッパー関数を定義してみましょう:

```
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Here we will extract the page title from the Request,
        // and call the provided handler 'fn'
    }
}
```

返された関数は、その外部で定義された値を囲むため、クロージャと呼ばれます。この場合、変数 `fn `(makeHandler の単一の引数) はクロージャで囲まれています。変数 `fn` は、保存、編集、または表示ハンドラのいずれかになります。

さて、`getTitle` からコードを取り出し、若干の修正を加えることでここで使うことができます:

```
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}
```

`makeHandler` が返すクロージャは `http.ResponseWriter` と `http.Request` を受け取る関数です (言い換えれば `http.HandlerFunc` です)。クロージャは、リクエストパスからタイトルを抽出し、`validPath` 正規表現でそれを検証します。タイトルが無効な場合、`http.NotFound` 関数を使用して、`ResponseWriter` にエラーが書き込まれます。タイトルが有効な場合は、`ResponseWriter`、`Request`、`title` を引数として同封のハンドラ関数 `fn` が呼び出されます。

これで、ハンドラ関数が `http` パッケージに登録される前に、`main` で `makeHandler` でラップすることができます:

```
func main() {
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

最後に、ハンドラから getTitle の呼び出しを削除しよりシンプルにしましょう:

```
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}
```

```
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}
```

```
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
```

## Try it out!
コードを再コンパイルし、アプリを実行してみましょう:

```
$ go build wiki.go
$ ./wiki
```

http://localhost:8080/view/ANewPage にアクセスすると、ページ編集フォームが表示されます。テキストを入力し、"Save" をクリックすると、新しく作成されたページにリダイレクトされるはずです。