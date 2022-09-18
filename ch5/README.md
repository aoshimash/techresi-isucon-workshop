# ISUCON本ゆる輪読会#5 [Chapter5 データベースのチューニング]

これは[ISUCON本ゆる輪読会#5](https://connpass.com/event/260243/)用の資料です。

[ISUCON本ゆる輪読会#4](https://tech-residence.connpass.com/event/258772/)のつづきになるので、ハンズオンに参加する場合は[前回の資料](https://github.com/aoshimash/techresi-isucon-workshop/blob/main/ch5/slides/slides.pdf)の通りmultipassのハンズオン環境を作成しておいてください。

## 演習環境の用意

次のコマンドで前回のVMを立ち上げ直す、

```
multipass start private-isu
```

### 操作の復習

- VMのIP確認

```
multipass list
```

- VMへのログイン

```
multipass shell private-isu
```

- ベンチマーク実行

VM内の`/home/isucon/private_isu.git/benchmarker/`で以下を実行

```
./bin/benchmarker -u ./userdata -t http://localhost
```

- MySQLへの接続

VM内で次のコマンド

```
sudo mysql -uroot -proot
```

## 前回のおさらい

pt-query-digestで遅いクエリを見つけ出し、インデックスをはることでスコアを改善することができた


## 今日やること

- インデックスの理解と活用
- N+1問題の発見と解消

## インデックスの理解と活用

### インデックスで検索が高速になる理由（ゆるふわ）

インデックスは辞典における「索引」みたいなもの。索引がない辞典を想像すると、目当ての項目を見つけるのにものすごく時間がかかりそう。でも、特定のルールに沿って並んでいる索引（あいうえお順など）があると、検索速度が飛躍的に上がる。そういうかんじのやつ。

なにも手がかりがない状態で探索するとなると線形探索（$\mathcal{O}(n)$）になっちゃうけど、データ群が特定のルールに沿って並んでいる（インデックスがはってある）なら、二分探索（ $\mathcal{O}(\log n)$）でいけるからそれだけですごく速くなりそうだと想像できる。

（実際にデータベースのインデックスも二分探索に適したBツリーがよく利用されているらしい。）

### インデックスがはられているか確認

前回は`comments`テーブルの`post_id`にインデックスを貼って検索が高速になるところまで確認したが、そういえばインデックスがちゃんとはられているのか確認していなかったので、まずはそれを確認。

MySQLにログインして、データベースを選択してから、

```sql
USE isuconp;
```

次のコマンドでテーブルの構造を確認することができる。

```sql
DESC comments;
```
結果
```
+------------+-----------+------+-----+-------------------+-------------------+
| Field      | Type      | Null | Key | Default           | Extra             |
+------------+-----------+------+-----+-------------------+-------------------+
| id         | int       | NO   | PRI | NULL              | auto_increment    |
| post_id    | int       | NO   | MUL | NULL              |                   |
| user_id    | int       | NO   |     | NULL              |                   |
| comment    | text      | NO   |     | NULL              |                   |
| created_at | timestamp | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
+------------+-----------+------+-----+-------------------+-------------------+
```

indexが設定されているかどうかは `Key` 列でわかる。`PRI`がPrimary Key で `MUL`がMultiple Key（重複可能なキー）。


ちなみに、インデックスを確認したいだけなら次のコマンドでもOK。

```sql
SHOW INDEX FROM comments;
```
結果
```
+----------+------------+-------------+--------------+-------------+-----------+-------------+----------+--------+------+------------+---------+---------------+---------+------------+
| Table    | Non_unique | Key_name    | Seq_in_index | Column_name | Collation | Cardinality | Sub_part | Packed | Null | Index_type | Comment | Index_comment | Visible | Expression |
+----------+------------+-------------+--------------+-------------+-----------+-------------+----------+--------+------+------------+---------+---------------+---------+------------+
| comments |          0 | PRIMARY     |            1 | id          | A         |       99505 |     NULL |   NULL |      | BTREE      |         |               | YES     | NULL       |
| comments |          1 | post_id_idx |            1 | post_id     | A         |       10075 |     NULL |   NULL |      | BTREE      |         |               | YES     | NULL       |
+----------+------------+-------------+--------------+-------------+-----------+-------------+----------+--------+------+------------+---------+---------------+---------+------------+
2 rows in set (0.01 sec)
```


### クエリの実行計画を確認

どのインデックスを使うのかとかは、オプティマイザというDBの機能がSQLクエリを解析してなんかいい感じにやってくれてる。（オプティマイザが何をしているのかより詳細に知りたい場合はこちら。 https://dev.mysql.com/doc/internals/en/optimizer-tracing.html）

`EXPLAIN`ステートメントでクエリの実行計画を確認して、意図したとおりにインデックスが使われているか確認することができる。

``` sql
EXPLAIN SELECT * FROM `comments` WHERE `post_id` = 100 ORDER BY `created_at` DESC LIMIT 3\G
```
結果
```
*************************** 1. row ***************************
           id: 1
  select_type: SIMPLE
        table: comments
   partitions: NULL
         type: ref
possible_keys: post_id_idx
          key: post_id_idx
      key_len: 4
          ref: const
         rows: 5
     filtered: 100.00
        Extra: Using filesort
1 row in set, 1 warning (0.00 sec)
```

`possible_keys`が選択可能なインデックス、`key`が実際に選択されたインデックス、`rows`が調査される行の見積もり。

なので、ちゃんと`pont_id_idx`がインデックスとして使われていて、調査する行も5行と小さい値になっていることがわかる。（インデックス貼る前と比較すればよかった...）


### 複合インデックス

さっきのクエリをもう一度確認すると、`Extra`の項目に`Using filesort`と書かれていることがわかる。これは、MySQL内部でsort処理が行われていることを示している。sort処理はデータベースにとって負担が大きい処理の一つなのでこれを解決する。（今回はsortする行が5行と少ないのでそこまで負担ではないはずだけどまあ練習で）

次のクエリでもともと貼ってあった`post_id`だけのインデックスを外して、 `post_id`と`created_at`の２つのカラムからなる複合インデックスを作る。

```sql
ALTER TABLE `comments` DROP INDEX `post_id_idx`, ADD INDEX `post_id_idx` (`post_id`, `created_at`);
```

![fig](./fig/multiple_column_indexes.png)
(「達人が教えるWebパフォーマンスチューニング」より)

`EXPLAIN`ステートメントの結果がどのように変わったのか確認する。

```sql
EXPLAIN SELECT * FROM `comments` WHERE `post_id` = 100 ORDER BY `created_at` DESC LIMIT 3\G
```
結果
```
*************************** 1. row ***************************
           id: 1
  select_type: SIMPLE
        table: comments
   partitions: NULL
         type: ref
possible_keys: post_id_idx
          key: post_id_idx
      key_len: 4
          ref: const
         rows: 5
     filtered: 100.00
        Extra: Backward index scan
1 row in set, 1 warning (0.00 sec)
```

`Extra`の`Using filesort`が`Backward index scan`に変わった。
上図のように昇順に並んでいるインデックスを逆向きに読んだことを表している。逆向きに読んでしまっているので、これを降順インデックスに変えればさらに処理が少なくなる可能性がある。

```sql
ALTER TABLE `comments` DROP INDEX post_id_idx, ADD INDEX post_id_idx(`post_id`,`created_at`DESC);
```

再び`EXPLAIN`を実行

```sql
EXPLAIN SELECT * FROM `comments` WHERE `post_id` = 100 ORDER BY `created_at` DESC LIMIT 3\G
```

結果

```
*************************** 1. row ***************************
           id: 1
  select_type: SIMPLE
        table: comments
   partitions: NULL
         type: ref
possible_keys: post_id_idx
          key: post_id_idx
      key_len: 4
          ref: const
         rows: 5
     filtered: 100.00
        Extra: NULL
1 row in set, 1 warning (0.01 sec)
```

`Backward index scan`が`NULL`になった。

### クラスターインデックスでのインデックスチューニング

#### インデックスの種類

MySQLをはじめとする多くのデータベースには、プライマリインデックスとセカンダリインデックスと呼ばれる２種類のインデックスがある。

- プライマリインデックス:
  - カラムに格納される値がすべてユニーク
  - int型でシーケンシャルな数字を自動付与することが多いが、UUIDなどをつけることもある
  - 2つのカラムデータの組み合わせがユニークであれば複合インデックスをプライマリインデックスにすることもできる
- セカンダリインデックス:
  - プライマリインデックス以外のインデックスは全てセカンダリインデックス

#### クラスターインデックスとは

MySQL（InnoDBストレージエンジン）のプライマリーキーはクラスターインデックスになっている。

クラスターインデックスとは、プライマリーインデックスのツリー構造の先に、データが含まれている構造のこと。

![](fig/primary_index_tree.png)
(「達人が教えるWebパフォーマンスチューニング」より)

プライマリーインデックスツリーの葉にはデータがある。

![](fig/secondary_index_tree.png)
(「達人が教えるWebパフォーマンスチューニング」より)

セカンダリインデックスツリーの葉はプライマリーキーの値。中のデータを参照するためには、プライマリインデックスツリーを走査する必要がある。

#### セカンダリインデックスの特徴を生かした検索の効率化

たとえば、「あるユーザーのコメント数を数えたい」という場合を考える。

クエリは次のようになる。

```sql
SELECT COUNT(*) FROM comments WHERE user_id = 123;
```

どのように実行されるのか確認しておく、

```sql
EXPLAIN SELECT COUNT(*) FROM comments WHERE user_id = 123\G
```
結果
```
*************************** 1. row ***************************
           id: 1
  select_type: SIMPLE
        table: comments
   partitions: NULL
         type: ALL
possible_keys: NULL
          key: NULL
      key_len: NULL
          ref: NULL
         rows: 99505
     filtered: 10.00
        Extra: Using where
1 row in set, 1 warning (0.01 sec)
```

このクエリを高速に行うためには、まず`user_id`にインデックスが必要なので作成する。

```sql
ALTER TABLE `comments` ADD INDEX `idx_user_id`(`user_id`);
```

今回のやりたいことは「`comments`テーブルから特定のユーザーがしたコメントの数を数える」ことなので、`user_id`のセカンダリインデックスツリーの調査だけで「コメントの数」はわかってしまう。つまりプライマリインデックスツリーまで走査する必要がないのでその分効率的。

このようにセカンダリインデックスに含まれる情報だけで結果が返せる最適化をCoverning Indexという。

Coverning Indexでクエリを解決できているかどうかは`EXPLAIN`で確認できる。

```sql
EXPLAIN SELECT COUNT(*) FROM comments WHERE user_id = 123\G
```

結果

```
*************************** 1. row ***************************
           id: 1
  select_type: SIMPLE
        table: comments
   partitions: NULL
         type: ref
possible_keys: idx_user_id
          key: idx_user_id
      key_len: 4
          ref: const
         rows: 100
     filtered: 100.00
        Extra: Using index
1 row in set, 1 warning (0.01 sec)
```

`Extra`に`Using index`がついていることから、Coverning Indexでクエリを解決できていることがわかる。

### アンチパターン

- インデックスが増えるとデータ更新時の負荷が高くなるので、インデックスの作りすぎもよくない。
- クエリ実行時に1つのテーブルに対して同時に使われるインデックスの個数は1個なので、複数のインデックスがあっても、データの絞り込みが速くなるわけではない。
- 様々な条件で検索する機能がある場合は、頻繁に使われるものだけインデックスを用意し、それ以外は「ORDER BY狙いのインデックス」を作成していくと良いらしい。

### その他のインデックス

MySQLがサポートするその他のインデックス

- 全文検索インデックス
  - データベース中に格納されるテキストデータから、特定の文字列を含む行の検索（`LIKE`クエリ）は全件走査してしまうので遅い。
  - MySQLでは全文検索インデックスを簡単に付与することができるが、検索や更新の負荷が高いので注意が必要。
  - MySQLでやらずにElasticsearchなどの専門の全文検索エンジンの利用も考慮するべき。
- 空間インデックス
  - 地図上の複数の座標を結んだ多角形の内側にある座標の検索などが簡単に、かつ高速に行うことができるようになる。

## N+1問題の発見と解消

### 事前準備

private-isuのappをruby実装からgo実装に変更する。（最初にやっておけばよかった...）

#### go実装への切り替え

```command
sudo systemctl stop isu-ruby
sudo systemctl disable isu-ruby
sudo systemctl start isu-go
sudo systemctl enable isu-go
```

動作確認

```command
systemctl status isu-go
```

activeになっていれば問題ない。

あとは、ブラウザからもアクセスできるか確認。

#### ファイルディスクリプタの上限変更

go実装に変更してベンチを回したら`too many open files`が出てしまったので、ファイルディスクリプタの上限を上げておく。

一応今の設定値確認

```
ulimit -n
```

1024だったので思い切って65536くらいまで上げてしまう。

`/etc/security/limits.conf`の末尾にこれを追加

```
* soft nofile 65536
* hard nofile 65536
```

VM再起動

```
sudo reboot
```

#### ベンチマーク再計測

appをgoに切り替えたので再び計測しておく。
もういちどVMにログインしてログのローテート

```command
sudo mv /var/log/mysql/mysql−slow.log /var/log/mysql/mysql−slow.log.bak3
sudo systemctl restart mysql
```

ベンチマーク実行

```command
cd /home/isucon/private_isu.git/benchmarker
./bin/benchmarker -u ./userdata -t http://localhost
```

slow-queryログをpt-query-digestにかける

```command
sudo pt-query-digest /var/log/mysql/mysql−slow.log > ~/pt-query-digest.log4
less ~/pt-query-digest.log4
```

### N+1問題の例

`~/pt-query-digest.log4`に次のようなログがある。

```
# Query 5: 1.76k QPS, 0.12x concurrency, ID 0x396201721CD58410E070DA9421CA8C8D at byte 69723609
# This item is included in the report because it matches --limit.
# Scores: V/M = 0.00
# Time range: 2022-09-18T07:51:10 to 2022-09-18T07:52:13
# Attribute    pct   total     min     max     avg     95%  stddev  median
# ============ === ======= ======= ======= ======= ======= ======= =======
# Count         15  110984
# Exec time      8      7s    15us    19ms    67us   144us   297us    26us
# Lock time     42   116ms       0     4ms     1us     1us    21us     1us
# Rows sent      0 108.38k       1       1       1       1       0       1
# Rows examine   0 108.38k       1       1       1       1       0       1
# Query size     3   4.01M      36      39   37.90   36.69    0.17   36.69
# String:
# Databases    isuconp
# Hosts        localhost
# Users        isuconp
# Query_time distribution
#   1us
#  10us  ################################################################
# 100us  ####
#   1ms  #
#  10ms  #
# 100ms
#    1s
#  10s+
# Tables
#    SHOW TABLE STATUS FROM `isuconp` LIKE 'users'\G
#    SHOW CREATE TABLE `isuconp`.`users`\G
# EXPLAIN /*!50100 PARTITIONS*/
SELECT * FROM `users` WHERE `id` = 932\G
```

`Exec time`をみると、95%のクエリが144マイクロ秒以下なので1クエリごとでみると高速だが、110984回も呼び出されているため5番目に負荷の大きいクエリになっている。

このクエリは [makePosts関数](https://github.com/catatsuy/private-isu/blob/096802b1d54481105624d7010c531e3b6328b170/webapp/golang/app.go#L174)で呼び出されているが、この関数の中では、

1. [`comments`テーブルから特定の`post_id`のコメントの最新3件を取得し](https://github.com/catatsuy/private-isu/blob/096802b1d54481105624d7010c531e3b6328b170/webapp/golang/app.go#L188)
2. [そのコメントに対して、さらにuser情報を引く](https://github.com/catatsuy/private-isu/blob/096802b1d54481105624d7010c531e3b6328b170/webapp/golang/app.go#L193-L198)

という処理が実装されている。つまり、最初の1回のSQL実行に対して、ループの中で何倍ものuser情報を求めるクエリが発生している。

このように、1回のクエリで得た結果の件数(N個)に対して、関連する情報を集めるため、N回以上のクエリを実行してしまうことでアプリケーションのレスポンス速度の低下やデータベースの負荷の原因になることをN+1問題という。

### N+1問題の見つけ方

APMやプロファイラ、フレームワークの機能でN+1に警告を出してくれたりするらしい。
しかし、今回はpt-query-digestの結果を見て実行数が多いクエリから探していく。

また、N+1を発見してもそれがボトルネックになっているのかは、確認が必要である。あまり呼ばれていない関数や速度が求められていない場所の高速化に時間をかけないように注意。

