package data

//gogen:config
type FamilyName struct {
	_          struct{} `file:"武将名字/姓.txt"`
	FamilyName string   `protofield:"-" key:"1"`
}

//gogen:config
type MaleGivenName struct {
	_    struct{} `file:"武将名字/男名.txt"`
	Name string   `protofield:"-" key:"1"`
}

//gogen:config
type FemaleGivenName struct {
	_    struct{} `file:"武将名字/女名.txt"`
	Name string   `protofield:"-" key:"1"`
}
