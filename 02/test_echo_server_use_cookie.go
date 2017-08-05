package main

import(
  "fmt"
  "log"
  "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("Set-Cookie", "VISIT=TRUE")
  if _, ok := r.Header["Cookie"]; ok {
    // 一度来たことがある人
    fmt.Fprintf(w, "<html><body>2回目以降</body></html>\n")
  } else {
    // はじめて来たことがある人
    fmt.Fprintf(w, "<html><body>初訪問</body></html>\n")
  }
}

func main() {
  var httpServer http.Server
  http.HandleFunc("/", handler)
  log.Println("start http listening :18888")
  httpServer.Addr = ":18888"
  log.Println(httpServer.ListenAndServe())
}
