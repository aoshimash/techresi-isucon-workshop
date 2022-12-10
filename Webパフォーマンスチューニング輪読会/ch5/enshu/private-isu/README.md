# techresi-isucon-workshop/ch5/enshu/private-isu

## Requirements

- multipass

```
brew install multipass
```

## Usage


### VMの起動

以下のコマンドで、アプリケーションサーバーもベンチマークサーバーも同梱されたVMを立てる。
リソースのミニマムはCPU=2,disk=16G,memory=2G程度。CPUコアが多いと構築にかかる時間が短縮できる。

```
multipass launch --name private-isu --cpus 4 --disk 16G --mem 4G --cloud-init standalone.cfg 20.04
```
(2CPUだと手元の環境で構築したときに30分以上かかったので、4CPUくらいがいいと思う。)

cloud-initの処理には時間がかかるため、上記のコマンドはタイムアウトすると思うが、VM内部では処理が進んでいるので問題ない。
VMのセットアップが完了したかどうかは次の手順で確認できる。

#### VMの初期化が完了しているか確認

VMにログインしてから

```
multipass shell private-isu
```

cloud-initのログを確認

```
sudo tail -f /var/log/cloud-init-output.log
```

完了までは20分程度はかかると思う。

### Webアプリへアクセス

まずは、次のコマンドをホストOSで実行してVMのIPを取得し、ブラウザからアクセスできるか確認する。

```
mulitpass list
```


#### ベンチマーカーの実行

VMにログイン

```
multipass shell private-isu
```

ベンチマーカーディレクトリに移動

```
cd /home/isucon/private_isu.git/benchmarker
```

実行

```
./bin/benchmarker -u ./userdata -t http://localhost/
```

### VMの停止

```
multipass stop private-isu
```

### VMの削除

```
multipass delete private-isu
```

## Reference

こちら(https://github.com/matsuu/cloud-init-isucon/tree/main/private-isu)の `standalone.cfg` に今回の演習に必要なパッケージを追加しただけ。
