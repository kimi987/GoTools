package ifacemock

//import (
//	"fmt"
//	"github.com/lightpaw/male7/gen/iface"
//	"reflect"
//	"runtime"
//	"testing"
//)
//
//func TestFuncName(t *testing.T) {
//
//	fmt.Println(GetFunctionName(ActionManager.Close))
//	fmt.Println(GetFunctionName(ActionManager.GetBaseLevel))
//	fmt.Println(getFunctionPointer(ActionManager.Close))
//	fmt.Println(getFunctionPointer(ActionManager.GetBaseLevel))
//
//	m1 := &MockActionManager{}
//	fmt.Println(GetFunctionName(m1.Close))
//	fmt.Println(GetFunctionName(m1.GetBaseLevel))
//	fmt.Println(getFunctionPointer(m1.Close))
//	fmt.Println(getFunctionPointer(m1.GetBaseLevel))
//
//	m2 := &MockActionManager{}
//	fmt.Println(GetFunctionName(m2.Close))
//	fmt.Println(GetFunctionName(m2.GetBaseLevel))
//	fmt.Println(getFunctionPointer(m2.Close))
//	fmt.Println(getFunctionPointer(m2.GetBaseLevel))
//
//	var im1 iface.ActionManager = m1
//	fmt.Println(GetFunctionName(im1.Close))
//	fmt.Println(GetFunctionName(im1.GetBaseLevel))
//	fmt.Println(getFunctionPointer(im1.Close))
//	fmt.Println(getFunctionPointer(im1.GetBaseLevel))
//
//	fmt.Println(m1.Close, m2.Close)
//}
//
//func TestMap(t *testing.T) {
//	m1 := &MockActionManager{}
//	m2 := &MockActionManager{}
//	m3 := &MockActionManager{}
//	var im1 iface.ActionManager = m1
//
//	dataMap := map[interface{}]int{}
//	dataMap[m1] = 1
//	dataMap[m2] = 2
//	dataMap[im1] = 3
//
//	fmt.Println(dataMap[m1])
//	fmt.Println(dataMap[m2])
//	fmt.Println(dataMap[im1])
//	fmt.Println(dataMap[m3])
//}
//
//func GetFunctionName(i interface{}) string {
//	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
//}
//
//func TestMock(t *testing.T) {
//
//	ActionManager.Close()
//	ActionManager.Mock(ActionManager.Close, func() {
//		fmt.Println(100)
//	})
//	ActionManager.Close()
//
//	m1 := &MockActionManager{}
//	m1.Close()
//
//	m1.Mock(m1.Close, func() {
//		fmt.Println(111)
//	})
//	m1.Close()
//
//	ActionManager.Close()
//
//	//// panic
//	//m1.Mock(m1.Close, func(int) {
//	//	fmt.Println(111)
//	//})
//	//m1.Close()
//}
