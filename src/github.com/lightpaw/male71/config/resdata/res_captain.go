package resdata

//gogen:config
type ResCaptainData struct {
	_ struct{} `file:"武将/武将.txt"`

	Id uint64

	object interface{}
}

func (d *ResCaptainData) GetObject() interface{} {
	return d.object
}

func (d *ResCaptainData) InitObject(obj interface{}) {
	d.object = obj
}
