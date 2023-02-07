# Writing Web Applications
[https://go.dev/doc/articles/wiki/ の日本語訳 + 問題

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

`GOPATH` の中にこのチュートリアル用の新しいディレクトリを作り、そこに `cd` します。
```
$ mkdir gowiki
$ cd gowiki
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

Go で書かれた非常にシンプルな web サーバの一例です:
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