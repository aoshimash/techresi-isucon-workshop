# 計測手法の検討

## そもそも計測したいものはなにか

- アクセスログ: 遅いリクエストを見つけたい
- MySQLログ: 遅いクエリを見つけたい
- システム情報: 負荷の高いプロセスを見つけたい。あとCPUとメモリの使用率、ディスクIOとかもみたい
- アプリケーションのプロファイリング: コードのどこが遅いのか見つけたい

## newrelicがISUCONで使えるか調査

### 調査結果

newrelicなら上の情報全部取得できた。
でも使いたくない。

### 理由

- 無料枠だとやはり色々制限がある
  - ユーザー数とか登録できるエンティティの数とか
  - 有料だと練習で使いにくくなる
  - 本番はNew Relicが無料ライセンス配ってくれるけど、普段使っていないツールは本番でも使えない
- 上位勢はあんまり使ってない
- 監視システムの構築が多少面倒
  - 時間を使うべきところはここじゃない
  - 監視システムの導入でトラブルと死ぬ

### 他の監視SaaS

- datadog
  - 無料枠だとやれることが限られるからnewrelicと同じ理由でなし
- 自分で作る
  - Grafana + Prometehus + Loki とかで自分でホストする
  - セルフホストはつらい

### まとめ

やっぱり計測は**構成がシンプル**であることと**使い慣れていること**が大事！

## 代替案

### アクセスログ

Nginxの遅いリクエストを見つけたい！

去年と一緒で[alp](https://github.com/tkuchiki/alp)を使うか、[kataribe](https://github.com/matsuu/kataribe)を使うか。取れる情報同じなのでどちらでもいい。

どちらを使うにしろnginxのログフォーマットはツールに合わせて変更が必要。

[kataribeはクエリパラメタが違うリクエストを別々に集計したりまとめたりすることができるらしい](https://blog.hog.as/entry/2020/09/13/052113)ので、一回kataribe使ってみたい。

### MySQLログ

スロウクエリの解析をして遅いクエリを見つけたい！

いつもどおり[pt-query-digest](https://docs.percona.com/percona-toolkit/pt-query-digest.html)でやる。
（スロウクエリファイルを作らなくても、MySQLプロトコルデータをtcpdumpでキャプチャすればpt-query-digest使えるらしい。知らなかった。知ったところでやらないけど。）

### システム情報

htop, ctop, gtop なんでもいいのでは。
[dstatは開発者がRedHatに怒って開発やめちゃったから](https://github.com/scottchiefbaker/dool)後継の[Dool](https://github.com/scottchiefbaker/dool)とかいいかも。

[netdata](https://www.netdata.cloud/)なら入れてもよさそう。（netdata用のポート開放が必要）

### アプリケーションプロファイリング

関数の実行時間などが調べられるので、コードのどこが遅いのかわかる。

Goなら素直に[pprof](https://github.com/google/pprof/)使うのが良さそう。

↑プロファイラーの使い方勉強会は別にやりましょう！


### その他

これらのツールを全部手動でたたくわけじゃなくて、計測ツールをまとめて実行して解析してグラフ化してレポート作るところまで自動化するつもりだからちょっとまってて。（[isucon-template](https://github.com/Tech-Residence/isucon-template)に追加してくよ。）

計測レポートはslack or discordに送ったほうが便利かな？
