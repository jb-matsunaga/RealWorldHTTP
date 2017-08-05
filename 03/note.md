
# ３章 Go言語による HTTP/1.0クライアントの実装

HTTPの基礎となる箱
- メソッドとパス
- ヘッダー
- ボディ
- ステータスコード

## 3.1 Go言語を使う理由

### Go言語のメリット

- 他の言語と比べてコンパクトな言語仕様と豊富な標準ライブラリ → __比較的学習コストが低め__
- コンパイルが動的スクリプト言語の実行並みに高速で、型チェックが確実に行われる →__ 高い品質を保てる__
- 実行速度も高速で、マルチコアの性能を引き出しやすく、省メモリ → __高速・省メモリ__
- __クロスコンパイル__が簡単
> 例えば Mac で Windows 用、 Linux 用のバイナリを一気に生成するといったことができます。
- アウトプットが単一バイナリになるため配布しやすい → __導入インストールが楽__

### 本書でGoを採用する理由

- 教育用言語として優れている
  - 言語仕様が他の言語よりも少なく、他の言語ユーザーが見ても挙動が理解しやすいため、疑似言語として優れている。
  - コンパイル言語であり、構文や型のチェックが行われるため、入力の間違いに気付きやすい
  - 標準ライブラリのみを使ってHTTPのアクセスを行うプログラムが作成できる
  - 実際に、さまざまなウェブサービスのCLIクライアントの実装言語として使われている

- C/C++<br>
Pros → 実行速度が高速、バイナリもコンパクト<br>
Cons → 環境整えるまでが大変

- Python, Ruby, Node.js<br>
Pros → コンパイル不要ライブラリ豊富<br>
Cons → 単一バイナリにはならない、速度、消費メモリで負ける

<br>
→ __Go言語はそつなくさまざまな課題を解決してくれる言語__

