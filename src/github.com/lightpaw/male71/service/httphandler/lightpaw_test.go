package httphandler

import (
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

//func TestAddYuanbao(t *testing.T) {
//
//	//id := 8041
//	//amount := 100
//
//	//addr := "111.230.185.163:7788"
//	addr := "118.126.107.11:7788"
//
//	//sendYuanbao(addr, id, amount)
//
//	//for i := 11201; i < 11301; i++ {
//	//	sendYuanbao(addr, i, 2000)
//	//}
//
//	//sendYuanbao(addr, 981+11200, 1)
//
//	ss := strings.Split(`
//m942`, "\n")
//
//	for _, v := range ss {
//		s := strings.TrimSpace(v)
//		if len(s) <= 0 {
//			continue
//		}
//
//		i, err := strconv.Atoi(s[1:])
//		if err != nil {
//			fmt.Println(addr, err)
//		} else {
//			fmt.Println(i)
//			sendYuanbao(addr, i+11200, 1)
//		}
//	}
//
//}
//
//func TestSendAll(t *testing.T) {
//	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/txstable"
//
//	db, err := sql.Open("mysql", dataSourceName)
//	if err != nil {
//		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
//		return
//	}
//
//	err = db.Ping()
//	if err != nil {
//		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
//		return
//	}
//
//	rows, err := db.Query("select id from hero")
//	if err != nil {
//		fmt.Println("查询hero失败，跳过测试", err.Error())
//		return
//	}
//
//	amount := 5000
//	addr := "111.230.185.163:7788"
//
//	for rows.Next() {
//		var id int
//		err = rows.Scan(&id)
//		if err != nil {
//			fmt.Println("hero扫描失败，跳过测试", err.Error())
//			return
//		}
//
//		sendYuanbao(addr, id, amount)
//
//	}
//}

func sendYuanbao(addr string, id, amount int) {
	key := []byte("suibianqige")

	timeStr := fmt.Sprintf("%v", time.Now().Unix())

	sign := computeHash(timeStr, key)

	url := fmt.Sprintf("http://%s/lightpaw/add_yuanbao?id=%v&amount=%v&time=%v&sign=%v", addr, id, amount, timeStr, sign)

	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(body))
}
