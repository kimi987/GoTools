package utils

import (
	"encoding/xml"
	"fmt"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/recordfile"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func readRf(st interface{}) *recordfile.RecordFile {
	rf, err := recordfile.New(st)
	if err != nil {
		log.Fatal("%v", err)
	}
	fn := reflect.TypeOf(st).Name() + ".txt"
	err = rf.Read("gamedata/" + fn)
	if err != nil {
		log.Fatal("%v: %v", fn, err)
	}

	return rf
}

func ReadXml(result interface{}, name string) {
	f, err := os.Open("conf/xml/" + name)

	defer f.Close()

	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(f)

	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(data, &result)

	if err != nil {
		panic(err)
	}

	fmt.Println("read xml data success")
}

func WriteXml(result interface{}, name string) {
	xmlOutPut, outPutErr := xml.MarshalIndent(result, "", "")
	if outPutErr == nil {
		//加入XML头
		headerBytes := []byte(xml.Header)

		xmlstring := string(xmlOutPut)

		xmlstring = strings.Replace(xmlstring, ">", ">\n", -1)

		xmlOutPut = []byte(xmlstring)
		//拼接XML头和实际XML内容
		xmlOutPutData := append(headerBytes, xmlOutPut...)
		//写入文件
		ioutil.WriteFile("conf/xml/"+name, xmlOutPutData, os.ModeAppend)

	} else {
		fmt.Println(outPutErr)
	}
}
