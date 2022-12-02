## 9.6 CPU利用率

チューニングを行う上で常に注意することが必要な値として**CPU利用率**が挙げられる。

CPU（Central Processing Unit）とは、汎用的な処理を行う部品。さまざまな演算を処理する機能を持っており、GPU（Graphics Processing Unit）などの専用チップが搭載されてない限りはCPU上で処理が行われる。

Linuxにおいても多くの演算処理をCPUで行っており、topコマンドで目視することができる。

```jsx
$ top

top - 22:58:33 up 2 days, 20:56,  1 user,  load average: 0.00, 0.00, 0.00
Tasks:  88 total,   1 running,  87 sleeping,   0 stopped,   0 zombie
%Cpu(s):  0.3 us,  0.0 sy,  0.0 ni, 99.0 id,  0.3 wa,  0.0 hi,  0.3 si,  0.0 st
MiB Mem :   1974.1 total,    829.8 free,    149.6 used,    994.7 buff/cache
MiB Swap:      0.0 total,      0.0 free,      0.0 used.   1735.2 avail Mem 

    PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND                                                                                                                                                                                   
  23077 ubuntu    20   0   15600   4348   3132 S   0.3   0.2   0:00.04 sshd                                                                                                                                                                                      
  23091 ubuntu    20   0   10680   3224   2672 R   0.3   0.2   0:00.09 top
```

topコマンド実行中に1を入力もしくはtopコマンドの引数に−1を指定することでそれぞれのコアの状態を確認できる。

```jsx
$ top -1

top - 23:04:18 up 2 days, 21:02,  1 user,  load average: 0.00, 0.00, 0.00
Tasks:  87 total,   1 running,  86 sleeping,   0 stopped,   0 zombie
%Cpu0  :  0.0 us,  0.3 sy,  0.0 ni, 99.7 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
MiB Mem :   1974.1 total,    828.9 free,    150.0 used,    995.2 buff/cache
MiB Swap:      0.0 total,      0.0 free,      0.0 used.   1734.8 avail Mem 

    PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND                                                                                                                                                                                   
  23233 ubuntu    20   0   10580   3104   2656 R   0.3   0.2   0:00.16 top                                                                                                                                                                                       
      1 root      20   0  103384  10908   7452 S   0.0   0.5   0:15.29 systemd
```

もともと%Cpu(s)という表記だった部分が%Cpu(0)になっている。

%Cpu(s)が表示されている場合は、すべてのコアにおける平均値が表示されている。

例えば「すべてのコアの平均値が50%」と「%Cpu(0)が100%、%Cpu(1)が0%」の状態を見分けるためにも複数行表示しておくことを推奨する。

### us - User : ユーザ空間におけるCPU利用率

「システムコールを利用するLinux OS上のアプリケーションが動作する部分」であるユーザ空間におけるCPU利用率を指す。

まさに実行されている個々のWebアプリケーションが多くのCPUを利用している場合に上昇する値。

### sy - Sysmtem : カーネル空間におけるCPU利用率

「Linux Kernel内の処理」であるカーネル空間におけるCPU利用率を指す。

プロセスのforkが多く発生している環境やコンテキストスイッチを行っている時間が長くなっている環境においてはカーネル空間の処理が大きくなるためsyの値が上昇する。

ユーザ空間とカーネル空間の違いに関しては[こちら](http://linux-dvr.biz/archives/39)の資料がわかりやすかった。

### ni - Nice : nice値（優先度）が変更されたプロセスのCPU利用率

Linuxのプロセスには複数の優先度が定義されている。

このプロセスの優先度を表すのがnice値である。psコマンドを用いて、nice値を確認できる。

```jsx
$ ps -axf -o pid,ppid,ni,args
    PID    PPID  NI COMMAND
（中略）
  1364       1   0 sshd: /usr/sbin/sshd -D [listener] 0 of 10-100 startups
  23124    1364   0  \_ sshd: ubuntu [priv]
  23219   23124   0      \_ sshd: ubuntu@pts/0
  ç   23219   0          \_ -bash
  23457   23220   0              \_ ps -axf -o pid,ppid,ni,args
（以下略）
```

NIと記載されている箇所がプロセスのnice値である。Linuxでは、−20（最高優先度）から19（最低優先度）までが定義されており、数字が小さい方が優先度が高い。

上記の場合のpsコマンド（PID=23457)のnice値は0であることがわかる。これはpsコマンドの親プロセスであるbashプロセス（PID=23457)のnice値が0であることから由来している。

