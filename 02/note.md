# HTTP/1.0のセマンティクス：ブラウザの基本機能の裏側

1章ではHTTPの基本機能が紹介された
- メソッドとパス
- ヘッダー
- ボディー
- ステータスコード

Webの高度化にともない多くの機能が追加されているが、
特にヘッダーで多くの機能が実現されている.
ブラウザがこれらの基本要素をどのように応用して
基本機能を実現しているのかをこの章では見ていきます.

## フォーム送信について
- シンプルなフォームの送信(x-www-form-urlencoded)
- フォームを使ったファイルの送信(multipart/form-data)

### シンプルなフォームの送信(x-www-form-urlencoded)
フォームを使ったPOST送信にはいくつかあるが、一番シンプルなものから紹介.

```html
<form method="post">
  <input type="text" name="title">
  <input type="text" name="atuhor">
  <input type="submit">
</form>
```

```sh
$ curl -v --http1.0 -d title="Books Title" -d author="Books Author" http://localhost:18888
```

例: https://beauty.hotpepper.jp/catalog/ladys/

検索結果の並び順を変えるとフォームを使ってPOST送信され、
`Content-Type:application/x-www-form-urlencoded` が設定される.

#### urlエンコードについて
ブラウザは `RFC1866` で定める変換フォーマットに従って変換をおこなうため、
以下の項目以外はエスケープが必要になる.
- アルファベット
- 数値
- アスタリスク
- ハイフン
- ピリオド
- アンダースコア

#### GET送信の場合
POSTと違い、GETの場合はボディではなく、クエリーとしてURLに付与される.

### フォームを使ったファイルの送信(multipart/form-data)
フォームを使ったPOST送信にはいくつかあるが、一番シンプルなものから紹介.

```html
<form method="post" enctype="multipart/type">
  <input type="text" name="title">
  <input type="text" name="atuhor">
  <input type="file" name="sample">
  <input type="submit">
</form>
```

```sh
$ curl -v --http1.0 -F title="TITLE" -F auhtor="AUTHOR" -F attachment-file@test.txt  http://localhost:18888
```

参考: <br>
フォームによるファイルアップロードの仕様 <br> https://www.javadrive.jp/servlet/fileupload_tutorial/index2.html

## 300番台のリダイレクトの懸念事項
リダイレクトには懸念事項もある.

1. URLは2000文字を目安にするべき（GETのクエリーには送信データ量に制限がある）
2. データがURLに入るため送信した内容がアクセスログに残りセキュリティ的に懸念がある

上記回避策としてフォームを利用したリダイレクトがある.

自動リダイレクトするフォーム
```html
<hrml>
<body action="next" onload="document.forms[0].submit()">
  <form method="post">
    <input type="hidden" name="id">
    <input type="hidden" name="password">
    ....
    <input type="submit">
  </form>
</body>
</hrml>
```

## コンテントネゴシエーション
サーバーとクライアントは別々に開発されているため、両者が期待する形式や設定が常に一致しているとは限らない.
そのため、1リクエストの中で両者がベストの設定を共有する仕組みが必要となり、それをコンテントネゴシエーションという.
コンテントネゴシエーションはヘッダーを利用する.ネゴシエーションする対象とヘッダーは以下の通りとなっている.

| リクエスト | レスポンス | ネゴシエーション対象 |
|:-----------|:------------|:------------|
| Accept | Content-Type | MIMEタイプ |
| Accept-Language | Content-Language | 表示言語 |
| Accept-Charset | Content-Type | 文字のキャラクターセット |
| Accept-Encoding | Content-Encoding | ボディーの圧縮 |

### Accept: ファイル種類の決定
ChromeのAcceptは以下の通り.
`accept:text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8`

