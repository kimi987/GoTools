// Code generated by protoc-gen-gogo.
// source: github.com/lightpaw/male7/pb/shared_proto/question.proto
// DO NOT EDIT!

package shared_proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// 回答状态
type HeroQuestionState int32

const (
	HeroQuestionState_QUESTION_INVALID_STATE HeroQuestionState = 0
	HeroQuestionState_QUESTION_RIGHT         HeroQuestionState = 1
	HeroQuestionState_QUESTION_WRONG         HeroQuestionState = 2
	HeroQuestionState_QUESTION_WAIT          HeroQuestionState = 3
)

var HeroQuestionState_name = map[int32]string{
	0: "QUESTION_INVALID_STATE",
	1: "QUESTION_RIGHT",
	2: "QUESTION_WRONG",
	3: "QUESTION_WAIT",
}
var HeroQuestionState_value = map[string]int32{
	"QUESTION_INVALID_STATE": 0,
	"QUESTION_RIGHT":         1,
	"QUESTION_WRONG":         2,
	"QUESTION_WAIT":          3,
}

func (x HeroQuestionState) String() string {
	return proto.EnumName(HeroQuestionState_name, int32(x))
}
func (HeroQuestionState) EnumDescriptor() ([]byte, []int) { return fileDescriptorQuestion, []int{0} }

// 答案的index从1开始递增
type QuestionProto struct {
	Id          int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Question    string   `protobuf:"bytes,2,opt,name=question,proto3" json:"question,omitempty"`
	RightAnswer string   `protobuf:"bytes,3,opt,name=right_answer,json=rightAnswer,proto3" json:"right_answer,omitempty"`
	WrongAnswer []string `protobuf:"bytes,4,rep,name=wrong_answer,json=wrongAnswer" json:"wrong_answer,omitempty"`
}

func (m *QuestionProto) Reset()                    { *m = QuestionProto{} }
func (m *QuestionProto) String() string            { return proto.CompactTextString(m) }
func (*QuestionProto) ProtoMessage()               {}
func (*QuestionProto) Descriptor() ([]byte, []int) { return fileDescriptorQuestion, []int{0} }

func (m *QuestionProto) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *QuestionProto) GetQuestion() string {
	if m != nil {
		return m.Question
	}
	return ""
}

func (m *QuestionProto) GetRightAnswer() string {
	if m != nil {
		return m.RightAnswer
	}
	return ""
}

func (m *QuestionProto) GetWrongAnswer() []string {
	if m != nil {
		return m.WrongAnswer
	}
	return nil
}

type QuestionSayingProto struct {
	Id      int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Content string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	Author  string `protobuf:"bytes,3,opt,name=author,proto3" json:"author,omitempty"`
}

func (m *QuestionSayingProto) Reset()                    { *m = QuestionSayingProto{} }
func (m *QuestionSayingProto) String() string            { return proto.CompactTextString(m) }
func (*QuestionSayingProto) ProtoMessage()               {}
func (*QuestionSayingProto) Descriptor() ([]byte, []int) { return fileDescriptorQuestion, []int{1} }

func (m *QuestionSayingProto) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *QuestionSayingProto) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *QuestionSayingProto) GetAuthor() string {
	if m != nil {
		return m.Author
	}
	return ""
}

type QuestionPrizeProto struct {
	Score int32       `protobuf:"varint,1,opt,name=score,proto3" json:"score,omitempty"`
	Prize *PrizeProto `protobuf:"bytes,2,opt,name=prize" json:"prize,omitempty"`
}

func (m *QuestionPrizeProto) Reset()                    { *m = QuestionPrizeProto{} }
func (m *QuestionPrizeProto) String() string            { return proto.CompactTextString(m) }
func (*QuestionPrizeProto) ProtoMessage()               {}
func (*QuestionPrizeProto) Descriptor() ([]byte, []int) { return fileDescriptorQuestion, []int{2} }

func (m *QuestionPrizeProto) GetScore() int32 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *QuestionPrizeProto) GetPrize() *PrizeProto {
	if m != nil {
		return m.Prize
	}
	return nil
}

type QuestionMiscProto struct {
	MaxTimes      int32 `protobuf:"varint,1,opt,name=max_times,json=maxTimes,proto3" json:"max_times,omitempty"`
	QuestionCount int32 `protobuf:"varint,2,opt,name=question_count,json=questionCount,proto3" json:"question_count,omitempty"`
}

