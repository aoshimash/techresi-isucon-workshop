# ISUCON本ゆる輪読会#8 

## 1. リバースプロキシを前段に置く理由
### 1.1 アプリケーションとプロセス・スレッドについて
- 負荷分散（ロードバランス）
- コンテンツのキャッシュ
- HTTPS通信の終端

### 1.2 リバースプロキシで得られるメリット
- 転送時のデータ圧縮
- リクエストとレスポンスのバッファリング
- リバースプロキシとアップストリームサーバー間でのコネクション管理

## 2. nginxのアーキテクチャ
### 2.1 C10K問題とは
### 2.2 nginxの設定確認
``` /etc/nginx/nginx.conf
include /etc/nginx/conf.d/*.conf;
include /etc/nginx/sites-enabled/*;
```

``` /etc/nginx/sites-available/isucon.conf
server {
    listen 80;

    client_max_body_size 10m;
    root /home/isucon/private_isu/webapp/public;
    location / {
        proxy_set_header Host $host;
        proxy_pass http://localhost8080;
    }
}
```

静的ファイルの配信をアプリケーションサーバーからnginx経由で行うように切り替えることでアプリケーションサーバーへの負荷を減らす。
``` /etc/nginx/sites-available/isucon.conf
server {
    listen 80;

    client_max_body_size 10m;
    root /home/isucon/private_isu/webapp/public;

    location /css/ {
        root /home/isucon/private_isu/webapp/public/;
        expires 1d;
    }

    location /js/ {
        root /home/isucon/private_isu/webapp/public/;
    }

    location / {
        proxy_set_header Host $host;
        proxy_pass http://localhost8080;
    }
}
```
公開するディレクトリをURLのパスにマッピングしている。

## 3. nginxを活用した高速化手法
### 3.1 転送時のデータ圧縮

```
gzip on;
    gzip_types text/css text/javascript application/javascript application/x-javascript application/json;
gzip_min_length 1k;
```

### 3.2 リクエスト・レスポンスのバッファリング


### 3.3 アップストリームサーバーとのコネクション管理
```/etc/nginx/sites-available/isucon.conf
location / {
    proxy_http_version 1.1;
    proxy_set_header connection "";
    proxy_pass http://app;

}
```

```
upstream app {
    server localhost:8080;
    keepalive 32;
    keepalive_requests 1000;
}
```

### 3.4 TLS通信の高速化

### 3.5 Linuxカーネルパラメータの設定
