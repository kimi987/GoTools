package rank_data

// 排行榜其他数据
//gogen:config
type RankMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"杂项/排行榜杂项.txt"`
	_ struct{} `proto:"shared_proto.RankMiscProto"`
	_ struct{} `protoconfig:"RankMisc"`

	RankCountPerPage uint64 `default:"5"`                    // 每页的数量
	MaxRankCount     uint64 `default:"10000" protofield:"-"` // 最大的排行人数
}
