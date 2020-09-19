package main

import(
	"net/http"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	// http.HandleFunc("/", IndexHandler)
	// http.HandleFunc("/validate", ValidateHandler)
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir("D:/goPath/src/selfTes/httpTest/download"))))

	err := http.ListenAndServe("0.0.0.0:5555", nil) //设置监听的端口
	if err != nil {
	   logrus.Fatal("ListenAndServe: ", err)
	}
	
}

 
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world 1111")
	logrus.Debug("abc")
}