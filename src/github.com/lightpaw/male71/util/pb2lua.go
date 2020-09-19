package util

import (
	"bytes"
	"github.com/lightpaw/protobuf/jsonpb"
	"github.com/lightpaw/protobuf/proto"
	"github.com/pkg/errors"
	"net/http"
	"io/ioutil"
	"github.com/lightpaw/male7/util/compress"
)

func Proto2LuaBytes(addr string, pb proto.Message) (string, []byte, error) {
	marshaler := jsonpb.Marshaler{}
	marshaler.EmitDefaults = true
	marshaler.EnumsAsInts = true
	jsonString, err := marshaler.MarshalToString(pb)
	if err != nil {
		return "", nil, errors.Wrapf(err, "proto to json fail")
	}

	jsonBytes := []byte(jsonString)
	key := Md5String(jsonBytes)

	resp, err := http.Get(addr + "/get?key=" + key)
	if err != nil {
		return "", nil, errors.Wrapf(err, "Proto2LuaBytes, get http.Get() fail")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, errors.Wrapf(err, "Proto2LuaBytes, http.Get() read body fail")
	}
	defer resp.Body.Close()

	if len(data) > 0 {
		return key, data, nil
	}

	compressBytes := compress.GzipCompress(jsonBytes)

	resp, err = http.Post(addr+"/generate?key="+key, "application/json", bytes.NewReader(compressBytes))
	if err != nil {
		return "", nil, errors.Wrapf(err, "Proto2LuaBytes, generate http.Post() fail")
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, errors.Wrapf(err, "Proto2LuaBytes, generate http.Post() read body fail")
	}
	defer resp.Body.Close()

	return key, data, nil
}
