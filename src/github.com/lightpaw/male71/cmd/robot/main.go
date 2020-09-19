package main

import (
	"net"
	"os"
	"github.com/lightpaw/male7/gen/pb/login"
	"bytes"
	"bufio"
	"encoding/base64"
	"strconv"
	"github.com/golang/snappy"
	"compress/gzip"
	"io/ioutil"
	"fmt"
	"github.com/lightpaw/male7/gen/pb/chat"
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/pb/shared_proto"
)

func main() {

	// 开启一个tcp连接
	//servAddr := "116.247.86.50:7775"
	servAddr := "192.168.1.5:7775"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// 先发送一个loginToken（对于机器人来说，这个是个固定值）
	loginToken, err := base64.URLEncoding.DecodeString("Y3NsSS8U1_iBIjP0dNgGCxPLBTfl8JO2YGYXQFCPSfjGA9uuUBcwq5YcOAUNXJTtHCn5UZicfNB128Y8C8iRtHAAAAAAAAAAAAAAAAAAAAAA")
	//loginToken, err := base64.URLEncoding.DecodeString("Y9KS-qv5KuUhBaJbBRsxUigBbdh2BXRP9e28LMXpmxoSHpOTpNW-0pXZAIso78d1TvWA8B3IBPXHv1KiEmqVm-0AAAAAAAAAAAAAAAAAAAAA")
	if err != nil {
		println("parse login token failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write(loginToken)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	// 然后开始正常的协议包交互
	connReader := bufio.NewReader(conn)
	testLogin(conn, connReader, 1)

	testChat(conn, connReader)

	//// 获取配置文件
	//_, err = conn.Write(newSendMsg(5, 3, nil))
	//if err != nil {
	//	println("Write to server failed:", err.Error())
	//	os.Exit(1)
	//}
	//if moduleId, msgId, protoBytes, err := readMsg(connReader); err != nil {
	//	println("Write to server failed:", err.Error())
	//	os.Exit(1)
	//} else {
	//	if moduleId == 0 && msgId == 1 {
	//		// 解压缩
	//		decodeBytes, _ := snappy.Decode(nil, protoBytes)
	//
	//		r := bufio.NewReader(bytes.NewBuffer(decodeBytes))
	//		// 重新读取 ModuleId 和 MsgId
	//		moduleId, n1, _ := readVarInt32(r)
	//		msgId, n2, _ := readVarInt32(r)
	//		protoBytes := decodeBytes[n1+n2:]
	//		if moduleId == 5 && msgId == 4 {
	//			proto := &misc.S2CConfigProto{}
	//			if err := proto.Unmarshal(protoBytes); err != nil {
	//				println("Write to server failed:", err.Error())
	//				os.Exit(1)
	//			}
	//		}
	//	}
	//}

	// 发送机器人登陆（跟正常登陆不一样，正常登陆id包含在loginToken中），回包跟正常登陆一致
	// c2s_robot_login = 1_13 // C2SRobotLoginProto
	//object := &login.C2SRobotLoginProto{Id: 4} // 机器人id从100001开始
	//protodata, _ := object.Marshal()
	//_, err = conn.Write(newSendMsg(1, 13, protodata))
	//if err != nil {
	//	println("Write to server failed:", err.Error())
	//	os.Exit(1)
	//}
	//
	//if moduleId, msgId, protoBytes, err := readMsg(connReader); err != nil {
	//	println("Write to server failed:", err.Error())
	//	os.Exit(1)
	//} else {
	//
	//	println("moduleId=", moduleId)
	//	println("msgId=", msgId)
	//
	//	// 没有完成新手引导的英雄，走新手引导流程
	//	// s2c_tutorial_progress = 1_17 // S2CTutorialProgressProto
	//	if moduleId == 1 && msgId == 17 {
	//		proto := &login.S2CTutorialProgressProto{}
	//		if err := proto.Unmarshal(protoBytes); err != nil {
	//			println("Write to server failed:", err.Error())
	//			os.Exit(1)
	//		}
	//
	//		println("新手引导进度: ", proto.Progress)
	//
	//		// 设置新建教程完成
	//		object := &login.C2SSetTutorialProgressProto{Progress: 10000, IsComplete: true}
	//		protodata, _ := object.Marshal()
	//		_, err = conn.Write(newSendMsg(1, 18, protodata))
	//		if err != nil {
	//			println("Write to server failed:", err.Error())
	//			os.Exit(1)
	//		} else {
	//			// 走创建角色流程 ...
	//			println("角色未创建")
	//
	//			// 设置新建教程完成
	//			object := &login.C2SCreateHeroProto{}
	//			protodata, _ := object.Marshal()
	//			_, err = conn.Write(newSendMsg(1, 3, protodata))
	//			if err != nil {
	//				println("Write to server failed:", err.Error())
	//				os.Exit(1)
	//			}
	//
	//			for {
	//				// 下面这段是演示代码，演示如何读取错误码
	//				// 1: 已经登陆了，不要重复登陆
	//				// 2: 发送上来的proto解析不了
	//				// 3: 发送的id无效
	//				// 4: 被T下线
	//				// 5: 服务器忙，请稍后再试
	//				// s2c_fail_login = 1_9 // 错误码
	//				if moduleId, msgId, protoBytes, err := readMsg(connReader); err != nil {
	//					println("Write to server failed:", err.Error())
	//					os.Exit(1)
	//				} else {
	//
	//					if moduleId == 1 && msgId == 4 {
	//						// 创建成功
	//						println("角色创建成功11111")
	//
	//					} else if moduleId == 1 && msgId == 9 {
	//						// 创建失败，查看错误码（读取一个varint32）
	//
	//						errorCode, _, _ := readVarInt32(bufio.NewReader(bytes.NewReader(protoBytes)))
	//						switch errorCode {
	//						case 1:
	//							println("已经登陆了，不要重复登陆")
	//						case 2:
	//							println("发送上来的proto解析不了")
	//						default:
	//							println("undown error code", errorCode)
	//						}
	//					}
	//				}
	//			}
	//
	//		}
	//
	//	} else if moduleId == 1 && msgId == 8 {
	//		// 已完成新手引导的玩家
	//		// s2c_login = 1_8 // S2CLoginProto
	//		proto := &login.S2CLoginProto{}
	//		if err := proto.Unmarshal(protoBytes); err != nil {
	//			println("Write to server failed:", err.Error())
	//			os.Exit(1)
	//		}
	//
	//		if proto.Created {
	//			// 已创建角色走登陆流程 ...
	//			println("角色已创建")
	//
	//			_, err = conn.Write(newSendMsg(1, 10, nil))
	//			if err != nil {
	//				println("Write to server failed:", err.Error())
	//				os.Exit(1)
	//			}
	//		} else {
	//
	//			// 走创建角色流程 ...
	//			println("角色未创建")
	//
	//			// 设置新建教程完成
	//			object := &login.C2SCreateHeroProto{}
	//			protodata, _ := object.Marshal()
	//			_, err = conn.Write(newSendMsg(1, 3, protodata))
	//			if err != nil {
	//				println("Write to server failed:", err.Error())
	//				os.Exit(1)
	//			}
	//
	//			// 下面这段是演示代码，演示如何读取错误码
	//			// 1: 已经登陆了，不要重复登陆
	//			// 2: 发送上来的proto解析不了
	//			// 3: 发送的id无效
	//			// 4: 被T下线
	//			// 5: 服务器忙，请稍后再试
	//			// s2c_fail_login = 1_9 // 错误码
	//			if moduleId, msgId, protoBytes, err := readMsg(connReader); err != nil {
	//				println("Write to server failed:", err.Error())
	//				os.Exit(1)
	//			} else {
	//
	//				if moduleId == 1 && msgId == 4 {
	//					// 创建成功
	//					println("角色创建成功")
	//				} else if moduleId == 1 && msgId == 9 {
	//					// 创建失败，查看错误码（读取一个varint32）
	//
	//					errorCode, _, _ := readVarInt32(bufio.NewReader(bytes.NewReader(protoBytes)))
	//					switch errorCode {
	//					case 1:
	//						println("已经登陆了，不要重复登陆")
	//					case 2:
	//						println("发送上来的proto解析不了")
	//					}
	//				}
	//			}
	//		}
	//
	//	} else {
	//		println("unkown msg:", moduleId, msgId)
	//		os.Exit(1)
	//	}
	//}
	//
	//conn.Close()
}

var (
	m1_4  = NewMs(1, 4)
	m1_8  = NewMs(1, 8)
	m1_9  = NewMs(1, 9)
	m1_11 = NewMs(1, 11)
	m1_17 = NewMs(1, 17)
	m1_18 = NewMs(1, 18)
)

func testLogin(conn net.Conn, reader *bufio.Reader, id int32) {
	var err error

	object := &login.C2SRobotLoginProto{Id: id} // 机器人id从100001开始
	protodata, _ := object.Marshal()
	_, err = conn.Write(newSendMsg(1, 13, protodata))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	createHero := func() {
		// 走创建角色流程 ...
		println("角色未创建")

		object := &login.C2SCreateHeroProto{}
		protodata, _ := object.Marshal()
		_, err = conn.Write(newSendMsg(1, 3, protodata))
		if err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		// 等等创建角色成功
		msg, protoBytes := waitMsg(reader, m1_4, m1_9)
		switch msg {
		case m1_4:
			// 创建成功
			println("角色创建成功")
		case m1_9:
			// 创建失败，查看错误码（读取一个varint32）

			errorCode, _, _ := readVarInt32(bufio.NewReader(bytes.NewReader(protoBytes)))
			switch errorCode {
			case 1:
				println("已经登陆了，不要重复登陆")
			case 2:
				println("发送上来的proto解析不了")
			default:
				println("undown error code", errorCode)
			}
			os.Exit(1)
		}
	}

	msg, protoBytes := waitMsg(reader, m1_17, m1_8)
	switch msg {
	case m1_17:
		// 设置新建教程完成
		if true {
			object := &login.C2SSetTutorialProgressProto{Progress: 10000, IsComplete: true}
			protodata, _ := object.Marshal()
			_, err = conn.Write(newSendMsg(1, 18, protodata))
			if err != nil {
				println("Write to server failed:", err.Error())
				os.Exit(1)
			}
		}
		createHero()
	case m1_8:
		// 已完成新手引导的玩家
		// s2c_login = 1_8 // S2CLoginProto
		proto := &login.S2CLoginProto{}
		if err := proto.Unmarshal(protoBytes); err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		if proto.Created {
			// 已创建角色走登陆流程 ...
			println("角色已创建")
		} else {
			createHero()
		}
	}

	// loaded
	_, err = conn.Write(newSendMsg(1, 10, nil))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	waitMsg(reader, m1_11)
	println("loaded")
}

var (
	m13_15 = NewMs(13, 15)
	m13_16 = NewMs(13, 16)
)

func testChat(conn net.Conn, reader *bufio.Reader) {
	var err error

	sendMsg := func(n int) {

		b := bytes.Buffer{}
		for i := 0; i < n; i++ {
			b.WriteString("哈")
		}

		object := &chat.C2SSendChatProto{
			ChatMsg: util.SafeMarshal(&shared_proto.ChatMsgProto{
				Text: b.String(),
			}),
		}
		protodata, _ := object.Marshal()
		println("聊天消息 ", len(protodata))
		_, err = conn.Write(newSendMsg(13, 14, protodata))
		if err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		waitMsg(reader, m13_15, m13_16)
	}

	sendMsg(10)
	sendMsg(50)
	sendMsg(100)
	sendMsg(500)
	sendMsg(1000)
}

func NewMs(moid, msid int) ms {
	return ms(m2i(moid, msid))
}

type ms int

func (s ms) ModuleId() int {
	return int(s / 10000)
}

func (s ms) MsgId() int {
	return int(s % 10000)
}

func m2i(moid, msid int) int {
	return moid*10000 + msid
}

func i2m(i int) (int, int) {
	return i / 10000, i % 10000
}

func waitMsg(reader *bufio.Reader, i ...ms) (m ms, proto []byte) {
	for {
		if moduleId, msgId, proto, err := readMsg(reader); err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		} else {
			if moduleId == 0 {
				switch msgId {
				case 0:
					proto, err = snappy.Decode(nil, proto)
					if err != nil {
						println("snappy uncompress fail:", err.Error())
						os.Exit(1)
					}
				case 1:
					r, err := gzip.NewReader(bytes.NewReader(proto))
					if err != nil {
						println("gzip uncompress fail:", err.Error())
						os.Exit(1)
					}

					proto, err = ioutil.ReadAll(r)
					if err != nil {
						println("gzip uncompress read fail:", err.Error())
						os.Exit(1)
					}
				default:
					println("Unkown compress type, " + strconv.Itoa(msgId))
					os.Exit(1)
				}

				moduleId, msgId, proto, err = readMsgWithLen(bufio.NewReader(bytes.NewReader(proto)), len(proto))
				if err != nil {
					println("Write to server failed:", err.Error())
					os.Exit(1)
				}
			}

			fmt.Println("收到消息", moduleId, msgId)

			for _, v := range i {
				if moduleId == v.ModuleId() && msgId == v.MsgId() {
					return v, proto
				}
			}
		}
	}
}

func newSendMsg(moduleId, msgId int, proto []byte) []byte {

	b := bytes.Buffer{}
	b.Write([]byte{0, 0, 0, 0}) // 消息长度占位
	b.Write(encodeVarInt32(moduleId))
	b.Write(encodeVarInt32(msgId))
	b.Write(proto)

	buf := b.Bytes()

	length := len(buf) - 3
	if length <= 127 { // 2 ^ 7 - 1
		buf[2] = uint8(length)
		return buf[2:]
	}

	if length <= 16383 { // 2 ^ 14 - 1
		buf[1], buf[2] = encodeVarint2(length)
		return buf[1:]
	}

	// 2 ^ 21
	buf[0], buf[1], buf[2] = encodeVarint3(length)
	return buf
}

func encodeVarInt32(n int) []byte {
	if n <= 127 { // 2 ^ 7 - 1
		return []byte{byte(n)}
	}

	if n <= 16383 { // 2 ^ 14 - 1
		x, y := encodeVarint2(n)
		return []byte{x, y}
	}

	// 2 ^ 21
	x, y, z := encodeVarint3(n)
	return []byte{x, y, z}
}

func encodeVarint2(n int) (uint8, uint8) {
	return (0x80 | uint8(n&0x7f)), uint8(n >> 7)
}

// first bit of the third byte is also considered data
func encodeVarint3(n int) (uint8, uint8, uint8) {
	return (0x80 | uint8(n&0x7f)), (0x80 | uint8((n>>7)&0x7f)), uint8(n >> 14)
}

func readMsg(conn *bufio.Reader) (moduleId int, msgId int, proto []byte, err error) {

	len, _, err := readVarInt32(conn)
	if err != nil {
		return 0, 0, nil, err
	}

	return readMsgWithLen(conn, len)
}

func readMsgWithLen(conn *bufio.Reader, len int) (moduleId int, msgId int, proto []byte, err error) {

	buf := make([]byte, len)
	if _, err := conn.Read(buf); err != nil {
		return 0, 0, nil, err
	}

	r := bufio.NewReader(bytes.NewReader(buf))

	var n1, n2 int
	moduleId, n1, _ = readVarInt32(r)
	msgId, n2, _ = readVarInt32(r)
	proto = buf[n1+n2:]

	return moduleId, msgId, proto, nil
}

func readVarInt32(r *bufio.Reader) (int, int, error) {

	readLen := 0

	n1, err := r.ReadByte()
	if err != nil {
		return 0, readLen, err
	}
	readLen++

	if n1 <= 127 {
		return int(n1), readLen, nil
	}

	n2, err := r.ReadByte()
	if err != nil {
		return 0, readLen, err
	}
	readLen++

	if n2 <= 127 {
		return decodeVarint2(n1, n2), readLen, nil
	}

	n3, err := r.ReadByte()
	if err != nil {
		return 0, readLen, err
	}
	readLen++

	return decodeVarint3(n1, n2, n3), readLen, nil
}

func decodeVarint2(n1, n2 byte) int {
	return (int(n2) << 7) | (int(n1) & 0x7f)
}

// first bit of the third byte is also considered as data
func decodeVarint3(n1, n2, n3 byte) int {
	return (int(n3) << 14) | ((int(n2) & 0x7f) << 7) | (int(n1) & 0x7f)
}
