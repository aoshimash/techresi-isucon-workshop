# techresi-isucon-workshop/ch6/enshu/private-isu

「Chapter6 リバースプロキシの利用」の勉強会で使う環境用のcloud-init。

## デプロイ

（eichan-serverでの利用を想定しているが、[networkインターフェイスが追加できるmultipass環境](https://multipass.run/docs/additional-networks)があれば使える。）

以下のコマンドでappサーバー3台とbenchサーバー1台の合計4台のVMが立ち上がる。


```
./create-vm.sh
```

次のコマンドで立ち上げたVMのステータスとIPを確認できる。

```
multipass list
```
<<<<<<< HEAD
(2CPUだと手元の環境で構築したときに30分以上かかったので、4CPUくらいがいいと思う。)

cloud-initの処理には時間がかかるため、上記のコマンドはタイムアウトすると思うが、VM内部では処理が進んでいるので問題ない。
VMのセットアップが完了したかどうかは次の手順で確認できる。

#### VMの初期化が完了しているか確認

VMにログインしてから

```
multipass shell private-isu-ch6
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
multipass shell private-isu-ch6
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
multipass stop private-isu-ch6
```

### VMの削除

```
multipass delete private-isu-ch6
```

## Reference

こちら(https://github.com/matsuu/cloud-init-isucon/tree/main/private-isu)の `standalone.cfg` に今回の演習に必要なパッケージを追加しただけ。

=======
>>>>>>> 24a16add3e5fe1431668cd37b02c2b32a695289d