またpsコマンドは以下のように優先度を下げた状態で実行することもできる。

```jsx
$ nice -n19 ps -axf -o pid,ppid,ni,args

  1364       1   0 sshd: /usr/sbin/sshd -D [listener] 0 of 10-100 startups
  23124    1364   0  \_ sshd: ubuntu [priv]
  23219   23124   0      \_ sshd: ubuntu@pts/0
  23220   23219   0          \_ -bash
  23475   23220  19              \_ ps -axf -o pid,ppid,ni,args
```

優先度を下げるメリットとして、より負荷の高いデーモンプロセスなどではあえて性能を下げることで他の優先すべきプロセスにCPUリソースを与えることができる。

### id -  Idle : 利用されてないCPU

idは「利用されてないCPU」の割合を示している。

### wa - Wait : I/O処理を待っているプロセスのCPU利用率

マルチスレッドの処理を行わない場合、プロセスがI/O処理を行っていると、そのプロセスI/O処理が終わるまで他の処理を行うことができない。waはプロセスの中でもディスクなどへのI/O処理を待っているプロセスCPU利用率である。

waの値が上がっている場合は、ディスクなどへの読み書き終了を待っているプロセスが多く存在していることを示している。

### hi - Hardware Interrupt : ハードウェア割り込みプロセスの利用率

9章4節「Linuxのネットワーク」で解説した、ハードウェア割り込みを利用しているプロセスの利用率を指す。

### si - Soft Interrupt : ソフト割り込みプロセスの利用率

9章4節「Linuxのネットワーク」で解説した、ソフト割り込みを利用しているプロセスの利用率を指す。

### st - Steal : ハイパーバイザによって利用されているCPU利用率

stは、パブリッククラウドなど仮想化された環境のLinux上で利用されるCPU利用率である。物理的なホストにインストールされたOS上でVMを起動する場合、物理的なCPUリソースをVMに割り当ててVMのプロセスを動作させる。このときVM上のLinuxが認識しているCPUは仮想環境によって作られたCPUであり、物理的なCPU演算が必要なときのみホスト側のCPUリソーつを使って演算を行う。

しかし、ホストOS側もコンテキストスイッチを用いてプロセスを動作させているため、どうしてもVMが必要なタイミングでCPUリソースを割り当てられない時間が存在する。特に同じCPUコアに複数のVMが割り当てられた上にVM内にCPU負荷の高いプロセスが存在した場合はほかのVMの影響を受けることが多くなる。そのような「利用できるはずができなかったCPU時間」の率を示しているのがst。

## 9.7 Linuxにおける効率的なシステム設計

本節はWebサービスを提供する際に頻出するLinuxのパラメータとその効果について説明する。

### ulimit

ulimit(user limit)は、プロセスが利用できるリソースの制限を設定する概念。各プロセスはどのリソースをどのくらい利用できるかについて、制限をかけている。

```jsx
$ ulimit -a
core file size          (blocks, -c) 0
data seg size           (kbytes, -d) unlimited
scheduling priority             (-e) 0
file size               (blocks, -f) unlimited
pending signals                 (-i) 7633
max locked memory       (kbytes, -l) 65536
max memory size         (kbytes, -m) unlimited
open files                      (-n) 1024
pipe size            (512 bytes, -p) 8
POSIX message queues     (bytes, -q) 819200
real-time priority              (-r) 0
stack size              (kbytes, -s) 8192
cpu time               (seconds, -t) unlimited
max user processes              (-u) 7633
virtual memory          (kbytes, -v) unlimited
file locks                      (-x) unlimited
```

