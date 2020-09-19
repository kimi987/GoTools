package main

import (
	"path/filepath"
	"path"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/tlog/gen"
)

func main() {
	dir, err := filepath.Abs(".")
	if err != nil {
		logrus.WithError(err).Panic("生成 tlog_gen.go 失败")
		return
	}

	svrDir := path.Join(dir, "tlog", "gen")
	xmlPath := path.Join(svrDir, "tlog.xml")
	//filePath := path.Join(svrDir, "tlog.template")
	//genFile(svrDir, xmlPath, filePath, "tlogservice.go")

	fullPath := path.Join("tlog", "tloggen.go")
	if err := gen.GenTlogService(xmlPath, fullPath); err != nil {
		logrus.WithError(err).Panicf("生成 %v 失败", fullPath)
	}

	//typeDir := path.Join(dir, "tlog", "tlog_context")
	//typePath := path.Join(typeDir, "tlog_type.template")
	//genFile(typeDir, xmlPath, typePath, "tlog_type_gen.go")
}
