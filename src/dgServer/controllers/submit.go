package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"net/http/cookiejar"
	"time"
)

const (
	OriginURL     = "https://huke.163.com/openapi"
	Appkey        = "fc4e5c36757b4ebd83ac0464cde2377a"
	AppSecret     = "3d6a7126ea4f40299957c8e5b1a0d83a"
	GetFieldPath  = "/customer/fields"
	AddCustomPath = "/customer/addCustomer"
)

const (
	AddCustomData = `{
		"follower": "\"\"",
		"status": 1,
		"data": [
			%s
		]
	}`
)

type ReqData struct {
	Appkey   string `json:"ur-appkey"`
	Sign     string `json:"ur-sign"`
	Curtime  string `json:"ur-curtime"`
	Checksum string `json:"ur-checksum"`
}

type Fields struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

//SumitUserData 上传用户数据
func SumitUserData() {
	// SubmitURL := fmt.Sprintf("%s%s", OriginURL, AddCustomPath)

}

//GetFieldValue 获取字段
func GetFieldValue() {
	requestURL := fmt.Sprintf("%s%s", OriginURL, GetFieldPath)
	reqData := &ReqData{
		Appkey:  Appkey,
		Sign:    md5Str("{}"),
		Curtime: fmt.Sprintf("%d", time.Now().Unix()),
	}
	originData := fmt.Sprintf("%s%s%s", reqData.Appkey, reqData.Sign, reqData.Curtime)
	Checksum := AesEncrypt(originData, AppSecret)
	// if err != nil {
	// 	fmt.Println("err = ", err)
	// 	return
	// }
	reqData.Checksum = Checksum
	head := fmt.Sprintf("application/json;charset=utf-8;ur-appkey=%s;ur-sign=%s;ur-curtime=%s;ur-checksum=%s", reqData.Appkey, reqData.Sign, reqData.Curtime, reqData.Checksum)
	// fmt.Println("head = ", head)
	result := &Fields{}
	HttpsPost(requestURL, head, reqData, result)
	fmt.Printf("Code = %d \n", result.Code)
	fmt.Printf("Msg = %s \n", result.Msg)
	fmt.Printf("Data = %s \n", result.Data)
}

//HttpsPost Post data
func HttpsPost(url string, head string, arg *ReqData, reply interface{}) (err error) {
	var (
		response *http.Response
		body     []byte
		// buf      *bytes.Buffer
	)

	// if arg != nil {
	// 	if b, ok := arg.([]byte); !ok {
	// 		if body, err = json.Marshal(arg); err != nil {
	// 			return
	// 		}
	// 	} else {
	// 		body = b
	// 	}
	// }
	fmt.Printf("url = %v \n", url)
	// fmt.Printf("body = %v \n", string(body))
	// buf = bytes.NewBuffer(body)
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}

	client.Jar, _ = cookiejar.New(nil)

	req, err := http.NewRequest("POST", url, strings.NewReader(string("{}")))
	if err != nil {
		// handle error
		fmt.Println("err = ", err)
		return
	}

	fmt.Printf("Content-Type = application/json \n")
	fmt.Printf("ur-appkey = %s \n", Appkey)
	fmt.Printf("ur-sign = %s \n", md5Str("{}"))
	fmt.Printf("ur-curtime = %s \n", arg.Curtime)
	fmt.Printf("ur-checksum = %s \n", head)

	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("charset", "utf-8")
	// req.Header.Set("ur-appkey", Appkey)
	// req.Header.Set("ur-sign", md5Str("{}"))
	// req.Header.Set("ur-curtime", arg.Curtime)
	// originData := fmt.Sprintf("%s%s%s", arg.Appkey, arg.Sign, arg.Curtime)
	// req.Header.Set("ur-checksum", AesEncrypt(originData, AppSecret))

	// appKey:34ff472473d741b881a7b46a3cfb5b59
	// appSecret:d83d0971ca5641a7a6b2b7c78dea7f49
	// sign:3d6055a498cab778429574604b810725
	// curtime:1592276090
	// checksum:849685B5801196A66ED06BD254399E5E9F33AB18550979A9B79AA6F545F78A481D8B5D299556E231FBE2F6D321A50FFEDEA029E69B173478A8DF8A0CBD76898067AF2D3E4DCAACCCEE0697D488ECCAA3
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "utf-8")
	req.Header.Set("ur-appkey", "34ff472473d741b881a7b46a3cfb5b59")
	req.Header.Set("ur-sign", "3d6055a498cab778429574604b810725")
	req.Header.Set("ur-curtime", "1592276090")
	originData := fmt.Sprintf("%s%s%s", "34ff472473d741b881a7b46a3cfb5b59", "3d6055a498cab778429574604b810725", "1592276090")

	str1 := AESEncrypt1(([]byte)(originData), ([]byte)("d83d0971ca5641a7a6b2b7c78dea7f49"))
	fmt.Printf("ur-checksum = %s \n", str1)
	req.Header.Set("ur-checksum", "849685B5801196A66ED06BD254399E5E9F33AB18550979A9B79AA6F545F78A481D8B5D299556E231FBE2F6D321A50FFEDEA029E69B173478A8DF8A0CBD76898067AF2D3E4DCAACCCEE0697D488ECCAA3")

	response, err = client.Do(req)

	if err != nil {
		// handle error
		fmt.Println("err = ", err)
		return
	}

	defer response.Body.Close()

	if body, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}
	fmt.Printf("result body = %v \n", string(body))
	return json.Unmarshal(body, reply)
}

func AESEncrypt(key []byte, text string) (string, bool) {
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", false
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", false
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	// fmt.Printf("ciphertext =  %s", string(ciphertext))
	return base64.URLEncoding.EncodeToString(ciphertext), true
}

func AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

func AESEncrypt1(src []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func md5Str(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
