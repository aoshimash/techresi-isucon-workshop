# techresi-isucon-workshop/ISUCONインフラ勉強会

## 今日の内容

ISUCON12の予選で実際に行った、**サーバーの初期セットアップ方法**と**デプロイスクリプトの作成方法** の説明と、実際にスクリプトを開発するときのTipsの紹介。

## ISUCON12当日のインフラ作業のおさらい

運営から環境が配られてから実際にコードの修正などの作業に入るまでには、色々とインフラ周りの準備が必要。
まずは、ISUCON12予選で実際に行った作業の全体の流れをおさらいする。

1. 運営から配布されるCloudFormationを使ってAWSコンソールから環境を構築
2. サービスの動作確認
   1. ブラウザからサービスを使うことができるか確認
   2. ベンチマークの実行
3. アプリの言語をPython(来年はGo)に切り替えて 2 を再度実施
4. git管理したいファイルをサーバーから抽出してgithubにアップする
  - `/home`と`/etc`を全部git管理しとけばいいと思う
    - ISUCON12ではDBのダンプファイルみたいな大きなファイルをgitに入れたくないから細かくファイル指定して必要なファイルだけコピーしていたが、必要なファイルの選別はめんどくさいから全部コピーしてくればいいと思う。
  - また、git管理にしたファイルは今後書き換えるファイルなのでちゃんとバックアップを取っておく
  - つまりこういうこと
    ```bash
    // サーバー側
    sudo cp -r /home /home.bak
    sudo cp -r /etc /etc.bak

    // ローカルPC側
    rsync -avz ubuntu@<サーバーIP>:/home $GITREPO_DIR/contents
    rsync -avz ubuntu@<サーバーIP>:/etc $GITREPO_DIR/contents
    ```
