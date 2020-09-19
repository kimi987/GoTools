package nsds
//
//import (
//	"cloud.google.com/go/datastore"
//	"context"
//	"fmt"
//	. "github.com/onsi/gomega"
//	"testing"
//)
//
//const kind = "kind"
//
//func TestNamespace(t *testing.T) {
//	RegisterTestingT(t)
//
//	ns := "nsc1"
//	nsc := getClient(ns)
//	if nsc == nil {
//		return
//	}
//
//	// 设置ns
//	Ω(nsc.NamespaceKey(datastore.IDKey(kind, 1, nil)).Namespace).Should(Equal(ns))
//	Ω(nsc.NamespaceKey(datastore.NameKey(kind, "name", nil)).Namespace).Should(Equal(ns))
//	Ω(nsc.NamespaceKey(datastore.IncompleteKey(kind, nil)).Namespace).Should(Equal(ns))
//
//	// 自动分配ns
//	keys, err := nsc.AllocateIDs(context.TODO(), []*datastore.Key{
//		datastore.IncompleteKey(kind, nil),
//		nsc.NamespaceKey(datastore.IncompleteKey(kind, nil)),
//	})
//	Ω(err).Should(Succeed())
//
//	for _, key := range keys {
//		Ω(key.ID).ShouldNot(Equal(int64(0)))
//		Ω(key.Namespace).Should(Equal(ns))
//
//		err := nsc.Get(context.TODO(), key, &struct{}{})
//		Ω(err).Should(Equal(datastore.ErrNoSuchEntity))
//	}
//
//	err = nsc.DeleteMulti(context.TODO(), keys)
//	Ω(err).Should(Succeed())
//
//}
//
//type entity struct {
//	Int int
//}
//
//func TestPut(t *testing.T) {
//	RegisterTestingT(t)
//
//	ns1 := "nsc1"
//	nsc1 := getClient(ns1)
//	if nsc1 == nil {
//		fmt.Println("nsc1 == nil")
//		return
//	}
//
//	key := datastore.IDKey(kind, 1, nil)
//	for i := 0; i <= 10; i++ {
//		_, err := nsc1.Put(context.TODO(), key, &entity{
//			Int: i,
//		})
//		Ω(err).Should(Succeed())
//	}
//
//	var x entity
//	err := nsc1.Get(context.TODO(), key, &x)
//	Ω(err).Should(Succeed())
//	Ω(x.Int).Should(Equal(10))
//
//	err = nsc1.Delete(context.TODO(), key)
//	Ω(err).Should(Succeed())
//
//	// 保存多个
//	keys := []*datastore.Key{
//		datastore.IDKey(kind, 1, nil),
//		datastore.IDKey(kind, 2, nil),
//		datastore.IDKey(kind, 3, nil),
//	}
//	_, err = nsc1.PutMulti(context.TODO(), keys, []*entity{
//		{Int: 1}, {Int: 2}, {Int: 3},
//	})
//	Ω(err).Should(Succeed())
//
//	err = nsc1.DeleteMulti(context.TODO(), keys)
//	Ω(err).Should(Succeed())
//}
//
//func TestGet(t *testing.T) {
//	RegisterTestingT(t)
//
//	key := datastore.IDKey(kind, 1, nil)
//
//	var clean1 func()
//	func() {
//		ns1 := "nsc1"
//		nsc1 := getClient(ns1)
//		if nsc1 == nil {
//			fmt.Println("nsc1 == nil")
//			return
//		}
//
//		x := &entity{}
//		_, err := nsc1.Put(context.TODO(), key, &entity{
//			Int: 1,
//		})
//		Ω(err).Should(Succeed())
//
//		Ω(nsc1.Get(context.TODO(), key, x)).Should(Succeed())
//		Ω(x.Int).Should(Equal(1))
//
//		clean1 = func() {
//			count, err := nsc1.Count(context.TODO(), datastore.NewQuery(kind))
//			Ω(err).Should(Succeed())
//			Ω(count).Should(Equal(1))
//
//			nsc1.Delete(context.TODO(), key)
//			Ω(nsc1.Get(context.TODO(), key, x)).Should(Equal(datastore.ErrNoSuchEntity))
//		}
//	}()
//
//	var clean2 func()
//	func() {
//		ns2 := "nsc2"
//		nsc2 := getClient(ns2)
//		if nsc2 == nil {
//			fmt.Println("nsc2 == nil")
//			return
//		}
//
//		x := &entity{}
//		_, err := nsc2.Put(context.TODO(), key, &entity{
//			Int: 2,
//		})
//		Ω(err).Should(Succeed())
//
//		Ω(nsc2.Get(context.TODO(), key, x)).Should(Succeed())
//		Ω(x.Int).Should(Equal(2))
//
//		clean2 = func() {
//			count, err := nsc2.Count(context.TODO(), datastore.NewQuery(kind))
//			Ω(err).Should(Succeed())
//			Ω(count).Should(Equal(1))
//
//			nsc2.Delete(context.TODO(), key)
//			Ω(nsc2.Get(context.TODO(), key, x)).Should(Equal(datastore.ErrNoSuchEntity))
//		}
//	}()
//
//	if clean1 != nil {
//		clean1()
//	}
//
//	if clean2 != nil {
//		clean2()
//	}
//}
//
//func TestQuery(t *testing.T) {
//	RegisterTestingT(t)
//
//	ns := "nsc1"
//	nsc := getClient(ns)
//	if nsc == nil {
//		return
//	}
//
//	// 自动分配ns
//	keys, err := nsc.AllocateIDs(context.TODO(), []*datastore.Key{
//		datastore.IncompleteKey(kind, nil),
//		nsc.NamespaceKey(datastore.IncompleteKey(kind, nil)),
//	})
//	Ω(err).Should(Succeed())
//
//	for _, key := range keys {
//		Ω(key.ID).ShouldNot(Equal(int64(0)))
//		Ω(key.Namespace).Should(Equal(ns))
//
//		// 不存在
//		err := nsc.Get(context.TODO(), key, &struct{}{})
//		Ω(err).Should(Equal(datastore.ErrNoSuchEntity))
//
//		// put值
//		_, err = nsc.Put(context.TODO(), key, &entity{
//			Int: int(key.ID),
//		})
//		Ω(err).Should(Succeed())
//	}
//
//	// 查询一下
//	var xs []*entity
//	keys2, err := nsc.GetAll(context.TODO(), datastore.NewQuery(kind), &xs)
//	Ω(err).Should(Succeed())
//	Ω(keys2).Should(ConsistOf(keys))
//
//	for i, k := range keys2 {
//		Ω(xs[i].Int).Should(Equal(int(k.ID)))
//	}
//
//	err = nsc.DeleteMulti(context.TODO(), keys)
//	Ω(err).Should(Succeed())
//}
//
//func TestTransaction(t *testing.T) {
//	RegisterTestingT(t)
//
//	key := datastore.IDKey(kind, 1, nil)
//
//	var clean1 func()
//	func() {
//		ns1 := "nsc1"
//		nsc1 := getClient(ns1)
//		if nsc1 == nil {
//			fmt.Println("nsc1 == nil")
//			return
//		}
//
//		x := &entity{}
//
//		_, err := nsc1.Put(context.TODO(), key, &entity{
//			Int: 1,
//		})
//		Ω(err).Should(Succeed())
//
//		nsc1.RunInTransaction(context.TODO(), func(tx *Transaction) error {
//			Ω(tx.Get(key, x)).Should(Succeed())
//			Ω(x.Int).Should(Equal(1))
//
//			return nil
//		})
//
//		Ω(nsc1.Get(context.TODO(), key, x)).Should(Succeed())
//		Ω(x.Int).Should(Equal(1))
//
//		clean1 = func() {
//			count, err := nsc1.Count(context.TODO(), datastore.NewQuery(kind))
//			Ω(err).Should(Succeed())
//			Ω(count).Should(Equal(1))
//
//			nsc1.RunInTransaction(context.TODO(), func(tx *Transaction) error {
//				tx.Delete(key)
//				return nil
//			})
//
//			Ω(nsc1.Get(context.TODO(), key, x)).Should(Equal(datastore.ErrNoSuchEntity))
//		}
//	}()
//
//	if clean1 != nil {
//		clean1()
//	}
//}
//
//func getClient(namespace string) *NamespaceClient {
//	projectID := "male7"
//
//	client, err := NewNamespaceClient("192.168.1.5:8432", projectID, namespace)
//	if err != nil {
//		fmt.Printf("Failed to create client: %v", err)
//		return nil
//	}
//
//	keys, err := client.GetAll(context.TODO(), datastore.NewQuery(kind).KeysOnly(), nil)
//	Ω(err).Should(Succeed())
//
//	err = client.DeleteMulti(context.TODO(), keys)
//	Ω(err).Should(Succeed())
//
//	Eventually(func() int {
//		keys, err := client.GetAll(context.TODO(), datastore.NewQuery(kind).KeysOnly(), nil)
//		Ω(err).Should(Succeed())
//
//		return len(keys)
//	}).Should(Equal(0))
//
//	return client
//}