func (m *QuestionMiscProto) Reset()                    { *m = QuestionMiscProto{} }
func (m *QuestionMiscProto) String() string            { return proto.CompactTextString(m) }
func (*QuestionMiscProto) ProtoMessage()               {}
func (*QuestionMiscProto) Descriptor() ([]byte, []int) { return fileDescriptorQuestion, []int{3} }

func (m *QuestionMiscProto) GetMaxTimes() int32 {
	if m != nil {
		return m.MaxTimes
	}
	return 0
}

func (m *QuestionMiscProto) GetQuestionCount() int32 {
	if m != nil {
		return m.QuestionCount
	}
	return 0
}

type HeroQuestionProto struct {
	UsedTimes              int32                    `protobuf:"varint,1,opt,name=used_times,json=usedTimes,proto3" json:"used_times,omitempty"`
	CurrentQuestion        []*HeroEachQuestionProto `protobuf:"bytes,2,rep,name=current_question,json=currentQuestion" json:"current_question,omitempty"`
	AllRightQuestionAmount int32                    `protobuf:"varint,3,opt,name=all_right_question_amount,json=allRightQuestionAmount,proto3" json:"all_right_question_amount,omitempty"`
}

func (m *HeroQuestionProto) Reset()                    { *m = HeroQuestionProto{} }
func (m *HeroQuestionProto) String() string            { return proto.CompactTextString(m) }
func (*HeroQuestionProto) ProtoMessage()               {}
func (*HeroQuestionProto) Descriptor() ([]byte, []int) { return fileDescriptorQuestion, []int{4} }

func (m *HeroQuestionProto) GetUsedTimes() int32 {
	if m != nil {
		return m.UsedTimes
	}
	return 0
}

func (m *HeroQuestionProto) GetCurrentQuestion() []*HeroEachQuestionProto {
	if m != nil {
		return m.CurrentQuestion
	}
	return nil
}

func (m *HeroQuestionProto) GetAllRightQuestionAmount() int32 {
	if m != nil {
		return m.AllRightQuestionAmount
	}
	return 0
}

type HeroEachQuestionProto struct {
	CurrentQuestion int32             `protobuf:"varint,1,opt,name=current_question,json=currentQuestion,proto3" json:"current_question,omitempty"`
	State           HeroQuestionState `protobuf:"varint,2,opt,name=state,proto3,enum=proto.HeroQuestionState" json:"state,omitempty"`
	Answer          int32             `protobuf:"varint,3,opt,name=answer,proto3" json:"answer,omitempty"`
}

func (m *HeroEachQuestionProto) Reset()                    { *m = HeroEachQuestionProto{} }
func (m *HeroEachQuestionProto) String() string            { return proto.CompactTextString(m) }
func (*HeroEachQuestionProto) ProtoMessage()               {}
func (*HeroEachQuestionProto) Descriptor() ([]byte, []int) { return fileDescriptorQuestion, []int{5} }

func (m *HeroEachQuestionProto) GetCurrentQuestion() int32 {
	if m != nil {
		return m.CurrentQuestion
	}
	return 0
}

func (m *HeroEachQuestionProto) GetState() HeroQuestionState {
	if m != nil {
		return m.State
	}
	return HeroQuestionState_QUESTION_INVALID_STATE
}

func (m *HeroEachQuestionProto) GetAnswer() int32 {
	if m != nil {
		return m.Answer
	}
	return 0
}

func init() {
	proto.RegisterType((*QuestionProto)(nil), "proto.QuestionProto")
	proto.RegisterType((*QuestionSayingProto)(nil), "proto.QuestionSayingProto")
	proto.RegisterType((*QuestionPrizeProto)(nil), "proto.QuestionPrizeProto")
	proto.RegisterType((*QuestionMiscProto)(nil), "proto.QuestionMiscProto")
	proto.RegisterType((*HeroQuestionProto)(nil), "proto.HeroQuestionProto")
	proto.RegisterType((*HeroEachQuestionProto)(nil), "proto.HeroEachQuestionProto")
	proto.RegisterEnum("proto.HeroQuestionState", HeroQuestionState_name, HeroQuestionState_value)
}
func (m *QuestionProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuestionProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(m.Id))
	}
	if len(m.Question) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(len(m.Question)))
		i += copy(dAtA[i:], m.Question)
	}
	if len(m.RightAnswer) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(len(m.RightAnswer)))
		i += copy(dAtA[i:], m.RightAnswer)
	}
	if len(m.WrongAnswer) > 0 {
		for _, s := range m.WrongAnswer {
			dAtA[i] = 0x22
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	return i, nil
}

func (m *QuestionSayingProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuestionSayingProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(m.Id))
	}
	if len(m.Content) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(len(m.Content)))
		i += copy(dAtA[i:], m.Content)
	}
	if len(m.Author) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(len(m.Author)))
		i += copy(dAtA[i:], m.Author)
	}
	return i, nil
}

