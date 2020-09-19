package collection

func NewRingList(capcity int) *RingList {
	return &RingList{
		Queue:   NewQueue(),
		capcity: capcity,
	}
}

type RingList struct {
	*Queue

	capcity int
}

func (r *RingList) Capcity() int {
	return r.capcity
}

func (r *RingList) SetCapcity(toSet int) {
	r.capcity = toSet
}

// Add puts an element on the end of the queue.
func (r *RingList) Add(elem interface{}) {
	if r.capcity <= 0 {
		return
	}

	if r.Length() >= r.capcity {
		r.Queue.Remove()
	}

	r.Queue.Add(elem)
}
