# 4.HTTP/1.1のシンタックス：高速化と安全性を求めた拡張

本章で取り上げる、プロトコルシンタックスとしてのHTTP/1.1の変更点は次の通り。
 - 通信の高速化
    - Keep-Aliveがデフォルトで有効に
    - パイプライニング
- TLSによる暗号化通信のサポート
- 新メソッドの追加
    - PUTとDELETEが必須のメソッドとなった
    - OPTION, TRACE, CONNECT メソッドが追加
- プロトコルのアップグレード
- 名前を使ったバーチャルホストのサポート
- サイズが事前にわからないコンテンツのチャンク転送エンコーディングのサポート

## 4.1 通信の高速化
キャッシュはコンテンツのリソースごとに通信を最適化する技術だが、Keep-Aliveと、パイプライニングはより汎用的な、すべてのHTTP通信を高速化する機能である。

### 4.1.1 Keep-Alive
Keep-Aliveは、HTTPの下のレイヤーであるTCP/IPの通信を効率化する仕組み。

Keep-Aliveを使わない場合は、ひとつのリクエストごとに通信を閉じるが、Keep-Aliveを使うと連続したリクエストの時に接続を再利用する。TCP/IPは接続までの待ち時間が減り、通信速度が上がったように感じる。

![](./img/keep-alive.png)
1回のTCP接続で複数のHTTPリクエストを処理する機能

Keep-Aliveによる通信は、クライアント、サーバーのどちらかが次のヘッダーを付与して接続を切るか、タイムアウトするまで接続が維持される。

```
Connection: Close
```

Keep-Aliveの接続時間は、クライアントとサーバーの両方が持っている。片方がTCP/IPの切断を行った瞬間に通信が完了するため、どちらか短い方が採用される。



>*参照*

>[持続接続](https://docs.oracle.com/javase/jp/6/technotes/guides/net/http-keepalive.html)

>[TCP/IP通信とは](http://research.nii.ac.jp/~ichiro/syspro98/tcpip.html)

>[Nginx の keep-alive の設定と検証](http://www.nari64.com/?p=579)