5. ansibleのinventoryファイルの更新
  - [`playbooks/hosts/isucon12q.yaml`](https://github.com/aoshimash/isucon12q-techresi/blob/main/playbooks/hosts/isucon12q.yaml)（ansibleの操作対象マシンの設定）を書き換える
6. 初期設定スクリプトを流す
  - [`playbook/setup.yaml`](https://github.com/aoshimash/isucon12q-techresi/blob/main/playbooks/setup.yaml)
  - このスクリプトでは次の操作をしている
    - ユーザー作成とSSH鍵の配置
    - 各種パッケージのインストール(alp, emacs, htop, ...)
    - newlelicを使うならここでagentのセットアップをする
  - このスクリプトは事前準備ができるので、ISUCON本番前には用意しておく
7. デプロイスクリプトの作成
  - [`playbooks/prepare_bench.yaml`](https://github.com/aoshimash/isucon12q-techresi/blob/main/playbooks/prepare_bench.yaml)
  - このスクリプトでは次の操作をしている
    - git管理しているファイルを操作対象マシンに配置
    - リスタートが必要なミドルウェアがある場合はリスタートさせる
    - ログファイルのローテーション
  - このスクリプトはある程度の事前準備はできるが、ファイルの配置などは実際のサーバーを見ないとわからないので、当日サーバーを見て修正する。

あとは、コードの修正修正 → デプロイ → ベンチマーク → 計測 の繰り返し。

## リポジトリの解説

ISUCONの試合で使うファイルは全部１つのリポジトリで管理している。

- ISUCON12予選: https://github.com/aoshimash/isucon12q-techresi
- 改良版: https://github.com/Tech-Residence/isucon-template

内容はほとんど同じだが、今日は改良版の方をベースに説明する

### リポジトリ構成

```
.
├── Dockerfile: ansibleを実行するためのDockerfile
├── LICENSE
├── Makefile
├── README.md
├── contents: サーバーからコピーしてきたファイル
│   ├── etc
│   └── home
├── docker-compose.yaml: Dockerfileをラップしている
├── mock
│   └── cloud-init.yaml: playbook開発用のモックサーバー
└── playbooks
    ├── ansible.cfg: ansibleの設定ファイル
    ├── inventory.yaml: インベントリファイル(操作対象マシンの情報)
    ├── prepare_bench.yaml: ベンチマークテスト実行前の準備を行うplaybook
    └── setup.yaml: サーバーの初期セットアップを行うplaybook
```

- ansibleの実行はdockerで行う
  - pythonのバージョン・ansible本体のバージョン・ansible moduleのバージョンが実行環境ごとにずれないようにコンテナ化している。
- dockerコマンドを直接叩くのではなくdocker-composeでラップしてからさらにmakeでラップしてる
  - 直接dockerコマンドを使おうとするとオプションがめちゃくちゃ多いから

### playbook詳解

`playbooks`ディレクトリの中身についてさらに深堀りしていく。

#### ansible.cfg

- https://github.com/Tech-Residence/isucon-template/blob/main/playbooks/ansible.cfg

これはansibleの設定ファイル。ansibleの挙動を設定するが、最初に一回設定したらその後さわることはあまりない。

ちなみに今は `host_key_checking = False` だけ設定している。
https://docs.ansible.com/ansible/latest/inventory_guide/connection_details.html#managing-host-key-checking

これは、ansibleがサーバーにSSH接続して操作する際のhost_keyチェックを無視させる設定。これをやってないとplaybook初回実行前に（known_hostsに鍵がないから）エラーがでる。

#### inventory

- [inventory.yaml](https://github.com/Tech-Residence/isucon-template/blob/main/playbooks/inventory.yaml)

操作対象のサーバー情報を記載している。実行時にどのinventoryファイルのどのホストに対してplaybookを実行するのか指定することで、操作対象のサーバーを切り替えることができる。

`all.hosts` 以下に操作対象のホスト情報一覧を記載している。

**inventoryファイルの書き方**
- https://docs.ansible.com/ansible/latest/inventory_guide/intro_inventory.html

####　playbookファイル

- [setup.yaml](https://github.com/Tech-Residence/isucon-template/blob/main/playbooks/setup.yaml)
- [prepare_bench.yaml](https://github.com/Tech-Residence/isucon-template/blob/main/playbooks/prepare_bench.yaml)

サーバーの設定が書いてあるファイル。たぶん一番編集するファイル。

各ディレクティブの説明

- name: playbookの名前
- hosts: playbookの操作対象となるhost名 or グループ名 (今回は全部のhostをallグループに所属させているので、allを指定すると全hostが操作対象になる。実行時オプションで操作対象を限定することもできる)
- become: サーバー内で操作するときに権限昇格システム(linuxサーバーであれば`sudo`など)を使ってroot権限でタスクを実行する　https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_privilege_escalation.html
- `vars`: 変数一覧。変数を別ファイルなどに切り出すこともできるが、この規模のplaybookの場合はファイルを分けると逆に見づらくなるので、変数は全部ここに書くことにしたい（playbookが大きくなってきた場合は方針変更あるかも）。
- `tasks`: `task`一覧。`module`と呼ばれる各種操作のユニットを並べていく。`module`は　https://docs.ansible.com/ansible/latest/collections/index.html#list-of-collections で検索できる。linuxの設定をするときに使いたいモジュールはほとんど[`ansible.builtin`コレクション](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/index.html#plugins-in-ansible-builtin)に入っている。

**補足**
- collectionとはplaybook, roll, module, plugin を含むことができるansibleコンテンツのディストリビューション形式。プログラミング言語で言うライブラリみたいなもの。各種OS・クラウド・ネットワーク機器向けのものが公開されているため、使いたいcollectionをインストールして使う。`ansible.builtin`コレクションは名前の通りansible本体に同梱されているので追加でインストールする必要はない。
ちなみに、ISUCONでは `ansible.builtin`コレクション以外を使うことはないと思うが、他のコレクションは [ansible galaxy](https://galaxy.ansible.com/)で探すことができる。
- [ISUCON12ではrollを使っていたが](https://github.com/aoshimash/isucon12q-techresi/tree/main/playbooks/roles)、この規模のplaybookで再利用性とか気にする必要がないので使うのをやめた。


## 参考資料

- [公式ドキュメント](https://docs.ansible.com/ansible/latest/getting_started/index.html)
  - ansibleについてわからないことがあればまずはここをみる。
- [VSCode拡張](https://marketplace.visualstudio.com/items?itemName=redhat.ansible)
  - とりあえずこれだけ入れておけばいいと思う。
