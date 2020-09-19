package resdata

// cost 兑换 prize

//gogen:config
type ExchangeData struct {
	_ struct{} `proto:"shared_proto.ExchangeDataProto"`

	Cost  *Cost  // 消耗
	Prize *Prize // 奖励
}
