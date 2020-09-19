package tag

// 打标签配置
//gogen:config
type TagMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"杂项/标签杂项.txt"`
	_ struct{} `proto:"shared_proto.TagMiscProto"`
	_ struct{} `protoconfig:"TagMisc"`

	MaxCount            uint64 `validator:"int>0" default:"50"`               // 最大的标签数量
	MaxCharCount        uint64 `validator:"int>0" default:"5"`                // 标签最多的字的数量, >就不允许了，一个汉字算一个字，一个英文字母算一个字
	MaxRecordCount      uint64 `validator:"int>0" default:"50"`               // 最大记录的标签日志数量
	MaxShowForViewCount uint64 `validator:"int>0" default:"5"`                // 展示给查看的标签数量
	MaxTagColorType     uint64 `validator:"int>0" default:"4" protofield:"-"` // 标签最多有多少种颜色
}
