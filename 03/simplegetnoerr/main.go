package main

import (
  "io/ioutil"
  "log"
  "net/http"
)

func main() {
  resp, _ := http.Get("http://localhost:18888") // resp変数に入っている要素が、http.Response型のオブジェクト。サーバから返ってきたさまざまな情報を全て格納している。
  defer resp.Body.Close() // ボディはbodyメンバ変数に格納されている。
  defer log.Println("defer")
  body, _ := ioutil.ReadAll(resp.Body)
  log.Println(string(body))
  log.Println("Status:", resp.Status) // 文字列で"200 OK"
  log.Println("StatusCode:", resp.StatusCode) // 数値で200
  log.Println("Headers:", resp.Header) // resp.Headerにはヘッダーの一覧が格納されている。ヘッダーは文字列配列のmap型
  log.Println("Content-Length:", resp.Header.Get("Content-Length"))
}

