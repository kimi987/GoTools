
package main
 
import (
	"math"
	"encoding/binary"
	"errors"
	"encoding/base32"
	"io/ioutil"
    "fmt"
	"net/http"
	"time"
)
 
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "hello world")
}

func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r.URL = ", r.URL)
	// body,err := r.URL

	// if err != nil {
	// 	fmt.Println("err = ", err)
	// 	return
	// }
	// defer body.Close()
	// fmt.Println("body = ", body)
	fmt.Fprintln(w, "isOk")
}
 
func main() {
	go func() {
		time.Sleep(time.Second * 5)
		couponId,sequence, err :=  decodeCoupon("UREQCAI")

		if err != nil {
			fmt.Println("err = ", err)
		}
		CheckCouponInvalide("UREQCAI", couponId, sequence)
	}()
	// http.HandleFunc("/", IndexHandler)
	// http.HandleFunc("/validate", ValidateHandler)

	// http.ListenAndServe("0.0.0.0:5555", nil)
	time.Sleep(time.Second * 100)
}

var ErrInvalidCoupon = errors.New("invalid coupon")
var encoding = base32.StdEncoding.WithPadding(base32.NoPadding)

//DecodeCoupon 解码礼品码
func decodeCoupon(couponCode string) (couponId uint64, sequence uint64, err error) {
	buf, err := encoding.DecodeString(couponCode)
	if err != nil {
	  return 0, 0, err
	}

	fmt.Println("buf = ", buf)
  
	id, n := binary.Uvarint(buf[2:])
	if n < 0 {
	  return 0, 0, ErrInvalidCoupon
	}

	fmt.Println("id = ", id)
  
	if id > math.MaxUint32 {
	  return 0, 0, ErrInvalidCoupon
	}
  
	seq, n := binary.Uvarint(buf[2+n:])
	if n < 0 {
	  return 0, 0, ErrInvalidCoupon
	}
  
	fmt.Println("seq = ", seq)

	return id, seq, nil
}


//CheckCouponInvalide 检查礼品码是否有效
func CheckCouponInvalide(couponCode string,couponId, seq uint64) error{
	url := fmt.Sprintf("http://192.168.1.5:7893/change_coupon?couponid=%d&sequence=%d&coupon=%s" ,couponId,seq,couponCode)
	fmt.Println("url = ", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err = ", err)
		return err
	}
    defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("CheckCouponInvalide body = %v \n", string(body))

	return nil
}