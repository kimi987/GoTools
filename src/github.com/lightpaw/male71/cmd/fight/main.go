package main

import (
	"io/ioutil"
	"fmt"
	"github.com/lightpaw/male7combatx"
	"github.com/lightpaw/male7/pb/shared_proto"
)

func main() {

	filename := ""

	b, err := ioutil.ReadFile(filename)
	fmt.Println(err)

	proto := &shared_proto.CombatXProto{}
	fmt.Println(proto.Unmarshal(b))

	combatx.PrintResult(proto)
}
