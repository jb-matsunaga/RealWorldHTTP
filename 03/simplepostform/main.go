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
