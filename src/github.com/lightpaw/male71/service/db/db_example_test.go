package db

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestQuery(t *testing.T) {
	RegisterTestingT(t)

	//datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	//Ω(err).Should(Succeed())
	//
	//sid := 40
	//client, err := nsds.NewNamespaceClient("192.168.1.5:8432", "male7", fmt.Sprintf("1-%d", sid))
	//if err != nil {
	//	fmt.Println("创建nsds.Client失败，跳过测试", err.Error())
	//	return
	//}

	//db, err := datastoredb.NewDatastoreDbServiceWithNsdsClient(datas, timeservice.NewDefaultTimeService(), client)
	//Ω(err).Should(Succeed())
	//
	//hero, err := db.LoadHero(context.TODO(), 5384)
	//Ω(err).Should(Succeed())iface_gen
	//
	//if hero == nil {
	//	fmt.Println("hero not exist")
	//} else {
	//	fmt.Println(hero.Id(), hero.Name())
	//}

}
