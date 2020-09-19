package guild_data

//gogen:config
type GuildClassTitleData struct {
	_ struct{} `file:"联盟/联盟职称.txt"`
	_ struct{} `proto:"shared_proto.GuildClassTitleDataProto"`
	_ struct{} `protoconfig:"GuildClassTitle"`

	Id uint64

	Name string

	Permission *GuildPermissionData `type:"sub"`
}
