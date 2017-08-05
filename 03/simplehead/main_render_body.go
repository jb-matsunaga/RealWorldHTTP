package main

import (
  "log"
  "io/ioutil"
  "net/http"
)

func main() {
  resp, err := http.Head("http://localhost:18888")
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  log.Println(string(body))
  log.Println("Status:", resp.Status)
  log.Println("Headers:", resp.Header)
}

