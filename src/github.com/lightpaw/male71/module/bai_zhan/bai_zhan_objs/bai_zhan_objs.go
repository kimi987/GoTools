package bai_zhan_objs

func NewBaiZhanObjs() *BaiZhanObjs {
	return &BaiZhanObjs{allBaiZhanObjs: Newbai_zhan_obj_map()}
}

type BaiZhanObjs struct {
	allBaiZhanObjs *bai_zhan_obj_map
}

func (objs *BaiZhanObjs) Count() int {
	return objs.allBaiZhanObjs.Count()
}

func (objs *BaiZhanObjs) GetBaiZhanObj(id int64) (obj *HeroBaiZhanObj) {
	obj, _ = objs.allBaiZhanObjs.Get(id)
	return
}

func (objs *BaiZhanObjs) AddBaiZhanObj(obj *HeroBaiZhanObj) {
	objs.allBaiZhanObjs.Set(obj.Id(), obj)
}

func (objs *BaiZhanObjs) Walk(f func(obj *HeroBaiZhanObj)) {
	objs.allBaiZhanObjs.IterCb(func(key int64, v *HeroBaiZhanObj) {
		f(v)
	})
}

type BaiZhanObjsFunc func(objs *BaiZhanObjs)
