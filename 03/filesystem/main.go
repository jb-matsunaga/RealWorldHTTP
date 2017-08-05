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
  resp, err := client.Get("file:///Users/sekitatsuhiko/air-design/Real-World-HTTP/03/server/server.go")
  if err != nil {
    panic(err)
  }
  dump, err := httputil.DumpResponse(resp, true)
  if err != nil {
    panic(err)
  }
  log.Println(string(dump))
}
