

## private-isuの場合
- [private-isu/webapp/golang](https://github.com/catatsuy/private-isu/tree/master/webapp/golang)


- [godoc](https://pkg.go.dev/)


- [net/html](https://pkg.go.dev/golang.org/x/net/html)


- [gsm](https://pkg.go.dev/github.com/bradleypeabody/gorilla-sessions-memcache#section-readme)


- [sqlx](https://pkg.go.dev/github.com/jmoiron/sqlx)


## データベースによく使用されるパッケージ
- database/sql
- github.com/jimoiron/sqlx
- github.com/ent/ent
- gorm.io/gorm
- github.com//kyleconroy/sqlc
- github.com/volatiletech/sqlboiler/v4


## ISUCONで使用されるパッケージ
### 標準パッケージ
| パッケージ | 年 | 用途 |
| --- | --- | --- |
|bytes|12f,11q|バイトを扱う|
|context|12q|コンテキストを扱う|
|crypto/ecdsa|11q|ECDSAを扱う|
|crypto/elliptic|10f|楕円曲線暗号を扱う|
|crypto/x509|10f|X.509を扱う|
|database/sql|12q,12f,11q,11f,10q,10f|データベースに接続する|
|database/sql/driver|12q|データベースに接続する|
|encoding/base64|10f|Base64を扱う|
|encoding/csv|12q,12f,10q|CSVファイルを扱う|
|encoding/json|12q,12f,11q,10q|JSONを扱う|
|encoding/pem|10f|PEMを扱う|
|errors|12q,11q,11f|エラーを扱う|
|fmt|12q,12f,11q,11f,10q,10f|フォーマットを扱う|
|io|12q,12f,11f|入出力を扱う|
|io/ioutil|11q,10q,10f|入出力を扱う|
|math|12f,11f|数学を扱う|
|math/rand|12f,11q,11f|乱数を扱う|
|net/http|12q,12f,11q,11f,10q|HTTPを扱う|
|net/url|11f|URLを扱う|
|os|12q,12f,11q,11f,10q|OSを扱う|
|os/exec|12q,12f,11q,11f,10q|コマンドを実行する|
|path/filepath|12q,10q|ファイルパスを扱う|
|reflect|12q|リフレクションを扱う|
|regexp|12q|正規表現を扱う|
|sort|12q,11q,11f|ソートを扱う|
|strconv|12q,12f,11q,11f,10q|文字列を扱う|
|strings|12q,12f,11q,11f,10q|文字列を扱う|
|sync|11f,10f|同期を扱う|
|time|12q,12f,11q,11f,10f|時間を扱う|

### 外部パッケージ
| パッケージ | 年 | 用途 |
| --- | --- | --- |
|github.com/dgrijalva/jwt-go|11q|JWTを扱う|
|github.com/go-sql-driver/mysql|12q,12f,11q,11f,10q|MySQLを扱う|
|github.com/gofrs/flock|12q|ファイルロックを扱う|
|github.com/golang/protobuf/proto|10f|Protocol Buffersを扱う|
|github.com/google/uuid|12f|UUIDを扱う|
|github.com/gorilla/session|11q,11f|セッションを扱う|
|github.com/jmoiron/sqlx|12f,11q,11f,10q,10f|データベースを扱う|
|github.com/labstack/echo/v4|12q,12f,11q,11f,10q|Webフレームワークを扱う|
|github.com/labstack/echo/v4/middleware|12q,12f,11f,10q|Webフレームワークのミドルウェアを扱う|
|github.com/labstack/gommon/log|12,11q,10q|ログを扱う|
|github.com/lestrrat-go/jwx/v2/jwa|12q|JWAを扱う|
|github.com/lestrrat-go/jwx/v2/jwk|12q|JWKを扱う|
|github.com/lestrrat-go/jwx/v2/jwt|12q|JWTを扱う|
|github.com/mattn/go-sqlite3|12q|SQLiteを扱う|
|github.com/oklog/ulid/v2|12f|ULIDを扱う|
|github.com/pkg/errors|11f|エラーを扱う|
|github.com/SherClockHolmes/webpush-go|10f|WebPushを扱う|
|golang.org/x/crypto/bcrypt|12f,11f|ハッシュを扱う|
|google.golang.org/protobuf/types/known/timestamppb|10f|timestampを扱う|

### 参考

[12q1](https://github.com/isucon/isucon12-qualify/blob/main/webapp/go/isuports.go)
[12q2](https://github.com/isucon/isucon12-qualify/blob/main/webapp/go/sqltrace.go)
[12f1](https://github.com/isucon/isucon12-final/blob/main/webapp/go/main.go)
[12f2](https://github.com/isucon/isucon12-final/blob/main/webapp/go/admin.go)
[11q](https://github.com/isucon/isucon11-qualify/blob/main/webapp/go/main.go)
[11f1](https://github.com/isucon/isucon11-final/blob/main/webapp/go/main.go)
[11f2](https://github.com/isucon/isucon11-final/blob/main/webapp/go/db.go
[11f3](https://github.com/isucon/isucon11-final/blob/main/webapp/go/util.go)
[10q](https://github.com/isucon/isucon10-qualify/blob/master/webapp/go/main.go)
[10f1](https://github.com/isucon/isucon10-final/blob/master/webapp/golang/notifier.go)
[10f2](https://github.com/isucon/isucon10-final/blob/master/webapp/golang/db.go)
[10f3]()

## Reference
### [達人が教えるWebパフォーマンスチューニング　〜ISUCONから学ぶ高速化の実践](https://amzn.to/3A3cZI8)