package entity

import (
	"github.com/lightpaw/male7/config/heroinit"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"math/rand"
	"time"
)

// 玩家成就
func newHeroTag(initData *heroinit.HeroInitData) *HeroTag {
	return &HeroTag{
		maxTagColorType:     u64.Int32(initData.MaxTagColorType),
		maxTagRecordCount:   initData.MaxTagRecordCount,
		maxShowForViewCount: initData.MaxShowForViewCount,
		tags:                make(map[string]*shared_proto.TagProto),
	}
}

// 玩家标签
type HeroTag struct {
	// 标签颜色类型数量
	maxTagColorType int32
	// 最大记录的标签日志数量
	maxTagRecordCount uint64
	// 展示给查看的标签数量
	maxShowForViewCount uint64
	// 标签数组
	tags map[string]*shared_proto.TagProto
	// 标签记录数组
	records []*shared_proto.TagRecordProto
}

func (h *HeroTag) TagCount() int {
	return len(h.tags)
}

func (h *HeroTag) RemoveTag(content string) (suc bool) {
	t := h.tags[content]
	if t == nil {
		return false
	}

	delete(h.tags, content)

	j := 0
	for _, record := range h.records {
		if record.Tag != content {
			// 相同内容，干掉
			h.records[j] = record
			j++
		}
	}

	h.records = h.records[:j]

	return true
}

func (h *HeroTag) Exist(tag string) bool {
	return h.tags[tag] != nil
}

func (h *HeroTag) AddTag(id []byte, name, flagName, tag string, time time.Time) (t *shared_proto.TagProto, record *shared_proto.TagRecordProto) {
	t = h.tags[tag]
	if t != nil {
		t.Count++
	} else {
		t = &shared_proto.TagProto{
			Tag:      tag,
			Count:    1,
			TagColor: rand.Int31n(h.maxTagColorType) + 1,
		}
		h.tags[tag] = t
	}

	record = &shared_proto.TagRecordProto{
		Id:       id,
		Name:     name,
		FlagName: flagName,
		Tag:      tag,
		Time:     timeutil.Marshal32(time),
	}
	h.addTagRecord(record)

	return
}

func (h *HeroTag) addTagRecord(proto *shared_proto.TagRecordProto) {
	if len(h.records) >= int(h.maxTagRecordCount) {
		// 删掉最左侧的
		copy(h.records[0:], h.records[1:])
		h.records[len(h.records)-1] = proto
	} else {
		h.records = append(h.records, proto)
	}
}

func (h *HeroTag) Encode() *shared_proto.HeroTagProto {
	proto := &shared_proto.HeroTagProto{}

	proto.Tags = make([]*shared_proto.TagProto, 0, len(h.tags))

	for _, t := range h.tags {
		proto.Tags = append(proto.Tags, t)
	}

	proto.Records = h.records

	return proto
}

func (h *HeroTag) unmarshal(proto *shared_proto.HeroTagProto) {
	if proto == nil {
		return
	}

	for _, t := range proto.Tags {
		h.tags[t.Tag] = t
	}
	h.records = proto.GetRecords()
}