> ##### オンラインチュートリアル
> [A Tour of Go](https://go-tour-jp.appspot.com/welcome/1)

> [A Tour of Goを終えたあなたにおすすめのGoを勉強するためのリソース \- すずけんメモ](http://suzuken.hatenablog.jp/entry/2017/07/21/121149)

> #### Goのデメリットはないのか
>
> [ksimka/go\-is\-not\-good: A curated list of articles complaining that go \(golang\) isn't good enough](https://github.com/ksimka/go-is-not-good)
>
> [Go言語がダメな理由 \| プログラミング \| POSTD](http://postd.cc/why-go-is-not-good/)
>
> - Generics (template) がない
> - 継承がない
> - 例外がない。まるで1970年代に設計されたかのようである。
> - 非知的なプログラマのためにデザインされている。

## 3.2 Go言語のAPIの構成
Go言語で提供されているHTTPのAPIは大きく分けて３つ

- 簡単に扱える `http.Get`, `http.Head`, `http.Post`, `http.PostForm`
- クッキーやプロキシを有効にして使うための `http.Client`
- 全機能に_アクセス_できる `http.Request/http.Client`
の組み合わせ

## 3.3 本章で取り上げるレシピ

| recipe                                        | method        | Golang API               |
| :--                                           | :--           | :--                      |
| [GET による情報取得](#3-4)                    | GET           | http.Get                 |
| [クエリー付き情報取得](#3-5)                  | GET           | http.Get                 |
| [HEADによるヘッダー取得](#3-6)                | HEAD          | http.Head                |
| [x-www-form-urlencodedでフォームの送信](#3-7) | POST          | http.PostForm            |
| [POSTでファイル送信](#3-8)                    | POST          | http.Post                |
| [multipart/form-dataでファイルの送信](#3-9)   | POST          | http.PostForm            |
| [クッキーの送受信](#3-10)                     | GET/HEAD/POST | http.Client              |
| [プロキシ](#3-11)                             | GET/HEAD/POST | http.Client              |
| [ファイルシステムへのアクセス](#3-12)         | GET/HEAD/POST | http.Client              |
| [自由なメソッド送信](#3-13)                   | なんでも      | http.Request/http.Client |
| [ヘッダーの送信](#3-14)                       | なんでも      | http.Request/http.Client |

<a id="3-4"></a>
## 3.4 GETメソッドの送信と、ボディ、ステータスコード、ヘッダーの送信

まずは、一番簡単なGETメソッドを送信するコード

```sh
$ curl http://localhost:18888
```

```sh
GET / HTTP/1.1
Host: localhost:18888
Accept: */*
User-Agent: curl/7.51.0
```

- [パッケージ \- はじめてのGo言語](http://cuto.unirita.co.jp/gostudy/post/go-package/)
- [A Tour of Go](https://go-tour-jp.appspot.com/basics/10)

```go
package main // ライブラリ以外のソースコードは必ずpackage mainから始まる

import ( // 必要なパッケージを取り込みます。Go言語ではそのファイルで宣言したパッケージ以外は使用できません。
  "io/ioutil"
  "log"
  "net/http"
)

func main() { // 全てのプログラムはmainパッケージのmain関数が最初に呼ばれます
  resp, err := http.Get("http://localhost:18888")
  // Go言語のエラー処理コードです。Go言語の関数は返り値としてエラーを返すので、それがnilかどうか確認します。例外処理はありません。
  if err != nil {
    // panicはエラーを表示させてプログラムを終了させる。
    // ライブラリ化する時はhttp.Get()のように、エラーを返り値の最後の項目として返すのがGoの流儀。
    panic(err)
  }
  // これは後処理コードです。defer を付与するとこの関数から抜けた後にこのぶんを実行します。ソケットからボディを読み込んだ後の処理です。
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body) // ボディの内容をバイト列として読み込みます
  if err != nil {
    panic(err)
  }
  log.Println(string(body))
}
```

#### エラーチェックを省いた最少のコード
```go
func main() {
  // resp変数に入っている要素が、http.Response型のオブジェクト。サーバから返ってきたさまざまな情報を全て格納している。
  resp, _ := http.Get("http://localhost:18888") // resp はhttp.Response型の変数
  defer resp.Body.Close() // ボディはbodyメンバ変数に格納されている。
  // io.Readerインターフェース化する
  body, _ := ioutil.ReadAll(resp.Body)
  log.Println(string(body))
}
```

Go言語では、データのシーケンシャルなデータの入出力を、`io.Reader`、`io.Writer`インターフェースとして抽象化しています。
ファイル、ソケットなど様々なところでしようされています。

```go
  log.Println("Status:", resp.Status) // 文字列で"200 OK"
  log.Println("StatusCode:", resp.StatusCode) // 数値で200
  log.Println("Headers:", resp.Header) // resp.Headerにはヘッダーの一覧が格納されている。ヘッダーは文字列配列のmap型
  log.Println("Content-Length:", resp.Header.Get("Content-Length"))
```

```sh
GET / HTTP/1.1
Host: localhost:18888
Accept-Encoding: gzip
User-Agent: Go-http-client/1.1
```


<a id="3-5"></a>
## 3.5 GETメソッド+クエリーの送信

GETメソッドでクエリーを送信する方法

```sh
$ curl -G --data-urlencode "query=hello world" http://localhost:18888
# クエリの中にスペースや、URLとして使えない文字がなければ`--data-urlencode`の代わりに、`--data`もしくは短縮系の`-d`も使えます。
# `-G`は`--get`の短縮形です。
```

```go
package main

import (
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
)

func main() {
  values := url.Values{ // クエリーの文字列を作成する、クエリー文字列はurl.Values型を使って宣言します。
    "query": {"hello world"},
  }

  resp, _ := http.Get("http://localhost:18888"+"?"+ values.Encode()) // values.Encode()を呼んで文字列にします。文字列のエスケープもこの関数が行う。
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  log.Println(string(body))
}
```

#### リクエスト

```sh
GET /?query=hello+world HTTP/1.1
Host: localhost:18888
Accept-Encoding: gzip
User-Agent: Go-http-client/1.1
```

<a id="3-6"></a>
## 3.6 HEADメソッドでヘッダーを取得

```sh
$ curl --head http://localhost:18888
```

```go
package main

import (
  "log"
  "net/http"
)

func main() {
  resp, err := http.Head("http://localhost:18888")
  if err != nil {
    panic(err)
  }
  log.Println("Status:", resp.Status)
  log.Println("Headers:", resp.Header)
}
```

#### リクエスト

```sh
HEAD / HTTP/1.1
Host: localhost:18888
User-Agent: Go-http-client/1.1
```


<a id="3-7"></a>
## 3.7 x-www-form-urlencoded形式のPOSTメソッドの送信

```sh
// -d または --data
$ curl -d test=value http://localhost:18888
```

```go
package main

import (
  "log"
  "net/http"
  "net/url"
)

func main() {
  values := url.Values{
    "test": {"value"},
  }
  resp, err := http.PostForm("http://localhost:18888", values)
  if err != nil {
    panic(err) // 送信失敗
  }
  log.Println("Status:", resp.Status)
}
```

#### リクエスト

```sh
POST / HTTP/1.1
Host: localhost:18888
Accept-Encoding: gzip
Content-Length: 10
Content-Type: application/x-www-form-urlencoded
User-Agent: Go-http-client/1.1

test=value
```

<a id="3-8"></a>
## 3.8 POSTメソッドで任意のボディを送信

Postメソッドを使うと任意のコンテンツをボディに入れて送信できます。<br>
HTTP/1.0のブラウザを使って送信することはできませんでしたが、HTTP/1.1 以降に登場したXMLHttpRequestを使って実現できるようになります。

```sh
$ curl -T main.go -H "Content-Type: text/plain" http://localhost:18888
```

```go
package main

import (
  "log"
  "net/http"
  "os"
)

func main() {
  file, err := os.Open("main.go")
  if err != nil {
    panic(err)
  }
  // os.Open()関数で作成されるos.Fileオブジェクトはio.Readerインターフェースを満たしているためそのままhttp.Post()に渡すことができる。
  resp, err := http.Post("http://localhost:18888", "text/plain", file)
  if err != nil {
    // 送信失敗
    panic(err)
  }
  log.Println("Status:", resp.Status)
}
```


#### リクエスト
```sh
POST / HTTP/1.1
Host: localhost:18888
Transfer-Encoding: chunked
Accept-Encoding: gzip
Content-Type: text/plain
User-Agent: Go-http-client/1.1

133
package main

import (
  "log"
  "net/http"
  "os"
)

func main() {
  file, err := os.Open("main.go")
  if err != nil {
    panic(err)
  }
  resp, err := http.Post("http://localhost:18888", "text/plain", file)
  if err != nil {
    // 送信失敗
    panic(err)
  }
  log.Println("Status:", resp.Status)
}

0
```

`Content-Type`ヘッダーの内容は`http.Post()`メソッドの第二引数に指定する。
`io.Reader` の形式で渡す。`os.Open()`関数は`os.File`オブジェクトは`io.Reader` インターフェースを満たしているため、そのまま`http.Post`に渡すことができる。


ファイルではなく、プログラム中で生成したテキストをhttp.Postに渡す場合は、`bytes.Buffer`, `strings.Reader`を使って`io.Reader`インターフェース化する。

```go
func main() {
  reader := strings.NewReader("テキスト")
  resp, err = http.Post("http://localhost:18888", "text/plain", reader)
}
```

<a id="3-9"></a>
## 3.9 multipart/form-data形式でファイルの送信

```
$ curl -F "name=Micheal Jackson" -F "thumbnail=@photo.jpg" http://localhost:18888
```

テキストデータと画像ファイルの2つのデータを送信しています。
```go
package main

import (
  "bytes"
  "io"
  "log"
  "mime/multipart"
  "net/http"
  "os"
)

func main() {
  var buffer bytes.Buffer // 格納するバッファを用意
  writer := multipart.NewWriter(&buffer)
  writer.WriteField("name", "Michael Jackson") // テキストデータ用フィールドはWriteField()メソッドを使って登録

  // ここから 画像ファイル読み込み操作 ここから
  fileWriter, err := writer.CreateFormFile("thumbnail", "photo.jpg")
  if err != nil {
    panic(err)
  }
  readFile, err := os.Open("photo.jpg") // 画像ファイルを開く
  if err != nil {
    panic(err)
  }
  defer readFile.Close()
  // io.Copy()を使って、ファイルの全コンテンツを、ファイル書き込み用のio.Writerにコピー
  io.Copy(fileWriter, readFile)
  // ここまで 画像ファイル読み込み操作 ここまで

  writer.Close() // 最後にマルチパートのio.Writerをクローズし、バッファに全てを書き込みます。

  resp, err := http.Post("http://localhost:18888", writer.FormDataContentType(), &buffer)
  writer.FormDataContentType() は "multipart/form-data; boundary=" + wwriter.Boundary() と同義
  if err != nil {
    panic(err)
  }
  log.Println("Status:", resp.Status)
}

```

### 3.9.1 送信するファイルに任意のMIMEタイプを設定する

前節のコードでファイルの送信はできるようになりました。しかし、
各ファイルのContent-Typeは事実上`void`型とも言える`application/octet-stream`型となってしまいます。
`textproto.MIMEHeader`を使うことで任意のMIMEタイプを設定できます。

```go
// package main

// import (
//   "bytes"
//   "io"
//   "log"
//   "mime/multipart"
//   "net/http"
//   "os"
   "net/textproto"
// )

// func main() {
//   var buffer bytes.Buffer
//   writer := multipart.NewWriter(&buffer)
//   writer.WriteField("name", "Michael Jackson")

  part := make(textproto.MIMEHeader)
  part.Set("Content-Type", "image/jpeg") // MIMEタイプにimage/jpegを設定
  part.Set("Content-Disposition", `form-data; name="thumbnail"; filename="photo.jpg"`)
  fileWriter, err := writer.CreatePart(part)
//   if err != nil {
//     panic(err)
//   }
//   readFile, err := os.Open("photo.jpg")
//   if err != nil {
//     panic(err)
//   }
//   io.Copy(fileWriter, readFile)

//   writer.Close()

//   resp, err := http.Post("http://localhost:18888", writer.FormDataContentType(), &buffer)
//   if err != nil {
//     panic(err)
//   }
//   log.Println("Status:", resp.Status)
// }
```

<a id="3-10"></a>
## 3.10 クッキーの送受信

```go
package main

import (
  "log"
  "net/http"
  "net/http/cookiejar"
  "net/http/httputil"
)

func main() {
  jar, err := cookiejar.New(nil) // クッキーを保存するcookiejarにインスタンスを作成します
  if err != nil {
    panic(err)
  }
  client := http.Client{ // クッキーを保存可能なhttp/Clientインスタンスを作成します
    Jar: jar,
  }
  for i := 0; i < 2; i++ {
    // クッキーは初回アクセスでクッキーを受信し、2回目以降のアクセスでクッキーをサーバーに対して送信する仕組みなので、2回アクセスしています
    // http.Get()の代わりに、作成したクライアントのGet()メソッドを使ってアクセスします。
    resp, err := client.Get("http://localhost:18888/cookie")
    if err != nil {
      panic(err)
    }
    dump, err := httputil.DumpResponse(resp, true)
    if err != nil {
      panic(err)
    }
    log.Println(string(dump))
  }
}
```

```sh
GET /cookie HTTP/1.1
Host: localhost:18888
Accept-Encoding: gzip
User-Agent: Go-http-client/1.1
```

`http.Client`構造体には、httpパッケージと同様に`Get()`, `Head()`, `Post()`, `PostForm()`が実装されており、
これまで紹介してきた使い方をそのまま踏襲して利用できます。
httpパッケージの関数のほとんどは、内部ではデフォルトの`http.Client`構造体のインスタンスの各メソッドへのエイリアスになっています。


#### go言語のオブジェクトの作成

3種類の文法＋1種類の作法。

```go
// 初期値を与えて作成。（スタンダード）
a := Struct{
  Member: "Value",
}

// new 関数で初期化
a: = new(Struct)

// make 関数で初期化
// 配列のスライス、map、チャネル専用
a := make(map[string]string)

// ライブラリが用意したビルダー関数で作成
// 内部では上記のどれかを利用
a := library.New()
```


## 3.11 プロキシの利用

プロキシを利用するときに使うのは`Transport`メンバー変数

```sh
$ curl -x http://localhost:18888 http://github.com
```

```go
package main

import (
  "log"
  "net/http"
  "net/http/httputil"
  "net/url"
)

func main() {
  proxyUrl, err := url.Parse("http://localhost:18888")
  if err != nil {
    panic(err)
  }
  client := http.Client{ // プロキシを利用するときも、http.Clientを使います。
    Transport: &http.Transport{
      Proxy: http.ProxyURL(proxyUrl),
    },
  }
  resp, err := client.Get("http://github.com")
    // client.Get先は外部サイトになっていますが、プロキシの向き先はローカルのテストサーバ。
    // 外部向けには直接リクエストは飛ばずに、ローカルのサーバが一旦リクエストをうけます。
    // このローカルサーバが直接レスポンスを返しているので、このコードではgithubへのアクセスは発生しない。
  if err != nil {
    panic(err)
  }
  dump, err := httputil.DumpResponse(resp, true)
  if err != nil {
    panic(err)
  }
  log.Println(string(dump))
}

```

デフォルトのhttp.Clientで使われるプロキシパラメータは、環境変数から情報を取得してきてプロキシを設定する処理になっている。<br>
環境変数`HTTP_PROXY`,`HTTPS_PROXY`が設定されている場合は、そちらに設定したプロキシにリクエストを送信する。<br>
また、`NO_PROXY`に設定を無視するホスト名を書いておくことで、設定したホストに対してはプロキシを通さず直接通信できます。


## 3.12 ファイルシステムへのアクセス
`http.Transport`にはスキーマ用のトランスポートを追加する`RegisterProtocol`メソッドがあります。
このメソッドに登録できる、ファイルアクセス用のバックエンド`http.NewFileTransport()`もあります。

```sh
$ curl file://path/to/file
```

```go
package main

import (
  "log"
  "net/http"
  "net/http/httputil"
)

func main() {
  transport := &http.Transport{}
  transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
  client := http.Client{
    Transport: transport,
  }
  resp, err := client.Get("file:///Users/sekitatsuhiko/go/server.go")
  if err != nil {
    panic(err)
  }
  dump, err := httputil.DumpResponse(resp, true)
  if err != nil {
    panic(err)
  }
  log.Println(string(dump))
}
```

## 3.13 自由なメソッドの送信

http.Client構造体のメソッドがサポートしているのは、GET, HEAD, POSTだけ。
これ以外のメソッドのリクエストを行うときは`http.Request`構造体のオブジェクトを使う。

```sh
$ curl -X DELETE http://localhost:18888
```

```go
package main

import (
  "log"
  "net/http"
  "net/http/httputil"
)

func main() {
  client := &http.Client{}
  // http.request構造体はhttp.NewRequest()というビルダー関数を使って生成する。
  // 関数の引数はメソッド、URL、ボディ。
  request, err := http.NewRequest("DELETE", "http://localhost:18888", nil)
  if err != nil {
    panic(err)
  }
  resp, err := client.Do(request)
  if err != nil {
    panic(err)
  }
  dump, err := httputil.DumpResponse(resp, true)
  if err != nil {
    panic(err)
  }
  log.Println(string(dump))
}
```


## 3.14 ヘッダーの送信

```
$ curl -H "Content-Type=image/jpeg" -d "@image.jpeg" http://localhost:18888
```

http.Request構造体はHeader フィールドをもっている。これはhttp.Response構造体のHeaderフィールドと同じものです。
追加するにはAdd()メソッドを使う。

```
// ヘッダーの追加
request.Header.Add("Content-Type", "image/jpeg")
```

curlコマンドのBASIC認証用のオプション、クッキー用オプションは、Go言語でも
`http.Request`構造体のメソッドとして提供されている。

```
$ curl --basic -u ユーザー名:パスワード

// BASIC認証
request.SetBasicAuth("ユーザー名", "パスワード")
```

```
$ curl -c ファイル

// クッキーを手動で一つ足す
request.AddCookie(&http.Cookie{Name:"test", Value:"value"})
```

## 3.15 国際化ドメイン

URLの国際化をGo言語でも変換することができます。
変換は`idna.ToASCII()`と`idna.ToUnicode()`関数で行います。
リクエスト前にドメイン名を`idna.ToASCII()`で変換することで日本語ドメインのサイトの情報を取得できます。

事前にpackage を読み込んでおく

```
// 事前にシェルの設定ファイルに $GOPATHを設定する必要あり
// -v は verbose オプション
$ go get -v golang.org/x/net/idna
```

```go
package main

import (
  "fmt"
  "golang.org/x/net/idna"
)

func main() {
  src := "握力王"
  ascii, err := idna.ToASCII(src)
  if err != nil {
    panic(err)
  }
  fmt.Printf("%s -> %s\n", src, ascii)
}
```


## 3.16 本章のまとめ

__goに入ればgoに従え__