わかりやすい画像を例に
image/webp<br>
image/apng<br>
*/*;q=0.8<br>

`q` は品質係数と呼ばれ、0〜1で設定し、デフォルトは1となっている.
これは優先度の高さを表し、サーバーは優先度が高く、対応しているフォーマットから返す.
もしお互いに一致しているものがなかったら、406 Not Acceptableエラーとなる.

補足:
- webp -> Googleが推奨するpngよりも2割程ファイルサイズが小さい画像フォーマット.
- apng -> アニメーションするpng. gifに取って代わる次世代の新しい画像フォーマット.

### Accept-Language: 表示言語の決定
GitHubのAccept-Languageは以下の通り.<br>
`Accept-Language:ja,en-US;q=0.8,en;q=0.6`

ja,en-US,enの順でリクエストを送る.
言語情報を収める箱として`Content-Language`があるが定義されているが
多くのWebサイトではそれを使っていなく、次のようにhtmlタグの中で返しているページが多いらしい.
GitHubは`ja`の日本語設定高くしているのって、時期に日本語対応してくれるのかなと思った.

```html
<html lang="en">
```

### Accept-Charset: キャラクターセットの決定
ヘッダーはご想像どおり以下のような感じ.<br>
`Accept-Charset: UTF-8,SHIFT_JIS;q=0.7,*;q=0.3`

全キャラクターセットのエンコーダをブラウザが内包しているため、
モダンブラウザはAccept-Charsetを送信していない.
なので、事前にネゴシエーションする必要がなくなったみたい.

HTMLの<meta http-equiv>はHTTPヘッダーと同じ指示をドキュメントの内部に埋め込んで返すための箱.

```html
<meta charset="UTF-8">
```

使用できるキャラクターセットはIANAで管理されている.<br>
著者が一番言いたげなのはUTF-8, SHIFT_JISなど区切り文字を
なぜ統一しなかったんだ...ということのように思えてきた.

## Accept-Encoding: ボディーの圧縮
ボディーの圧縮は圧縮による通信速度の向上のためにも積極的にやってください.
圧縮してそれを展開してっていう処理を入れても、圧縮しない時よりも、Webページ表示までにかかるトータルの処理速度は早くなる.

- ブラウザ利用者にとって転送速度がはやくなる
- データ量が減ることで通信料金も安くなる

リクエスト
`accept-encoding:gzip, deflate, br`

レスポンス
`content-encoding:gzip`

```sh
$ curl --http1.0 --compressed http://localhost:18888
```

補足:
- Brotli(br) -> gzipよりも効率が良い新しい圧縮フォーマット

## クッキー
Webサイトの情報をブラウザ側に保存する仕組み.<br>
クッキーもHTTPのヘッダーをインフラとして実装されている.

```sh
$ curl -v --http1.0 -c cookie.txt -b cookie.txt http://localhost:18888

-c: 指定したファイルにcookieを保存
-b: 指定したファイルと指定したkey:valueをcookieとして送信
```

### クッキーの間違った使い方
#### 永続性の問題
シークレットモード、改ざん、etc..、確実に保持されるものではないので、
なくなっても困らないもの、サーバーの情報から復元できるもの以外は保存しないこと.

#### 容量の問題
容量制限がある(4Kバイトほど)、リクエスト/レスポンスのヘッダーに入るため速度が劣化するなどの問題もある.

#### セキュリティの問題
secureオプションを付与すればHTTPSでしか送受信されないが、HTTPの場合は平文で送受信される.
パスワードや個人情報など見えて困るものを入れると情報漏洩の問題となるし、書き換えも容易なので誤動作につながるような情報は入れないようにすること.

### クッキーに制約をあたえる
クッキーは特定のサービスを利用するためのトークンとして利用することが多いため、
クッキーを本来必要としないサーバに送信することはセキュリティリスクを高めることにつながる.
そのため、送信先の制御や寿命を設定するような属性が存在し、HTTPクライアントはこれらの属性を解釈し、クッキーの送信を制御する責務がある.

#### Expires, Max-Age属性
クッキーの寿命を設定する属性.

#### Domain属性
クライアントからクッキーを送信する対象のサーバー.
省略時はクッキーを発行したサーバー.

#### Path属性
クライアントからクッキーを送信する対象のサーバーのパス.
省略時はクッキーを発行したサーバーのパス.

#### Secure属性
HTTPS接続での安全な接続以外は、クライアントからサーバーへのクッキー送信をしない.

#### HttpOnly属性
CookieをJavaScriptから触ることができないようにする属性.
クロスサイトスクリプティングなど、悪意のあるJavaScriptが実行されるリスクを守る.

#### SameSite属性
RFCには存在しない属性.
リクエストの起因となったオリジンと、リクエスト先のオリジンが異なるようなリクエストには送信しないように設定する属性.<br>
参考: http://qiita.com/flano_yuki/items/b87b2c28db0b056665ef#same-site-cookies

## 認証とセッション
ユーザ名とパスワードを入力してログインすることを認証といい、
サービス側はその認証情報をもとに、誰がアクセスしているかを特定する.

- Basic認証
- Digest認証
- クッキーを使ったセッション管理

### Basic認証
ユーザ名とパスワードをbase64エンコーディングしたもの.
SSL/TLSを使っていないと、通信を傍受されたら簡単にユーザ名とパスワードが漏洩する.

```sh
$ curl -v --http1.0 --basic -u user:pass http://localhost:18888
```

ヘッダーには`Authorization: Basic dXNlcjpwYXNz`が付与される.

### Digest認証
Basic認証より強固なものが、ハッシュ関数（A->Bは簡単だが、B->Aは簡単には計算できない）を利用したDigest認証.
認証対象画面へアクセスすると401認証エラーを返し、承認ダイアログを表示する.

`www-Authenticate: Digest username="ユーザ名" realm="エリア名", nonce="0123456789", algorithm=MD5, qop="auth"`

### クッキーを使ったセッション管理
#### Basic・Digest認証は現在はあまり使われていない理由
- 特定フォルダ配下を認証しないと見せないという使い方しかできない.<br>
  ログインしてもしなくても良い画面とかできない
- ログイン画面がカスタムできない.<br>
  画像表示: https://point.recruit.co.jp/member/login/
- 明示的にログアウトできない.
- ログインした端末の識別ができない.<br>
  Googleとかだと新しい端末からアクセスしたら警告メールがきますよね

#### クッキーとセッションを使った認証管理が主流
認証の流れとしては以下の通り
1. クライアントは、フォームでユーザ名とパスワードを送信<br>
   ユーザ名とパスワードは直接送信するのでSSL/TLSでないと簡単に漏洩するので注意
2. サーバ側ではユーザ名とパスワードで認証し、問題がなければセッショントークンを発行し、DBに保存
3. トークンはクッキーとしてクライアントに返却
こっちの方がわかりやすい:<br> http://qiita.com/hththt/items/07136ad74127999df271

## プロキシ
プロキシはHTTPなどの通信を中継する仕組み.

### プロキシサーバ導入のメリット
1. 送信元を隠す<br>
  Proxy を介して通信する事で、WEBページにとってはProxy が送信元として捉えられる.
  この為、送信元を隠す目的で使用される.
  例えば、アメリカのProxyを使えば、アメリカからアクセスしているように思わせる事も可能.
2. セキュリティ<br>
  大企業などでは、社内からInternet を閲覧する際に、必ずと言って良いほど、Proxy を使って通信をさせてる.
  コレは、Proxy を中継させる事で、社員のアクセス管理を行う事が可能となり、不要なサイトへのアクセスを禁じたり、アクセスの統計を取得したりする目的で使われている.
3. キャッシュサーバ<br>
  Proxy の役割は、『代理サーバ』とは別に『キャッシュサーバ』としての役割を持つ.
  通常、クライアントのブラウザなどにもこの『キャッシュ』機能は搭載されているが、その巨大版です。

```sh
$ curl -v --http1.0 -x http://localhost:18888/helloworld/ -U user:pass http://example.com/helloworld/
```
> GET http://example.com/helloworld/ HTTP/1.0　/　Host: example.com

```sh
$ curl -v --http1.0 http://localhost:18888/helloworld/
```
> GET /helloworld/ HTTP/1.0　/　Host: localhost:18888

## キャッシュ
Webサイトのリッチ化に伴い、読み込みファイル数・サイズが増加傾向にある.
通信回線が速くなっても、毎回読み込んでいては表示速度が遅くなってしまう.
キャッシュは、ダウンロード済みで内容に変化がなければ、新たに読み込むのではなく、
ダウンロード済みのものを表示してパフォーマンスをあげる仕組みである.

### 更新日時によるキャッシュ
HTTP1.0のキャッシュの仕組み.
当時は静的コンテンツがメインだったので、コンテンツの新旧を比較するだけで事足りていた.

#### 仕組み
サーバーは、コンテンツの最終更新時刻をレスポンスヘッダの`Last-Modified`に入れて送信する.
ブラウザ側はこのコンテンツの最終更新時刻を覚えておき、
次回リクエストした際にリクエストヘッダの`If-Modified-Since`の中に含めて送信する.
サーバー側のコンテンツに変更がなければサーバーはステータスコード304 コンテンツ未更新ステータスコードを送り
ブラウザはキャッシュされたコンテンツを表示させるという仕組み.

### Expiresによるキャッシュ
更新日時によるキャッシュの場合、キャッシュの有効性を確認するためにどうしても通信が発生する.
その通信自体をなくす仕組みがHTTP1.0に追加されたExpiresによるキャッシュ.

#### 仕組み
Expiresはレスポンスヘッダのひとつで、新しいファイルが存在するかどうかを確認することなく、Expiresの期限内であればブラウザでキャッシュ済みのファイルを強制的に適用する.

コンテンツが日々変わるページには不向きで、あまり変更の入らないcssやjsファイルとかによく使われる.<br>
https://beauty.hotpepper.jp/svcSA/

### Pragma: no-cache
クライアントからプロキシに対して指示を送ることもある.
実装依存の命令を含むリクエストヘッダーの置き場として、
HTTP1.0からPragmaヘッダーが定義された(no-cacheしか設定できない).
リクエストしたコンテンツがキャッシュされていたとしても、
本来のサーバーまでリクエストを届けるようにするための仕組み.

HTTP/1.1ではCache-Controlにマージされました.

### ETagの追加
同じ画面にアクセスしてもユーザの状態によって、コンテンツの内容が異なる場合がある.
例えば、ログイン前後や、会員ステータス（プレミアム会員、会員、会員ではない）でコンテンツを出し分ける場合など.
このように動的に変更する要素が増えれば増えるほど、どの日時を根拠にキャッシュの有効性を判断すれば良いのか判断が難しくなる.
その場合に使用できるのがHTTP/1.1で定義されたETagヘッダー.

#### 仕組み
ETagはレスポンスヘッダのひとつで、
初回アクセス時に、サーバーはレスポンスにETagヘッダーにハッシュ値を付与する.
ブラウザは2度目以降のアクセス時にIf-None-Matchヘッダーにダウンロード済みのETagの値をつけてリクエストする.
サーバーはIf-None-Matchとこれから送りたいファイルのETagとを比較し、同じなら304 コンテンツ未更新を返し、ブラウザはキャッシュされたコンテンツを表示させるという仕組み.

`ETag:"6307c-553dda1253000"`
参考: https://www.ponparemall.com/

### Cache-Control(1)
ETagと同時期にHTTP/1.1で追加されたのが、Cache-Control.
Expiresよりも優先されるキャッシュの仕組み.
まずはサーバーからレスポンスとして送付されるヘッダーについての紹介.

| ディレクティブ | 説明 |
|:-----------|:------------|
| public | 複数ユーザーで共有できるようにキャッシュしてよい(プロキシOK, ブラウザOK) |
| private | 特定ユーザーだけが使えるようにキャッシュしてよい(ブラウザOK) |
| max-age=n | キャッシュの鮮度を秒で設定.86400を指定すると、1日キャッシュが有効でサーバーに問い合わせることなくキャッシュを利用する.それ以降はサーバーに問い合わせを行い304 コンテンツ未更新が返ってきた時のみキャッシュを利用する. |
| s-maxage=n | max-ageと同等だが、共有(プロキシの)キャッシュに対する設定値 |
| no-cache | 一度キャッシュに記録されたコンテンツは、現在でも有効か否かをサーバに問い合わせて確認がとれない限り再利用してはならない、という意味 |
| no-store | キャッシュしない |

参考：
- キャッシュについて整理してみた<br>
  http://qiita.com/karore/items/2dc6ab8347c940ea4648#pencil2-httpレスポンスヘッダ-で制御
- max-ageとs-maxageの違いについて<br>
  https://suin.io/534
- セキュリティに絡むとこ...<br>
  http://tech.mercari.com/entry/2017/06/22/204500

### Cache-Control(2)
#### リクエスト時にプロキシに対してキャッシュに関する指示をだす
Cache-Controlヘッダーをリクエストヘッダーに含めることでプロキシへ指示することができる.
まずはクライアント側からリクエストヘッダーで使える設定値を紹介.

| ディレクティブ | 説明 |
|:-----------|:------------|
| no-cache | Pragma: no-cacheと同等 |
| no-store | レスポンスのno-storeと同じで、プロキシサーバにキャッシュを削除するように指示する |
| max-age | プロキシで保存されたキャッシュの有効期限を設定 |
| max-stale | キャッシュされたリソースの有効期限が切れていても指定時間内であればキャッシュを受け入れるよう設定 |
| min-fresh | 指定された時間は新鮮であるものを返すようにプロキシに要求する.逆に、指定時間内で有効期限が切れるリソースはレスポンスとして返せない. |
| no-transform | コンテンツを改変しないようにプロキシに要求する.例えば、画像の圧縮とか |
| only-if-chached | キャッシュからのみデータを取得するという指示で、設定時には初回をのぞいてオリジンサーバーへのアクセスは一切行わない |

#### レスポンス時にプロキシに対してキャッシュに関する指示をだす
レスポンスヘッダーでサーバーがプロキシに対して送信するキャッシュ制御の指示についての紹介.
補足をすると、Cache-Control(1)で紹介したディレクティブもすべてプロキシに対して有効.

| ディレクティブ | 説明 |
|:-----------|:------------|
| no-transform | プロキシがコンテンツを改変するのを抑制 |
| must-revalidate | no-cacheとほぼ同じだが、プロキシへの司令となる.キャッシュが期限切れだった場合、オリジンサーバでの確認無しにキャッシュを利用してはならない |
| proxy-revalidate | must-revalidateと同じだが、共有キャッシュのみに対する要請 |

参考
- HTTPヘッダーフィールド2<br>
  http://d.hatena.ne.jp/s-kita/20080927/1222508924

### Vary
ETagの説明では同じURLでも個人ごとに結果が異なるケースについて紹介された.
同じURLでもクライアントによって返す結果が異なることを示す場合はVaryを使用する.

#### 仕組み
例えば、サーバがクライアントのUser-Agentリクエストヘッダーによって返す内容を変えているとする.この場合、サーバから返されたデータをproxyが素直にキャッシュしてしまうと、別のUser-Agentがproxyにアクセスしたときにサーバが意図しないデータがクライアントに返されてしまう危険がある.

その場合、レスポンスヘッダに`Vary: User-Agent`が指定されていれば、
proxyはUser-Agentによって内容が変わることを知ることができるので、
キャッシュをしないとか、User-Agent毎に異なるキャッシュを保持するといった対応が
可能となる.

## リファラー
ユーザーがどの経路からWebサイトに来たのかをサーバーが把握するために、
クライアントがサーバーに送るヘッダーをリファラーという.

`Referer:https://www.google.co.jp`

注意:<br>
`referrer` ではなく `referer`.<br>
スペルミスのままRFCに定義されちゃったらしい.

### リファラーを制御する(1)
ユーザの通信内容を秘密にするHTTPSがHTTP/1.1から追加された.
保護された通信内容が保護されていない通信経路に漏れるのを防ぐため、
クライアントはリファラーの送信に制御を加えることをRFCで定められている.

| アクセス元 | アクセス先 | 送信するかどうか |
|:-----------|:------------|:------------|
| HTTPS | HTTPS | する |
| HTTPS | HTTP | しない |
| HTTP | HTTPS | する |
| HTTP | HTTP | する |

### リファラーを制御する(2)
リファラーを制御する(1)を厳密に適用すると、サービス間に支障が出ることもあり、
リファラーポリシーなるものが提案され、現在ドラフトステータスとなっている.<br>
困る事例: http://web-tan.forum.impressrd.jp/e/2015/04/14/19750

リファラーポリシーは以下の方法で設定できる.
- Referrer-Policyヘッダー
- &lt;meta name="referrer" content="設定値"&gt;
- aタグなどいくつかの要素のreferrerpolicy属性および、rel="noreferrer"属性
- Content-Security-Policyヘッダー(廃止された模様)

リファラーポリシーとして設定できる値には次のようなものがある
- no-referrer<br>
  一切おくらない
- no-referrer-when-downgrade<br>
  現在のデフォルト動作と同じで、HTTPS→HTTP時は送信しない
- same-origin<br>
  同一ドメイン内のリンクに対してのみ、リファラーを送信
- origin<br>
  詳細ページではなく、トップページからリンクされたものとしてドメイン名だけを送信
- strict-origin<br>
  originと同じだが、HTTPS→HTTP時は送信しない
- origin-when-crossorigin<br>
  同じドメイン内ではフルのリファラーを、別ドメインにはトップのドメイン名だけを送信
- strict-origin-when-crossorigin<br>
  origin-when-crossoriginと同じだが、HTTPS→HTTP時は送信しない
- unsafe-url<br>
  常に送信

リファラーポリシーについてはこの記事がよくまってる<br>
http://qiita.com/wakaba@github/items/707d72f97f2862cd8000

## 検索エンジン向けのコンテンツのアクセス制御
クローラー向けのアクセス制御の方法として、主に2つの手法が広く使われている.
- robots.txt<br>
  robots.txtは、主に検索エンジンの巡回を指示するファイル.
- サイトマップ<br>
  検索エンジンにサイト内のURLや動画の情報を告知するファイル.

<br>参照: https://digital-marketing.jp/seo/sitemap-xml-and-robots-txt/

### robots.txt
- robots.txtファイルの用意
- robots meta tagの設定

robots.txtはクロール最適化のために行うもの.
一方、robots meta tagは一つ一つのページのインデックスを最適化するために行うもの.<br>
参照: https://bazubu.com/robots-txt-16678.html

robots.txtファイルはドメインのルートディレクトリに設置する必要がある.
つまり、サブディレクトリ型のレンタルウェブスペースに設置されているサイトでrobots.txtファイルを使うことはできない.<br>
参照: http://whitehatseo.jp/robots-txtの記述例と使い方を解説します/

#### robots.txtファイルの用意
robots.txtは、サーバーのコンテンツ提供者が、
クローラーに対してアクセスの許可・不許可を伝えるためのプロトコル.
robots.txtは以下のような形式で読み込みを禁止するクローラーの名前と場所を指定する.

```txt
User-agent: *
Disallow: /cgi-bin/
Disallow: /tmp/
```

上記は、全クローラーに対して、/cgi-bin/フォルダと/tmp/フォルダへのアクセスを禁止している例となっている.
`User-agent: Googlebot`のように、特定のクローラーに対しての指定も可能.

`例: https://www.facebook.com/robots.txt`

#### robots meta tagの設定
robots.txtと同じような内容をHTMLのメタタグに記述できる.
robots.txtの方が優先されるが、こちらの方が細かく指定可能となっている.

`<meta name-"robots" content="noindex" />`

content属性にはさまざまなディレクティブが記述できる.
Googlebotが解釈するディレクティブの詳細はGoogleのWebサイトに記載されている.
代表的なものを以下に示す.

| ディレクティブ | 意味 |
|:-----------|:------------|
| noindex | 検索エンジンがインデックスするのを拒否する |
| nofollow | クローラーがこのページ内のリンクを辿るのを拒否する |
| noarchive | ページ内のコンテンツをキャッシュするのを拒否する |

同じディレクティブはHTTPのX-Robots-Tagヘッダーにも記述できる
`X-Robots-Tag: noindex, nofollow`

### サイトマップ
サイトマップはWebサイトに含まれるページ一覧とそのメタデータを提供するXMLファイル.
Flashを使って作られたコンテンツや、JavaScriptを多用して作られた動的ページからのリンクなど、
クローラーの実装によってはページが発見できない場合でもサイトマップによって補完できます.

```xml
<?xml version="1.0" encoding="utf-8">
<urlset xmls="http://www.sitemaps.org/schemas/sitemap/0.9"
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9
    http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd">
  <url>
    <loc>http://example.com</loc>
    <lastmod>2006-11-18</lastmod>
  </url>
</xml>
```

このurlタグを登録したいページ数分作成します.
locタグには絶対URLを指定する.

サイトマップは前述のrobots.txt内にも記載可能。
また、各検索エンジンに対してXMLファイルをアップロードする方法がある。
Googleの場合はSearch Consoleサイトマップツールを使う

`Sitemap: https://beauty.hotpepper.jp/sitemap.xml`