ulimitはプロセス単位で設定されており、ulimit -aを実行すると表示されるのは現在起動しているシェルプロセスの制限。

## 9.8 Linuxカーネルパラメータ

今回はWebサービスを提供する際に利用するカーネルパラメータを取り上げ、設定を変更する方法を紹介する。

### net.core.somaxconn

net.core.somaxconnはTCPソケットが受け付けた接続要求を格納する、キューの最大長のこと。backlog > net.core.somaxconnのとき、キューの大きさは暗黙にnet.core.somaxconnに切り詰められる。

backlogとは接続要求のキューのこと。

```jsx
sysctlによってカーネルパラメータの設定値を確認する
$ sysctl net.core.somaxconn
net.core.somaxconn = 4096
```

この値をさらに大きくすることで、接続数を増やすことができる。今回は倍の8192まで増加させる。

```jsx
$ sudo sysctl -w net.core.somaxconn=8192
net.core.somaxconn = 8192

$ sysctl net.core.somaxconn
net.core.somaxconn = 8192

恒久的に書き換える場合は、/etc/sysctl.conf、または/etc/sysctl.d/配下のファイルに記載する
$ tail /etc/sysctl.conf
#net.ipv4.conf.all.log_martians = 1
#

###################################################################
# Magic system request Key
# 0=disable, 1=enable all, >1 bitmask of sysrq functions
# See https://www.kernel.org/doc/html/latest/admin-guide/sysrq.html
# for what other values do
#kernel.sysrq=438

net.core.somaxconn = 8192

# sudo sysctl -pコマンドで更新する
$ sudo sysctl -p
```

### net.ipv4.ip_local_port_range

LinuxでTCP/UDPの通信を行う際、サーバー側のポートはHTTPなら80,HTTPSなら443がよく利用される。これらはSystem Portsと呼ばれており、1〜1023番までがこれにあたる。

対して、パケットの通信を行うためにはクライアント側にもポートが必要である。クライアント側のポートなど、動的に利用できるポート領域はEphemeral Portsと呼ばれており、Linux5.4環境においては32768〜60999番ポートが利用されることがデフォルトで設定されている。ipv4.ip_local_port_rangeは、動的確保するポートの範囲を設定するカーネルパラメータである。

```jsx
$ sysctl net.ipv4.ip_local_port_range
net.ipv4.ip_local_port_range = 32768    60999
```

net.ipv4.ip_local_port_rangeの設定を極端に小さくする

```jsx
curlコマンドで何度かgoogleに接続するとエラーになる
$ curl -vvv https://www.google.com/
*   Trying 142.250.207.36:443...
* TCP_NODELAY set
* Immediate connect fail for 142.250.207.36: Cannot assign requested address
*   Trying 2404:6800:4004:824::2004:443...
* TCP_NODELAY set
* Immediate connect fail for 2404:6800:4004:824::2004: Network is unreachable
* Closing connection 0
curl: (7) Couldn't connect to server
```

ローカルポートを割り当てることに失敗するため、Cannot assign requested addressというエラーが出力された。

利用可能なポートが2つだけ存在しているため、curlコマンドを複数回実行することで、ポートの確保に成功しHTTPリクエストが成功したり、ポートが確保できずHTTPリクエストが失敗したりする様子が確認できる。

これはデフォルトの設定でもポートが足りなくなる場合がある。

```jsx
一時的に書き換える場合はsysctl -wコマンドを利用する
$ sudo sysctl -w net.ipv4.ip_local_port_range="1024 65535"
net.ipv4.ip_local_port_range = 1024 65535

$ sysctl net.ipv4.ip_local_port_range
net.ipv4.ip_local_port_range = 1024     65535
```

