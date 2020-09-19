package lock

import (
	"context"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestLockService_Lock(t *testing.T) {
	RegisterTestingT(t)
	service := NewLockService(&provider{}, 10*time.Second)

	obj1, unlock, err := service.Lock(123)
	Ω(err).Should(Succeed())
	obj1.(*lo).changed = 23
	unlock()

	time.Sleep(1 * time.Second)
	service.doSaveLoop(500 * time.Millisecond)
	Ω(obj1.(*lo).saved).Should(BeNumerically("==", 23))

	service.doEvictLoop(3*time.Second, 0)
	Ω(service.entries.Has(123)).Should(BeTrue())

	service.doEvictLoop(500*time.Millisecond, 0)
	Ω(service.entries.Has(123)).Should(BeFalse())
}

func TestLockService_EvictMustSave(t *testing.T) {
	RegisterTestingT(t)
	service := NewLockService(&provider{}, 10*time.Second)

	obj1, unlock, err := service.Lock(123)
	Ω(err).Should(Succeed())
	obj1.(*lo).changed = 23
	unlock()

	time.Sleep(1 * time.Second)

	service.doEvictLoop(500*time.Millisecond, 0)
	Ω(service.entries.Has(123)).Should(BeFalse())

	Ω(obj1.(*lo).saved).Should(BeNumerically("==", 23), "evict must be saved first")
}

func TestLockService_Check_Bug(t *testing.T) {
	RegisterTestingT(t)
	service := NewLockService(&provider{}, 10*time.Second)

	obj1, _, err := service.Lock(123)
	Ω(err).Should(Succeed())
	obj1.(*lo).changed = 23

	//service.doCheckLoop()
}

func TestCreate(t *testing.T) {
	RegisterTestingT(t)
	service := NewLockService(&provider{}, 10*time.Second)

	id := int64(110)
	o := &lo{id: id}
	unlock, err := service.Create(id, o)
	Ω(err).Should(Succeed())
	Ω(o.created).Should(Equal(true))
	o.changed = 23
	unlock()

	obj1, unlock, err := service.Lock(id)
	Ω(err).Should(Succeed())
	Ω(obj1.(*lo).changed).Should(Equal(int64(23)))
	unlock()
}

// ----------- setup -------------
type provider struct {
}

func (provider) Name() string {
	return "testService"
}

func (provider) GetObject(ctx context.Context, id int64) (LockObject, error) {
	if id == 110 {
		return nil, ErrEmpty
	}
	return &lo{id: id}, nil
}

func (provider) SaveObject(ctx context.Context, id int64, obj LockObject) error {
	o := obj.(*lo)
	o.saved = o.changed
	return nil
}

func (provider) CreateObject(ctx context.Context, obj LockObject) error {
	o := obj.(*lo)
	o.created = true
	return nil
}

type lo struct {
	id      int64
	changed int64
	saved   int64
	created bool
}

func (l lo) Marshal() ([]byte, error) {
	return nil, nil
}
