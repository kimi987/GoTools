package main

import (
	"os"
	"path"
	"fmt"
	"net/http"
)

func main() {
	root := "temp/"
	os.MkdirAll(path.Dir(root), os.ModePerm)

	http.HandleFunc("/download", downloadFunc)

	// 回放服务器
	prefix := "/replay/"
	http.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(root))))

	// 监控服务器，promethus metrics

	http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
}

func downloadFunc(w http.ResponseWriter, r *http.Request){
	w.Write([]byte(downloadHtml))
}

const downloadHtml = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>下载安装</title>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
<style type="text/css">
.jumbotron {
  text-align: center;
  border-bottom: 1px solid #e5e5e5;
}
.jumbotron .btn {
  padding: 14px 24px;
  font-size: 21px;
}

</style>
</head>
<body>
    <div class="container">
      <div class="jumbotron">
        <h1>七雄争霸3D</h1>
        <p class="lead">当你看到这个页面的时候</p>
        <p class="lead">不要犹豫，直接点我</p>
        <p><a class="btn btn-lg btn-success" href="/replay/male7_last.apk" role="button">点我</a></p>
        <h6>注：微信用户请在右上角选择“在浏览器打开”，再点我</h6>
      </div>
    </div>
</body>
</html>
`