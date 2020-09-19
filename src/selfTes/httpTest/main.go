package main

import(
	"net/http"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	http.HandleFunc("/", IndexHandler)
	// http.HandleFunc("/validate", ValidateHandler)
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir("D:/goPath/src/selfTes/httpTest/download"))))

	http.ListenAndServe("0.0.0.0:5555", nil)
	
}

 
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
	logrus.Debug("abc")
}