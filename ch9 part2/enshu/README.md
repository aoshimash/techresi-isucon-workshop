# ch9/enshu

## Requirements

- multipass


## Usage

インスタンスの立ち上げ

```
$ multipass launch --name enshu --cpus 1 --disk 8G --mem 2G --cloud-init enshu.cfg 20.04
```

インスタンスの一覧を表示（状態確認）

```
$ mulitpass list
```

インスタンスへのログイン


```
$ multipass shell enshu
```

インスタンスの削除

```
$ multipass delete enshu
```


削除済みインスタンスをすべて取り除く

```
$ multipass purge
```
