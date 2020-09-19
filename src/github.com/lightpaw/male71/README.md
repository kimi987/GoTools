### 七雄争霸服务器 [![Build Status](https://travis-ci.com/lightpaw/male7.svg?token=XwRhaKjJPghm9pUZqVRu&branch=master)](https://travis-ci.com/lightpaw/male7)

必须在同目录有`server.yaml`配置文件和`conf`配置文件夹.

新同学一上来需要安装gogen和luagen开发工具

首先clone项目generator

    git clone git@github.com:lightpaw/generator.git

安装 gogen（服务器代码生成器）

编译

    cd cmd/gogen
    go install

安装proto3编译环境（proto 3.2.0） https://github.com/google/protobuf

安装gogo-protobuf

    go get github.com/lightpaw/protobuf/protoc-gen-gofast

测试，在male7的根路径下执行

    cd male7
    gogen
    
安装 luagen（客户端代码生成器）

编译

    cd cmd/luagen
    go install
    
安装protobuf的lua编译环境（待补充）

---

[服务器配置说明](/vendor/github.com/lightpaw/config/README.md)