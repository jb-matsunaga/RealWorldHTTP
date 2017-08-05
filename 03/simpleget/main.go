package main // ライブラリ以外のソースコードは必ずpackage mainから始まる

import ( // 必要なパッケージを取り込みます。Go言語ではそのファイルで宣言したパッケージ以外は使用できません。
  "io/ioutil"
  "log"
  "net/http"
)

func main() { // 全てのプログラムはmainパッケージのmain関数が最初に呼ばれます
  resp, err := http.Get("http://localhost:18888")
  if err != nil { // Go言語のエラー処理コードです。Go言語の関数は返り値としてエラーを返すので、それがnilかどうか確認します。例外処理はありません。
    panic(err) // panicはエラーを表示させてプログラムを終了させる
  }
  defer resp.Body.Close() // これは後処理コードです。defer を付与するとこの関数から抜けた後にこのぶんを実行します。ソケットからボディを読み込んだ後の処理です。
  body, err := ioutil.ReadAll(resp.Body) // ボディの内容をバイト列として読み込みます
  if err != nil {
    panic(err)
  }
  log.Println(string(body))

}