```jsx
恒久的に書き換える場合は、/etc/sysctl.conf、または/etc/sysctl.d/配下のファイルに記載する
$ tail /etc/sysctl.conf
#net.ipv4.conf.all.log_martians = 1
#

###################################################################
# Magic system request Key
# 0=disable, 1=enable all, >1 bitmask of sysrq functions
# See https://www.kernel.org/doc/html/latest/admin-guide/sysrq.html
# for what other values do
#kernel.sysrq=438

net.core.somaxconn = 8192
net.ipv4.ip_local_port_range=1024 65535

# sudo sysctl -pコマンドで更新する
$ sudo sysctl -p
```

ipv4という文字列があるが、ipv6においてもこのパラメータを参照している。

Webサーバーが通信を待ち受ける際に、通常であればTCPの80番ポートへ接続することが多いがSocketファイルという特殊なファイルを生成し、そのファイルを通して接続を待機することもできる。

nginxの設定ファイルである80番ポートでなく、UNIX domain socketを用いる例を紹介

```jsx
server {
  80番ポートで接続を待機する際の設定(#を付けてコメントアウト済)
  # listen 80;
  ## /var/run/nginx.sock で接続を待機する際の設定
  listen unix: /var/run/nginx.sock;
<以下略>
```

この設定を利用した場合に、curlコマンドを用いてHTTPリクエストを送る例は以下

```jsx
$ curl --unix-socket /var/run/nginx.sock exemple.com
```

このようにUNIX domain socketを用いることで、サーバーのポートを消費せず通信をおこなうことができる。

private-isuにおける、nginxとWebアプリケーションの間の接続においてもUNIX domain socketを利用できる。

private-isuの初期状態では、Webアプリケーションは0.0.0.0:8080をListen(2)し、nginxが受け取ったリクエストをhttp://localhost:8080にプロキシする構成になっている。

次に、nginxとwebapp間の接続にUNIX domain socketを利用する設定例を紹介

Webアプリケーション側の設定を変更し、tmp/webapp.socketにてlisten(2)するように変更する。

Goの実装ではapp.go内でhttp.ListenAndServer()関数によってlisten(2)するアドレスをしている。

初期状態は以下

```jsx
log.Fatal(http.ListenAndServer(":8000", mux)
```

リスト 4 Go 実装を書き換える

```jsx
リスト 4 Go 実装を書き換える
## "/tmp/webapp.sock" listen (2) $3
listener, err = net.Listen ("unix", "/tmp/webapp.sock")
if err != nil {
  log.Fatalf("Failed to listen on /tmp/webapp.sock: %s.", err)
}

defer func() {
  err = listener.Close()
    if err != nil {
      log.Fatalf("Failed to close listener: %s.", err)
    }
  }
}()

## systemdなどから送信されるシグナルを受け取る
c = make (chan os.Signal, 2)
signal. Notify(c, os. Interrupt, syscall. SIGTERM)
go func() {
  <-C
  err = listener.Close()
  if err != nil {
    log.Fatalf("Failed to close listener: %s.", err)
  }
}()
log.Fatal (http.Serve (listener, mux))
```

このように書き換えた上でWebアプリケーションを再起動すると/tmp/webapp.sockでlisten(2)を行う。nginxでは、 proxy_passによってhttp://localhost:8000にプロキシする設定になっている。

リスト5 初期状態

```jsx
server {
<省略＞
  location / {
		proxy_set_header Host Shost;
    proxy_set_header X-Real-IP Sremote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_pass http://localhost:8080;
  }
}
```

upstreamディレクティブによって /tmp/webapp.sockをアップストリームサーバーとして指定し、
UNIX domain socket に接続。 

リスト6に修正を加えた設定を示します。

```jsx
リスト6 /tmp/webapp.sockをアップストリームサーバーとして指定した設定
upstream webapp {
	server unix:/tmp/webapp.sock;
}
server {
<省略>
  location / {
    proxy_set_header Host Shost;
    proxy_set_header X-Real-IP Sremote_addr;
    proxy_set_header X-Forwarded - For $proxy_add_x_forwarded_for;
    proxy_pass http://webapp;
  }
}
```

上記の設定を加えた上でnginx を再起動することで､ nginxとWebアプリケーション間の通信が
UNIX domain socket を介して行われ、一定の高速化を見込める。