func (m *QuestionPrizeProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuestionPrizeProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Score != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(m.Score))
	}
	if m.Prize != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(m.Prize.Size()))
		n1, err := m.Prize.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func (m *QuestionMiscProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuestionMiscProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.MaxTimes != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(m.MaxTimes))
	}
	if m.QuestionCount != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(m.QuestionCount))
	}
	return i, nil
}

func (m *HeroQuestionProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HeroQuestionProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.UsedTimes != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(m.UsedTimes))
	}
	if len(m.CurrentQuestion) > 0 {
		for _, msg := range m.CurrentQuestion {
			dAtA[i] = 0x12
			i++
			i = encodeVarintQuestion(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.AllRightQuestionAmount != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(m.AllRightQuestionAmount))
	}
	return i, nil
}

func (m *HeroEachQuestionProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HeroEachQuestionProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.CurrentQuestion != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(m.CurrentQuestion))
	}
	if m.State != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(m.State))
	}
	if m.Answer != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintQuestion(dAtA, i, uint64(m.Answer))
	}
	return i, nil
}

func encodeFixed64Question(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Question(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintQuestion(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *QuestionProto) Size() (n int) {
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovQuestion(uint64(m.Id))
	}
	l = len(m.Question)
	if l > 0 {
		n += 1 + l + sovQuestion(uint64(l))
	}
	l = len(m.RightAnswer)
	if l > 0 {
		n += 1 + l + sovQuestion(uint64(l))
	}
	if len(m.WrongAnswer) > 0 {
		for _, s := range m.WrongAnswer {
			l = len(s)
			n += 1 + l + sovQuestion(uint64(l))
		}
	}
	return n
}

func (m *QuestionSayingProto) Size() (n int) {
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovQuestion(uint64(m.Id))
	}
	l = len(m.Content)
	if l > 0 {
		n += 1 + l + sovQuestion(uint64(l))
	}
	l = len(m.Author)
	if l > 0 {
		n += 1 + l + sovQuestion(uint64(l))
	}
	return n
}

func (m *QuestionPrizeProto) Size() (n int) {
	var l int
	_ = l
	if m.Score != 0 {
		n += 1 + sovQuestion(uint64(m.Score))
	}
	if m.Prize != nil {
		l = m.Prize.Size()
		n += 1 + l + sovQuestion(uint64(l))
	}
	return n
}

func (m *QuestionMiscProto) Size() (n int) {
	var l int
	_ = l
	if m.MaxTimes != 0 {
		n += 1 + sovQuestion(uint64(m.MaxTimes))
	}
	if m.QuestionCount != 0 {
		n += 1 + sovQuestion(uint64(m.QuestionCount))
	}
	return n
}

func (m *HeroQuestionProto) Size() (n int) {
	var l int
	_ = l
	if m.UsedTimes != 0 {
		n += 1 + sovQuestion(uint64(m.UsedTimes))
	}
	if len(m.CurrentQuestion) > 0 {
		for _, e := range m.CurrentQuestion {
			l = e.Size()
			n += 1 + l + sovQuestion(uint64(l))
		}
	}
	if m.AllRightQuestionAmount != 0 {
		n += 1 + sovQuestion(uint64(m.AllRightQuestionAmount))
	}
	return n
}

func (m *HeroEachQuestionProto) Size() (n int) {
	var l int
	_ = l
	if m.CurrentQuestion != 0 {
		n += 1 + sovQuestion(uint64(m.CurrentQuestion))
	}
	if m.State != 0 {
		n += 1 + sovQuestion(uint64(m.State))
	}
	if m.Answer != 0 {
		n += 1 + sovQuestion(uint64(m.Answer))
	}
	return n
}

