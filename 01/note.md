# HTTP1.0のシンタックス：基本となる4つの要素

## HTTPの歴史

### バージョン
- 1990: HTTP/0.9
- 1996: HTTP/1.0
- 1997: HTTP/1.1
- 2015: HTTP/2

### RFC - HTTPの仕様策定文書
IETFが作ったRFCというルール文書で定義される。
ex) RFC XXXX(バージョン名)

## Goでのテストサーバーとcurl

Goの基礎: http://qiita.com/kazusa-qooq/items/40f9ea3e72406d845b10

通信内容をそのまま表示するだけのエコーサーバー

`go run filename.go` で実行

```go
package main

import (
  "fmt"
  "log"
  "net/http"
  "net/http/httputil"
)

func handler(w http.ResponseWriter, r *http.Request)  {
  dump, err := httputil.DumpRequest(r, true)
  if err != nil {
    http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
    return
  }
  fmt.Println(string(dump))
  fmt.Fprintf(w, "<html><body>hello</body></html>\n")
}

func main()  {
  var httpServer http.Server
  http.HandleFunc("/", handler)
  log.Println("start http listening :18888")
  httpServer.Addr = ":18888"
  log.Println(httpServer.ListenAndServe())
}
```

curlコマンド

```sh
curl --http1.0 http://localhost:18888/greeting
=> <html><body>hello</body></html>
```

## 0.9 to 1.0

要素  | 0.9 | 1.0
------|-----|----
メソッドとパス   | △ メソッドは無し    | ◯
ヘッダー         |   ✕  | ◯
ボディ           |  ◯   | ◯
ステータスコード |   ✕  | ◯

## ヘッダー
- 元々はメールシステムで使われていたアイデア。
- **フィールド名: 値** という形式
- 大文字・小文字の区別なし
- `X-`のものは開発者が任意に追加可能

#### curlコマンドでの確認

```sh
➜  ~ curl --http1.0 -v -H "X-Test: Hello" http://localhost:18888/greeting
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 18888 (#0)
> GET /greeting HTTP/1.0
> Host: localhost:18888
> User-Agent: curl/7.51.0
> Accept: */*
> X-Test: Hello
>
* HTTP 1.0, assume close after body
< HTTP/1.0 200 OK
< Date: Mon, 03 Jul 2017 01:27:32 GMT
< Content-Length: 32
< Content-Type: text/html; charset=utf-8
<
<html><body>hello</body></html>
* Curl_http_done: called premature == 0
```

### Content-Type
MIMEタイプを指定することで、ファイルの種類を特定する。

WebサーバがHTMLを送信する場合は、

```sh
Content-Type: text/html; charset=utf-8
X-Content-Type-Options: nosniff  // セキュリティ観点で付与する
```

## メソッド

指定されたアドレスにあるリソースに対する操作をサーバーに指示する。

### 主要なメソッド
- GET
- POST
- PUT
- DELETE
- HEAD

#### curlコマンドでの確認
--request, もしくは短縮形の-Xを使う

```
curl --http1.0 -v -X POST http://localhost:18888/greeting
```

## ステータスコード
サーバーがリクエストに対してどのように応答したかがわかる

 番 | 意味
--|--
100番台  |  処理中の情報の伝達
200番台  |  成功時のレスポンス
300番台  |  サーバーからクライアントへの命令。リダイレクトやキャッシュの使用
400番台  |  リクエストに関するエラー
500番台  |  サーバー内部のエラー



### リダイレクト

- 300番台の一部はリダイレクト処理をブラウザに大して指示するステータスコード
- 300番以外の時は、ブラウザはLocationヘッダの値を見て再度リクエストを送信する

ステータスコード  |  用途
--|--
301/308  |  一般的なリダイレクト
302/307  |  一時的な移動。モバイルサイトやメンテページへのジャンプ
303  |  ログイン時のリダイレクト等

#### curlコマンドでの確認
```
curl -L http://localhost:18888
```

nginxの設定例 https://nginx.org/en/docs/http/converting_rewrite_rules.html

## URL

### URLとURI
URI: URN（名前の付け方のルールを含む <br>
URL: 場所のみ

### URLの構造
例
> https://www.oreilly.co.jp/index.html

- スキーマ: https
- ホスト名: www.oreilly.co.jp
- パス: index.html

### ボディ
- ヘッダとの間に空行を挟んでそれ以降全てはボディー
- Content-Lengthにはボディのバイト数が入る

#### curlコマンドでのbodyの送信
-dオプションを使って以下のように書く(デフォルトではContent-Type: application/x-www-form-urlencoded)。-dオプションだとPOSTになる

```
curl -d "{\"hello\": \"world\"}" -H "Content-Type: application/json" http://localhost:18888

// jsonファイルの内容を送る場合
curl -d @test.json -H "Content-Type: application/json" http://localhost:18888
```

#### bodyを送るべきでない場合
GETメソッドでも送ることは可能だが、非セマンティックとされる
