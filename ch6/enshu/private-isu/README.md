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