func sovQuestion(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozQuestion(x uint64) (n int) {
	return sovQuestion(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QuestionProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuestion
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QuestionProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuestionProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Question", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuestion
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Question = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RightAnswer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuestion
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RightAnswer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WrongAnswer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuestion
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.WrongAnswer = append(m.WrongAnswer, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuestion(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuestion
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QuestionSayingProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuestion
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QuestionSayingProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuestionSayingProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Content", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuestion
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Content = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Author", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuestion
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Author = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuestion(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuestion
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QuestionPrizeProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuestion
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QuestionPrizeProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuestionPrizeProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Score", wireType)
			}
			m.Score = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Score |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prize", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuestion
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Prize == nil {
				m.Prize = &PrizeProto{}
			}
			if err := m.Prize.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuestion(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuestion
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QuestionMiscProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuestion
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QuestionMiscProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuestionMiscProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxTimes", wireType)
			}
			m.MaxTimes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxTimes |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field QuestionCount", wireType)
			}
			m.QuestionCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.QuestionCount |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuestion(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuestion
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *HeroQuestionProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuestion
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HeroQuestionProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HeroQuestionProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UsedTimes", wireType)
			}
			m.UsedTimes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UsedTimes |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentQuestion", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuestion
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CurrentQuestion = append(m.CurrentQuestion, &HeroEachQuestionProto{})
			if err := m.CurrentQuestion[len(m.CurrentQuestion)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AllRightQuestionAmount", wireType)
			}
			m.AllRightQuestionAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AllRightQuestionAmount |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuestion(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuestion
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *HeroEachQuestionProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuestion
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HeroEachQuestionProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HeroEachQuestionProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentQuestion", wireType)
			}
			m.CurrentQuestion = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CurrentQuestion |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= (HeroQuestionState(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Answer", wireType)
			}
			m.Answer = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Answer |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuestion(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthQuestion
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuestion(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuestion
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuestion
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthQuestion
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowQuestion
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipQuestion(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthQuestion = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuestion   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/lightpaw/male7/pb/shared_proto/question.proto", fileDescriptorQuestion)
}

var fileDescriptorQuestion = []byte{
	// 527 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xc1, 0xae, 0xd2, 0x40,
	0x14, 0x7d, 0x03, 0xf6, 0xf9, 0xde, 0xe5, 0x81, 0x30, 0x2a, 0x41, 0x54, 0x82, 0x4d, 0x8c, 0xe8,
	0x02, 0x12, 0x34, 0x51, 0x97, 0xa8, 0x04, 0x48, 0x94, 0xe7, 0x2b, 0x55, 0x12, 0x37, 0xcd, 0x50,
	0x26, 0x50, 0xd3, 0x76, 0xb0, 0x9d, 0x86, 0xa7, 0x3b, 0xb7, 0x7e, 0x81, 0x5f, 0xe2, 0x37, 0xb8,
	0xf4, 0x13, 0x0c, 0xfe, 0x88, 0x99, 0x99, 0x0e, 0xd0, 0xc8, 0xc2, 0x55, 0x73, 0xce, 0xbd, 0xf7,
	0xdc, 0x73, 0xef, 0x9d, 0xc2, 0xb3, 0x85, 0xc7, 0x97, 0xc9, 0xac, 0xed, 0xb2, 0xa0, 0xe3, 0x7b,
	0x8b, 0x25, 0x5f, 0x91, 0x75, 0x27, 0x20, 0x3e, 0x7d, 0xda, 0x59, 0xcd, 0x3a, 0xf1, 0x92, 0x44,
	0x74, 0xee, 0xac, 0x22, 0xc6, 0x59, 0xe7, 0x53, 0x42, 0x63, 0xee, 0xb1, 0xb0, 0x2d, 0x21, 0x36,
	0xe4, 0xa7, 0xfe, 0xe4, 0xff, 0x05, 0x66, 0x24, 0xa6, 0xaa, 0xd8, 0xfc, 0x8a, 0xa0, 0x78, 0x91,
	0xea, 0xbd, 0x95, 0x72, 0x25, 0xc8, 0x79, 0xf3, 0x1a, 0x6a, 0xa2, 0x96, 0x61, 0xe5, 0xbc, 0x39,
	0xae, 0xc3, 0x89, 0x6e, 0x58, 0xcb, 0x35, 0x51, 0xeb, 0xd4, 0xda, 0x62, 0x7c, 0x0f, 0xce, 0x22,
	0xd1, 0xca, 0x21, 0x61, 0xbc, 0xa6, 0x51, 0x2d, 0x2f, 0xe3, 0x05, 0xc9, 0xf5, 0x24, 0x25, 0x52,
	0xd6, 0x11, 0x0b, 0x17, 0x3a, 0xe5, 0x4a, 0x33, 0x2f, 0x52, 0x24, 0xa7, 0x52, 0xcc, 0x29, 0x5c,
	0xd7, 0x16, 0x26, 0xe4, 0xb3, 0x17, 0x2e, 0x0e, 0x1b, 0xa9, 0xc1, 0x55, 0x97, 0x85, 0x9c, 0x86,
	0x3c, 0xf5, 0xa1, 0x21, 0xae, 0xc2, 0x31, 0x49, 0xf8, 0x92, 0x69, 0x03, 0x29, 0x32, 0x27, 0x80,
	0x77, 0xb3, 0x79, 0x5f, 0xa8, 0xd2, 0xbd, 0x01, 0x46, 0xec, 0xb2, 0x88, 0xa6, 0xd2, 0x0a, 0xe0,
	0x07, 0x60, 0xac, 0x44, 0x8e, 0xd4, 0x2e, 0x74, 0x2b, 0x6a, 0x3f, 0xed, 0x5d, 0x9d, 0xa5, 0xe2,
	0xe6, 0x14, 0x2a, 0x5a, 0xf4, 0x8d, 0x17, 0xbb, 0x4a, 0xf3, 0x36, 0x9c, 0x06, 0xe4, 0xd2, 0xe1,
	0x5e, 0x40, 0xe3, 0x54, 0xf7, 0x24, 0x20, 0x97, 0xb6, 0xc0, 0xf8, 0x3e, 0x94, 0xf4, 0xc6, 0x1c,
	0x97, 0x25, 0xa9, 0x7f, 0xc3, 0x2a, 0x6a, 0xf6, 0xa5, 0x20, 0xcd, 0x1f, 0x08, 0x2a, 0x43, 0x1a,
	0xb1, 0xec, 0x39, 0xee, 0x02, 0x24, 0x31, 0x9d, 0x67, 0xa4, 0x4f, 0x05, 0xa3, 0xb4, 0x07, 0x50,
	0x76, 0x93, 0x28, 0xa2, 0x21, 0x77, 0xf6, 0xae, 0x94, 0x6f, 0x15, 0xba, 0x77, 0xd2, 0x09, 0x84,
	0x64, 0x9f, 0xb8, 0xcb, 0x8c, 0xac, 0x75, 0x2d, 0xad, 0xd2, 0x2c, 0x7e, 0x0e, 0xb7, 0x88, 0xef,
	0x3b, 0xea, 0x9c, 0x5b, 0xbb, 0x24, 0x90, 0x7e, 0xf3, 0xb2, 0x6d, 0x95, 0xf8, 0xbe, 0x25, 0xe2,
	0xba, 0xa8, 0x27, 0xa3, 0xe6, 0x37, 0x04, 0x37, 0x0f, 0x76, 0xc1, 0x0f, 0x0f, 0xb8, 0x53, 0x23,
	0xfc, 0xd3, 0xbf, 0x0d, 0x46, 0xcc, 0x09, 0x57, 0xfb, 0x2f, 0x75, 0x6b, 0x7b, 0xee, 0xb7, 0x8f,
	0x43, 0xc4, 0x2d, 0x95, 0x26, 0x6f, 0xbe, 0x7b, 0x74, 0x86, 0x95, 0xa2, 0x47, 0x1f, 0xb3, 0x4b,
	0x94, 0x35, 0xb8, 0x0e, 0xd5, 0x8b, 0x77, 0xfd, 0x89, 0x3d, 0x3a, 0x1f, 0x3b, 0xa3, 0xf1, 0xfb,
	0xde, 0xeb, 0xd1, 0x2b, 0x67, 0x62, 0xf7, 0xec, 0x7e, 0xf9, 0x08, 0x63, 0x28, 0x6d, 0x63, 0xd6,
	0x68, 0x30, 0xb4, 0xcb, 0x28, 0xc3, 0x4d, 0xad, 0xf3, 0xf1, 0xa0, 0x9c, 0xc3, 0x15, 0x28, 0xee,
	0xb8, 0xde, 0xc8, 0x2e, 0xe7, 0x5f, 0x34, 0x7f, 0x6e, 0x1a, 0xe8, 0xd7, 0xa6, 0x81, 0x7e, 0x6f,
	0x1a, 0xe8, 0xfb, 0x9f, 0xc6, 0xd1, 0x10, 0x7d, 0x38, 0xdb, 0xff, 0xd3, 0x66, 0xc7, 0xf2, 0xf3,
	0xf8, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xde, 0xed, 0xb2, 0x46, 0xde, 0x03, 0x00, 0x00,
}
