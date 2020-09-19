package datastoredb
//
//import (
//	"cloud.google.com/go/datastore"
//	"context"
//	"fmt"
//	"github.com/lightpaw/logrus"
//	"github.com/lightpaw/male7/config"
//	"github.com/lightpaw/male7/config/kv"
//	"github.com/lightpaw/male7/config/resdata"
//	"github.com/lightpaw/male7/entity"
//	"github.com/lightpaw/male7/gen/pb/util"
//	"github.com/lightpaw/male7/pb/server_proto"
//	"github.com/lightpaw/male7/pb/shared_proto"
//	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
//	"github.com/lightpaw/male7/service/db/isql"
//	"github.com/lightpaw/male7/service/db/nsds"
//	"github.com/lightpaw/male7/service/timeservice"
//	"github.com/lightpaw/male7/util/atomic"
//	"github.com/lightpaw/male7/util/must"
//	"github.com/lightpaw/male7/util/u64"
//	"github.com/pkg/errors"
//	"google.golang.org/api/iterator"
//	"sync"
//	"github.com/lightpaw/male7/entity/cb"
//	"github.com/lightpaw/male7/util/ctxfunc"
//)
//
//func NewDatastoreDbService(datas *config.ConfigDatas, serverConfig *kv.IndividualServerConfig, timeService *timeservice.TimeService) (*DatastoreDbService, error) {
//
//	// namespace 使用平台id+区服id来区分
//	namespace := fmt.Sprintf("%d-%d", serverConfig.PlatformID, serverConfig.ServerID)
//
//	client, err := nsds.NewNamespaceClient(serverConfig.DatastoreHost, serverConfig.ProjectID, namespace)
//	if err != nil {
//		return nil, errors.Wrapf(err, "NewDatastoreDbService 错误")
//	}
//
//	return NewDatastoreDbServiceWithNsdsClient(datas, timeService, client)
//}
//
//func NewDatastoreDbServiceWithNsdsClient(datas *config.ConfigDatas, timeService *timeservice.TimeService, client *nsds.NamespaceClient) (*DatastoreDbService, error) {
//
//	db := &DatastoreDbService{
//		datas:       datas,
//		timeService: timeService,
//		client:      client,
//	}
//
//	var maxId uint64
//	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
//		maxId, err = db.maxGuildLogId(ctx)
//		return
//	})
//	if err != nil {
//		return nil, errors.Wrapf(err, "db.maxGuildLogId() 错误")
//	}
//	db.guildLogIdGen = atomic.NewUint64(maxId)
//
//	return db, nil
//}
//
//type DatastoreDbService struct {
//	datas       *config.ConfigDatas
//	timeService *timeservice.TimeService
//
//	client *nsds.NamespaceClient
//
//	guildLogIdGen *atomic.Uint64
//}
//
//func (d *DatastoreDbService) Close() error {
//	return nil
//}
//
//const (
//	kv_kind         = "kv"
//	user_kind       = "user"
//	heroname_kind   = "heroname"
//	hero_kind       = "hero"
//	guild_kind      = "guild"
//	guild_logs_kind = "guild_logs"
//	mail_kind       = "mail"
//	baizhan_kind    = "baizhan"
//	farm_kind       = "farm"
//)
//
//type config_entity struct {
//	V []byte `datastore:",noindex"`
//}
//
//func (d *DatastoreDbService) LoadKey(ctx context.Context, key server_proto.Key) ([]byte, error) {
//
//	keyName, ok := server_proto.Key_name[int32(key)]
//	if !ok {
//		return nil, errors.Errorf("datastore.LoadKey(%v) key不存在", key)
//	}
//
//	k := datastore.NameKey(kv_kind, keyName, nil)
//	v := config_entity{}
//	err := d.client.Get(ctx, k, &v)
//	if err != nil && err != datastore.ErrNoSuchEntity {
//		return nil, errors.Wrapf(err, "datastore.LoadKey(%s) 失败", keyName)
//	}
//
//	return v.V, nil
//}
//
//func (d *DatastoreDbService) SaveKey(ctx context.Context, key server_proto.Key, data []byte) error {
//
//	keyName, ok := server_proto.Key_name[int32(key)]
//	if !ok {
//		return errors.Errorf("datastore.SaveKey(%v) key不存在", key)
//	}
//
//	k := datastore.NameKey(kv_kind, keyName, nil)
//	_, err := d.client.Put(ctx, k, &config_entity{
//		V: data,
//	})
//	if err != nil {
//		return errors.Wrapf(err, "datastore.SaveKey() fail")
//	}
//
//	return nil
//}
//
//type user_entity struct {
//	Misc []byte `datastore:",noindex"`
//}
//
//func (d *DatastoreDbService) UpdateUserMisc(ctx context.Context, id int64, proto *server_proto.UserMiscProto) error {
//	data, err := proto.Marshal()
//	if err != nil {
//		return err
//	}
//
//	_, err = d.client.Put(ctx, datastore.IDKey(user_kind, id, nil), &user_entity{
//		Misc: data,
//	})
//	if err != nil {
//		return errors.Wrapf(err, "datastore.UpdateUserMisc() 失败: %d, %#v", id, proto)
//	}
//
//	return nil
//}
//
//func (d *DatastoreDbService) LoadUserMisc(ctx context.Context, id int64) (*server_proto.UserMiscProto, error) {
//	entity := user_entity{}
//	err := d.client.Get(ctx, datastore.IDKey(user_kind, id, nil), &entity)
//	if err != nil {
//		if err == datastore.ErrNoSuchEntity {
//			return &server_proto.UserMiscProto{}, nil
//		}
//
//		return nil, errors.Wrapf(err, "DatastoreDbService.LoadUserMisc(%d)", id)
//	}
//
//	proto := &server_proto.UserMiscProto{}
//	if len(entity.Misc) > 0 {
//		if err := proto.Unmarshal(entity.Misc); err != nil {
//			return nil, err
//		}
//	}
//
//	return proto, nil
//}
//
//func (d *DatastoreDbService) UpdateSettings(ctx context.Context, id int64, settings uint64) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) FindSettingsOpen(ctx context.Context, settingType shared_proto.SettingType, ids []int64) (result []int64, err error) {
//	err = ErrUnsupport
//	return
//}
//
//type heroname_entity struct {
//	Id int64 `datastore:",noindex"`
//}
//
//func (d *DatastoreDbService) HeroNameExist(ctx context.Context, name string) (bool, error) {
//
//	heroId, err := d.HeroId(ctx, name)
//	if err != nil {
//		return false, errors.Wrapf(err, "DB查询英雄名字是否存在失败")
//	}
//
//	return heroId != 0, nil
//
//}
//
//func (d *DatastoreDbService) HeroId(ctx context.Context, name string) (heroId int64, err error) {
//	if len(name) <= 0 {
//		return 0, errors.Errorf("根据英雄名字查询id，名字为空")
//	}
//
//	dst := heroname_entity{}
//	err = d.client.Get(ctx, datastore.NameKey(heroname_kind, name, nil), &dst)
//	if err != nil {
//		if err == datastore.ErrNoSuchEntity {
//			return 0, nil
//		}
//		return 0, errors.Wrapf(err, "DB查询英雄名字对应的玩家id失败")
//	}
//
//	return dst.Id, nil
//}
//
//var errHeroNameExist = errors.Errorf("英雄名字已经存在")
//
//func (d *DatastoreDbService) setHeroNameIfAbsent(id int64, newName string, tx *nsds.Transaction) error {
//	entity := heroname_entity{}
//	newKey := datastore.NameKey(heroname_kind, newName, nil)
//
//	err := tx.Get(newKey, &entity)
//	switch err {
//	case nil:
//		return errHeroNameExist
//	case datastore.ErrNoSuchEntity:
//		// 不存在
//	default:
//		return errors.Wrapf(err, "查询英雄名字是否存在失败")
//	}
//
//	// 到这里，新的那个名字没有被使用，尝试插入新的值
//	entity.Id = id
//	if _, err := tx.Put(newKey, &entity); err != nil {
//		return errors.Wrapf(err, "英雄改名，Put失败")
//	}
//
//	return nil
//}
//
//func (d *DatastoreDbService) UpdateHeroName(ctx context.Context, id int64, originName, newName string) bool {
//
//	if len(newName) <= 0 {
//		return false
//	}
//
//	_, err := d.client.RunInTransaction(ctx, func(tx *nsds.Transaction) error {
//		return d.setHeroNameIfAbsent(id, newName, tx)
//	})
//
//	if err != nil {
//		logrus.WithError(err).Debugf("英雄改名")
//		return false
//	}
//
//	// 尝试删除老的值
//	if len(originName) > 0 {
//		_, err := d.client.RunInTransaction(ctx, func(tx *nsds.Transaction) error {
//			originKey := datastore.NameKey(heroname_kind, originName, nil)
//			dst := heroname_entity{}
//			if err := tx.Get(originKey, &dst); err != nil {
//				if err != datastore.ErrNoSuchEntity {
//					return nil
//				}
//				return err
//			}
//
//			if dst.Id == id {
//				if err := tx.Delete(originKey); err != nil {
//					return err
//				}
//			}
//
//			return nil
//		})
//		if err != nil {
//			logrus.WithError(err).Debugf("英雄改名，删除原来的名字失败")
//		}
//	}
//
//	return true
//}
//
//func (d *DatastoreDbService) LoadNoGuildHeroListByName(ctx context.Context, text string, index, size uint64) (heros []*entity.Hero, err error) {
//	return []*entity.Hero{}, ErrUnsupport
//}
//
//func (d *DatastoreDbService) UpdateHeroGuildId(ctx context.Context, id, guildId int64) error {
//	return ErrUnsupport
//}
//
//type hero_entity struct {
//	Name       string `datastore:",noindex"`
//	HeroData   []byte `datastore:",noindex"`
//	BaseRegion int64
//}
//
//var errHeroExist = errors.Errorf("创建英雄，英雄已经存在")
//
//func (d *DatastoreDbService) CreateHero(ctx context.Context, hero *entity.Hero) (bool, error) {
//
//	heroData, err := hero.EncodeServer().Marshal()
//	if err != nil {
//		return false, errors.Wrapf(err, "CreateHero ServerProto Marshal heroData fail")
//	}
//
//	heroId := hero.Id()
//	heroName := hero.Name()
//	baseRegion := hero.BaseRegion()
//
//	_, err = d.client.RunInTransaction(ctx, func(tx *nsds.Transaction) error {
//
//		// 查询这个玩家是否创建了角色
//		heroKey := datastore.IDKey(hero_kind, heroId, nil)
//		if err := tx.Get(heroKey, &hero_entity{}); err != datastore.ErrNoSuchEntity {
//			if err == nil {
//				// 角色已经存在
//				err = errHeroExist
//			}
//			return err
//		}
//
//		if err := d.setHeroNameIfAbsent(heroId, heroName, tx); err != nil {
//			// 名字已经存在
//			return err
//		}
//
//		if _, err := tx.Put(heroKey, &hero_entity{
//			Name:       heroName,
//			HeroData:   heroData,
//			BaseRegion: baseRegion,
//		}); err != nil {
//			return err
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		return false, errors.Wrapf(err, "datastore.CreateHero() 失败")
//	}
//
//	return true, nil
//}
//
//func (d *DatastoreDbService) SaveHero(ctx context.Context, hero *entity.Hero) error {
//
//	heroData, err := hero.EncodeServer().Marshal()
//	if err != nil {
//		return errors.Wrapf(err, "SaveHero Proto Marshal fail")
//	}
//
//	heroId := hero.Id()
//	heroName := hero.Name()
//	baseRegion := hero.BaseRegion()
//
//	heroKey := datastore.IDKey(hero_kind, heroId, nil)
//
//	if _, err := d.client.Put(ctx, heroKey, &hero_entity{
//		Name:       heroName,
//		HeroData:   heroData,
//		BaseRegion: baseRegion,
//	}); err != nil {
//		return errors.Wrapf(err, "datastore.SaveHero() 失败")
//	}
//
//	return nil
//}
//
//func (d *DatastoreDbService) LoadHero(ctx context.Context, id int64) (*entity.Hero, error) {
//
//	heroKey := datastore.IDKey(hero_kind, id, nil)
//
//	entity := hero_entity{}
//	if err := d.client.Get(ctx, heroKey, &entity); err != nil {
//		if err == datastore.ErrNoSuchEntity {
//			return nil, nil
//		}
//
//		return nil, errors.Wrapf(err, "DatastoreDbService.LoadHero(%d)", id)
//	}
//
//	return d.parseHero(id, entity.Name, entity.HeroData)
//}
//
//func (d *DatastoreDbService) parseHero(id int64, name string, heroData []byte) (*entity.Hero, error) {
//	var heroProto *server_proto.HeroServerProto
//	if len(heroData) > 0 {
//		heroProto = &server_proto.HeroServerProto{}
//		err := heroProto.Unmarshal(heroData)
//		if err != nil {
//			return nil, errors.Wrapf(err, "DbSerbive.parseHero(%d), Unmarshal HeroServerProto fail", id)
//		}
//	}
//
//	hero := entity.UnmarshalHero(id, name, d.datas.HeroInitData(), heroProto, d.datas, d.timeService.CurrentTime())
//
//	return hero, nil
//}
//
//func (d *DatastoreDbService) LoadHeroCount(ctx context.Context) (uint64, error) {
//	q := datastore.NewQuery(hero_kind)
//	count, err := d.client.Count(ctx, q)
//	return u64.FromInt(count), err
//}
//
//func (d *DatastoreDbService) LoadAllHeroData(ctx context.Context) ([]*entity.Hero, error) {
//
//	q := datastore.NewQuery(hero_kind)
//
//	array := make([]*entity.Hero, 0)
//	for t := d.client.Run(ctx, q); ; {
//		var x hero_entity
//		key, err := t.Next(&x)
//		if err != nil {
//			if err == iterator.Done {
//				break
//			}
//
//			return nil, errors.Wrapf(err, "datastore.LoadAllHeroData() 失败")
//		}
//
//		hero, err := d.parseHero(key.ID, x.Name, x.HeroData)
//		if err != nil {
//			return nil, errors.Wrapf(err, "datastore.LoadAllHeroData() 失败")
//		}
//
//		array = append(array, hero)
//	}
//
//	return array, nil
//}
//
//func (d *DatastoreDbService) LoadAllRegionHero(ctx context.Context) ([]*entity.Hero, error) {
//
//	q := datastore.NewQuery(hero_kind).Filter("BaseRegion >", 0)
//
//	array := make([]*entity.Hero, 0)
//	for t := d.client.Run(ctx, q); ; {
//		var x hero_entity
//		key, err := t.Next(&x)
//		if err != nil {
//			if err == iterator.Done {
//				break
//			}
//
//			return nil, errors.Wrapf(err, "datastore.LoadAllRegionHero() 失败")
//		}
//
//		hero, err := d.parseHero(key.ID, x.Name, x.HeroData)
//		if err != nil {
//			return nil, errors.Wrapf(err, "datastore.LoadAllRegionHero() 失败")
//		}
//
//		array = append(array, hero)
//	}
//
//	return array, nil
//}
//
//// mail
//
//type mail_entity struct {
//	Receiver  int64
//	Data      []byte `datastore:",noindex"`
//	Keep      bool
//	Readed    bool
//	HasReport bool
//	ReportTag int32
//	HasPrize  bool
//	Collected bool
//}
//
//func (d *DatastoreDbService) MaxMailId(ctx context.Context) (uint64, error) {
//
//	q := datastore.NewQuery(mail_kind).Order("-__key__").Limit(1).KeysOnly()
//
//	key, err := d.getFirstKey(ctx, q)
//	if err != nil {
//		if err == datastore.ErrNoSuchEntity {
//			return 0, nil
//		}
//
//		return 0, errors.Wrapf(err, "获取最大的邮件id失败")
//	}
//
//	return uint64(key.ID), nil
//}
//
//func (d *DatastoreDbService) getFirstKey(ctx context.Context, q *datastore.Query) (*datastore.Key, error) {
//
//	t := d.client.Run(ctx, q)
//	key, err := t.Next(nil)
//	if err == iterator.Done {
//		err = datastore.ErrNoSuchEntity
//	}
//
//	return key, err
//}
//
//func boolFilter(q *datastore.Query, fieldname string, filter int32) *datastore.Query {
//
//	switch filter {
//	case 0:
//		return q
//	case 1:
//		return q.Filter(fieldname, false)
//	case 2:
//		return q.Filter(fieldname, true)
//	default:
//		logrus.Errorf("boolFilterString invalid filter: %v", filter)
//		return q
//	}
//}
//
//// 工程没有地方调用到此函数，是否考虑去除@Albert Fan
//func (d *DatastoreDbService) LoadHeroMailList(ctx context.Context, heroId int64, minMailId uint64, keep, readed, has_report, has_prize, collected int32, count uint64) ([]*shared_proto.MailProto, error) {
//
//	count = u64.Min(u64.Max(count, d.datas.MiscConfig().MailMinBatchCount), d.datas.MiscConfig().MailMaxBatchCount)
//
//	q := datastore.NewQuery(mail_kind)
//	if minMailId > 0 {
//		q = q.Filter("__key__ <", d.client.NamespaceKey(datastore.IDKey(mail_kind, int64(minMailId), nil)))
//	}
//
//	q = q.Filter("Receiver =", heroId)
//	q = boolFilter(q, "Keep =", keep)
//	q = boolFilter(q, "Readed =", readed)
//	q = boolFilter(q, "HasReport =", has_report)
//	q = boolFilter(q, "HasPrize =", has_prize)
//	q = boolFilter(q, "Collected =", collected)
//
//	q = q.Order("-__key__").Limit(int(count))
//
//	var xs []*mail_entity
//	_, err := d.client.GetAll(ctx, q, &xs)
//	if err != nil {
//		return nil, errors.Wrapf(err, "datastore.LoadHeroMailList() 失败")
//	}
//
//	var datas []*shared_proto.MailProto
//	for _, x := range xs {
//		mailProto := &shared_proto.MailProto{}
//		if err := mailProto.Unmarshal(x.Data); err != nil {
//			return nil, errors.Wrapf(err, "datastore.LoadHeroMailList() MailProto.Unmarshal失败")
//		}
//
//		mailProto.Keep = x.Keep
//		mailProto.Read = x.Readed
//		mailProto.HasPrize = mailProto.Prize != nil
//		mailProto.Collected = x.Collected
//
//		datas = append(datas, mailProto)
//	}
//
//	return datas, nil
//}
//
//func (d *DatastoreDbService) LoadMailCountHasPrizeNotCollected(ctx context.Context, heroId int64) (int, error) {
//	q := datastore.NewQuery(mail_kind).Filter("Receiver =", heroId).Filter("HasPrize =", true).Filter("Collected =", false)
//
//	return d.client.Count(ctx, q)
//}
//
//
//
//func (d *DatastoreDbService) LoadMailCountHasReportNotReaded(ctx context.Context, heroId int64, reportTag int32) (int, error) {
//	q := datastore.NewQuery(mail_kind).Filter("Receiver =", heroId).Filter("HasReport =", true).Filter("ReportTag =", reportTag).Filter("Readed =", false)
//
//	return d.client.Count(ctx, q)
//}
//
//func (d *DatastoreDbService) LoadMailCountNoReportNotReaded(ctx context.Context, heroId int64) (int, error) {
//	q := datastore.NewQuery(mail_kind).Filter("Receiver =", heroId).Filter("HasReport =", false).Filter("Readed =", false)
//
//	return d.client.Count(ctx, q)
//}
//
//func (d *DatastoreDbService) loadMailEntity(ctx context.Context, id uint64) (*datastore.Key, *mail_entity, error) {
//	key := datastore.IDKey(mail_kind, int64(id), nil)
//
//	var x mail_entity
//	err := d.client.Get(ctx, key, &x)
//	if err != nil {
//		if err == datastore.ErrNoSuchEntity {
//			return key, nil, nil
//		}
//		return nil, nil, err
//	}
//
//	return key, &x, nil
//}
//
//func (d *DatastoreDbService) IsCollectableMail(ctx context.Context, id uint64) (bool, error) {
//	_, entity, err := d.loadMailEntity(ctx, id)
//	if err != nil {
//		return false, errors.Wrapf(err, "datastore.LoadMailPrize() 失败")
//	}
//
//	return entity.HasPrize && !entity.Collected, nil
//}
//
//func (d *DatastoreDbService) LoadCollectMailPrize(ctx context.Context, id uint64, heroId int64) (*resdata.Prize, error) {
//
//	_, entity, err := d.loadMailEntity(ctx, id)
//	if err != nil {
//		return nil, errors.Wrapf(err, "datastore.LoadMailPrize() 失败")
//	}
//
//	if entity != nil && len(entity.Data) > 0 && entity.Receiver == heroId &&
//		entity.HasPrize && !entity.Collected {
//
//		proto := &shared_proto.MailProto{}
//		if err := proto.Unmarshal(entity.Data); err != nil {
//			return nil, errors.Wrapf(err, "datastore.LoadCollectMailPrize() MailProto.Unmarshal失败")
//		}
//
//		if proto.Prize != nil {
//			return resdata.UnmarshalPrize(proto.Prize, d.datas), nil
//		}
//	}
//
//	return nil, nil
//}
//
//func (d *DatastoreDbService) LoadMail(ctx context.Context, id uint64) ([]byte, error) {
//	_, entity, err := d.loadMailEntity(ctx, id)
//	if err != nil {
//		return nil, errors.Wrapf(err, "datastore.LoadMailPrize() 失败")
//	}
//
//	return entity.Data, nil
//}
//
//
//func (d *DatastoreDbService) CreateMail(ctx context.Context, id uint64, receiver int64, data []byte, keep, has_report, has_prize bool, report_tag int32, time int64) error {
//
//	readed := false
//	collected := false
//
//	key := datastore.IDKey(mail_kind, int64(id), nil)
//
//	_, err := d.client.RunInTransaction(ctx, func(tx *nsds.Transaction) error {
//
//		err := tx.Get(key, &mail_entity{})
//		if err != datastore.ErrNoSuchEntity {
//			if err != nil {
//				return err
//			}
//
//
//			return errors.Errorf("邮件已经存在, id: %d", id)
//		}
//
//		_, err = tx.Put(key, &mail_entity{
//			Receiver:  receiver,
//			Data:      data,
//			Keep:      keep,
//			Readed:    readed,
//			HasReport: has_report,
//			ReportTag: report_tag,
//			HasPrize:  has_prize,
//			Collected: collected,
//		})
//
//		return err
//	})
//
//	if err != nil {
//		return errors.Wrapf(err, "datastore.CreateMail() 失败")
//	}
//
//	return nil
//}
//
//func (d *DatastoreDbService) DeleteMail(ctx context.Context, id uint64, heroId int64) error {
//
//	err := d.deleteMailFunc(ctx, id, heroId, func(entity *mail_entity) (success bool) {
//		return true
//	})
//
//	if err != nil {
//		return errors.Wrapf(err, "datastore.DeleteMail() 失败")
//	}
//
//	return nil
//}
//
//func (d *DatastoreDbService) deleteMailFunc(ctx context.Context, id uint64, heroId int64, f func(entity *mail_entity) (success bool)) error {
//	key := datastore.IDKey(mail_kind, int64(id), nil)
//
//	_, err := d.client.RunInTransaction(ctx, func(tx *nsds.Transaction) error {
//
//		var x mail_entity
//		if err := tx.Get(key, &x); err != nil {
//			if err == datastore.ErrNoSuchEntity {
//				return nil
//			}
//			return err
//		}
//
//		if x.Receiver != heroId {
//			return errOtherMailFail
//		}
//
//		if !f(&x) {
//			// rollback
//			return errUpdateFail
//		}
//
//		if err := tx.Delete(key); err != nil {
//			return err
//		}
//
//		return nil
//	})
//
//	if err == errUpdateFail || err == errOtherMailFail {
//		return nil
//	}
//
//	return err
//}
//
//var errOtherMailFail = errors.Errorf("操作的是其他玩家的邮件")
//var errUpdateFail = errors.Errorf("更新错误")
//
//func (d *DatastoreDbService) updateMailFunc(ctx context.Context, id uint64, heroId int64, f func(entity *mail_entity) (success bool)) error {
//	key := datastore.IDKey(mail_kind, int64(id), nil)
//
//	_, err := d.client.RunInTransaction(ctx, func(tx *nsds.Transaction) error {
//
//		var x mail_entity
//		if err := tx.Get(key, &x); err != nil {
//			if err == datastore.ErrNoSuchEntity {
//				return nil
//			}
//			return err
//		}
//
//		if x.Receiver != heroId {
//			return errOtherMailFail
//		}
//
//		if !f(&x) {
//			// rollback
//			return errUpdateFail
//		}
//
//		if _, err := tx.Put(key, &x); err != nil {
//			return err
//		}
//
//		return nil
//	})
//
//	if err == errUpdateFail || err == errOtherMailFail {
//		return nil
//	}
//
//	return err
//}
//
//func (d *DatastoreDbService) UpdateMailKeep(ctx context.Context, id uint64, heroId int64, keep bool) error {
//
//	err := d.updateMailFunc(ctx, id, heroId, func(entity *mail_entity) (success bool) {
//		if entity.Keep != keep {
//			entity.Keep = keep
//			return true
//		}
//
//		return false
//	})
//
//	if err != nil {
//		return errors.Wrapf(err, "datastore.UpdateMailKeep() 失败")
//	}
//
//	return nil
//}
//
//func (d *DatastoreDbService) UpdateMailRead(ctx context.Context, id uint64, heroId int64, read bool) error {
//
//	err := d.updateMailFunc(ctx, id, heroId, func(entity *mail_entity) (success bool) {
//		if entity.Readed != read {
//			entity.Readed = read
//			return true
//		}
//
//		return false
//	})
//	if err != nil {
//		return errors.Wrapf(err, "datastore.UpdateMailRead() 失败")
//	}
//
//	return nil
//}
//
//func (d *DatastoreDbService) UpdateMailCollected(ctx context.Context, id uint64, heroId int64, collected bool) error {
//
//	err := d.updateMailFunc(ctx, id, heroId, func(entity *mail_entity) (success bool) {
//		if entity.Collected != collected {
//			entity.Collected = collected
//			return true
//		}
//
//		return false
//	})
//	if err != nil {
//		return errors.Wrapf(err, "datastore.UpdateMailCollected() 失败")
//	}
//
//	return nil
//}
//
//func (d *DatastoreDbService) ReadMultiMail(ctx context.Context, heroId int64, mailIds []uint64, hasReport bool) (*resdata.Prize, error) {
//
//	// 对id进行去重
//	idMap := make(map[uint64]struct{})
//
//	prizeList := make([]*shared_proto.PrizeProto, len(mailIds))
//
//	wg := sync.WaitGroup{}
//	for i, id := range mailIds {
//		if _, exist := idMap[id]; !exist {
//			idMap[id] = struct{}{}
//
//			wg.Add(1)
//			go func(id uint64, i int) {
//				defer wg.Done()
//
//				err := d.updateMailFunc(ctx, id, heroId, func(x *mail_entity) (success bool) {
//					if !hasReport {
//						if x.HasPrize && !x.Collected {
//							// 没有领取的邮件，自动领取
//							mailProto := &shared_proto.MailProto{}
//							if err := mailProto.Unmarshal(x.Data); err != nil {
//								logrus.Error("datastore.ReadMultiMail() MailProto.Unmarshal失败", err)
//								return false
//							}
//
//							x.Collected = true
//							prizeList[i] = mailProto.Prize
//
//							x.Readed = true
//							return true
//						}
//					}
//
//					if !x.Readed {
//						x.Readed = true
//						return true
//					}
//
//					return false
//				})
//
//				if err != nil {
//					prizeList[i] = nil
//				}
//
//				if err != nil {
//					logrus.Error("datastore.ReadMultiMail() Transaction失败", err)
//				}
//			}(id, i)
//		}
//	}
//	wg.Wait()
//
//	var prize *resdata.Prize
//	if !hasReport {
//		var prizeBuilder *resdata.PrizeBuilder
//		for _, v := range prizeList {
//			if v != nil {
//				if prizeBuilder == nil {
//					prizeBuilder = resdata.NewPrizeBuilder()
//				}
//				prizeBuilder.Add(resdata.UnmarshalPrize(v, d.datas))
//			}
//		}
//
//		if prizeBuilder != nil {
//			prize = prizeBuilder.Build()
//		}
//	}
//
//	return prize, nil
//}
//
////	// 对id进行去重
////	idMap := make(map[uint64]struct{})
////	newIds := make([]uint64, 0, len(mailIds))
////	for _, id := range mailIds {
////		if _, exist := idMap[id]; !exist {
////			idMap[id] = struct{}{}
////			newIds = append(newIds, id)
////		}
////	}
////
////	// 每个transaction处理25个数据
////	batchCount := 25
////	n := (len(newIds) + batchCount - 1 ) / batchCount
////
////	prizeList := make([]*shared_proto.PrizeProto, len(newIds))
////
////	wg := sync.WaitGroup{}
////	for i := 0; i < n; i++ {
////		startIndex := i * batchCount
////		endIndex := imath.Min((i+1)*batchCount, len(newIds))
////
////		var keys []*datastore.Key
////		for i := startIndex; i < endIndex; i++ {
////			keys = append(keys, datastore.IDKey(mail_kind, int64(newIds[i]), nil))
////		}
////
////		wg.Add(1)
////		go func(keys []*datastore.Key, batchIndex, startIndex int) {
////			defer wg.Done()
////			_, err := d.client.RunInTransaction(ctx, func(tx *nsds.Transaction) error {
////
////				array := make([]*mail_entity, len(keys))
////				if err := tx.GetMulti(keys, array); err != nil {
////					if err != datastore.ErrNoSuchEntity {
////						return err
////					}
////				}
////
////				var putKey []*datastore.Key
////				var putArray []*mail_entity
////				for i, x := range array {
////					if x == nil {
////						continue
////					}
////
////					if x.Receiver == heroId {
////						changed := false
////						if !hasReport {
////							if x.HasPrize && !x.Collected {
////								// 没有领取的邮件，自动领取
////								mailProto := &shared_proto.MailProto{}
////								if err := mailProto.Unmarshal(x.Data); err != nil {
////									logrus.Error("datastore.ReadMultiMail() MailProto.Unmarshal失败", err)
////									continue
////								}
////
////								x.Collected = true
////								prizeList[startIndex+i] = mailProto.Prize
////
////								changed = true
////							}
////						}
////
////						if !x.Readed {
////							x.Readed = true
////
////							changed = true
////						}
////
////						if changed {
////							putKey = append(putKey, keys[i])
////							putArray = append(putArray, x)
////						}
////					}
////				}
////
////				if len(putKey) > 0 {
////					if _, err := tx.PutMulti(putKey, putArray); err != nil {
////						for i := range array {
////							prizeList[startIndex+i] = nil
////						}
////						return err
////					}
////				}
////
////				return nil
////			})
////
////			if err != nil {
////				logrus.Error("datastore.ReadMultiMail() Transaction失败", err)
////			}
////		}(keys, i, startIndex)
////	}
////	wg.Wait()
////
////	var prize *resdata.Prize
////	if !hasReport {
////		var prizeBuilder *resdata.PrizeBuilder
////		for _, v := range prizeList {
////			if v != nil {
////				if prizeBuilder == nil {
////					prizeBuilder = resdata.NewPrizeBuilder()
////				}
////				prizeBuilder.Add(resdata.UnmarshalPrize(v, d.datas))
////			}
////		}
////
////		if prizeBuilder != nil {
////			prize = prizeBuilder.Build()
////		}
////	}
////
////	return prize, nil
////}
//
//func (d *DatastoreDbService) DeleteMultiMail(ctx context.Context, heroId int64, mailIds []uint64, hasReport bool) (err error) {
//	// 一个个来吧
//	idMap := make(map[uint64]struct{})
//
//	wg := sync.WaitGroup{}
//	for _, id := range mailIds {
//		if _, exist := idMap[id]; !exist {
//			idMap[id] = struct{}{}
//
//			wg.Add(1)
//			go func(id uint64) {
//				defer wg.Done()
//
//				err := d.deleteMailFunc(ctx, id, heroId, func(x *mail_entity) (success bool) {
//					// 收藏的邮件不删
//					if x.Keep {
//						return false
//					}
//
//					if !hasReport {
//						// 有奖励未领的不删（战报不检查奖励）
//						if x.HasPrize && !x.Collected {
//							return false
//						}
//					}
//
//					return true
//				})
//
//				if err != nil {
//					logrus.Error("datastore.ReadMultiMail() Transaction失败", err)
//				}
//			}(id)
//		}
//	}
//	wg.Wait()
//
//	return nil
//}
//
////// 对id进行去重
////idMap := make(map[uint64]struct{})
////newIds := make([]uint64, 0, len(mailIds))
////for _, id := range mailIds {
////	if _, exist := idMap[id]; !exist {
////		idMap[id] = struct{}{}
////		newIds = append(newIds, id)
////	}
////}
////
////// 每个transaction处理25个数据
////batchCount := 25
////n := (len(newIds) + batchCount - 1 ) / batchCount
////
////wg := sync.WaitGroup{}
////for i := 0; i < n; i++ {
////	startIndex := i * batchCount
////	endIndex := imath.Min((i+1)*batchCount, len(newIds))
////
////	var keys []*datastore.Key
////	for i := startIndex; i < endIndex; i++ {
////		keys = append(keys, datastore.IDKey(mail_kind, int64(newIds[i]), nil))
////	}
////
////	wg.Add(1)
////	go func(keys []*datastore.Key, batchIndex, startIndex int) {
////		defer wg.Done()
////		_, err := d.client.RunInTransaction(ctx, func(tx *nsds.Transaction) error {
////
////			array := make([]*mail_entity, len(keys))
////			if err := tx.GetMulti(keys, array); err != nil {
////				if err != datastore.ErrNoSuchEntity {
////					return err
////				}
////			}
////
////			var toRemoveKeys []*datastore.Key
////			for i, x := range array {
////				if x == nil {
////					continue
////				}
////
////				if x.Receiver != heroId {
////					// 不是自己的邮件
////					continue
////				}
////
////				// 收藏的邮件不删
////				if x.Keep {
////					continue
////				}
////
////				// 有奖励未领的不删
////				if x.HasPrize && !x.Collected {
////					continue
////				}
////
////				toRemoveKeys = append(toRemoveKeys, keys[i])
////			}
////
////			if len(toRemoveKeys) > 0 {
////				if err := tx.DeleteMulti(toRemoveKeys); err != nil {
////					if err != datastore.ErrNoSuchEntity {
////						return err
////					}
////				}
////			}
////
////			return nil
////		})
////
////		if err != nil {
////			logrus.Error("datastore.ReadMultiMail() Transaction失败", err)
////		}
////	}(keys, i, startIndex)
////}
////wg.Wait()
//
//// guild
//type guild_entity struct {
//	Data []byte `datastore:",noindex"`
//}
//
//func (d *DatastoreDbService) MaxGuildId(ctx context.Context) (int64, error) {
//
//	q := datastore.NewQuery(guild_kind).Order("-__key__").Limit(1).KeysOnly()
//
//	key, err := d.getFirstKey(ctx, q)
//	if err != nil {
//		if err == datastore.ErrNoSuchEntity {
//			return 0, nil
//		}
//
//		return 0, errors.Wrapf(err, "获取最大的帮派id失败")
//	}
//
//	return key.ID, nil
//}
//
//func (d *DatastoreDbService) CreateGuild(ctx context.Context, id int64, data []byte) error {
//
//	key := datastore.IDKey(guild_kind, id, nil)
//
//	_, err := d.client.RunInTransaction(ctx, func(tx *nsds.Transaction) error {
//
//		err := tx.Get(key, &guild_entity{})
//		if err != datastore.ErrNoSuchEntity {
//			if err != nil {
//				return err
//			}
//
//			return errors.Errorf("帮派已经存在, id: %d", id)
//		}
//
//		_, err = tx.Put(key, &guild_entity{
//			Data: data,
//		})
//
//		return err
//	})
//
//	if err != nil {
//		return errors.Wrapf(err, "datastore.CreateGuild() 失败")
//	}
//
//	return nil
//}
//
//func (d *DatastoreDbService) LoadAllGuild(ctx context.Context) ([]*sharedguilddata.Guild, error) {
//
//	array := make([]*sharedguilddata.Guild, 0)
//	for t := d.client.Run(ctx, datastore.NewQuery(guild_kind)); ; {
//		var x guild_entity
//		key, err := t.Next(&x)
//		if err == iterator.Done {
//			break
//		}
//
//		if err != nil {
//			return nil, errors.Wrapf(err, "datastore.LoadAllSharedHeroData() 失败")
//		}
//
//		g, err := d.parseGuild(key.ID, x.Data)
//		if err != nil {
//			return nil, errors.Wrapf(err, "datastore.LoadAllSharedHeroData() 失败")
//		}
//
//		array = append(array, g)
//	}
//
//	return array, nil
//}
//
//func (d *DatastoreDbService) parseGuild(id int64, guildData []byte) (*sharedguilddata.Guild, error) {
//
//	if len(guildData) <= 0 {
//		return nil, errors.Errorf("DbSerbive.parseGuild(%d), guildData len = 0", id)
//	}
//
//	var proto *server_proto.GuildServerProto
//
//	proto = &server_proto.GuildServerProto{}
//	err := proto.Unmarshal(guildData)
//	if err != nil {
//		return nil, errors.Wrapf(err, "DbSerbive.parseGuild(%d), Unmarshal GuildServerProto fail", id)
//	}
//
//	return sharedguilddata.UnmarshalGuild(id, proto, d.datas, d.timeService.CurrentTime())
//}
//
//func (d *DatastoreDbService) LoadGuild(ctx context.Context, id int64) (*sharedguilddata.Guild, error) {
//
//	key := datastore.IDKey(guild_kind, id, nil)
//
//	var x guild_entity
//	err := d.client.Get(ctx, key, &x)
//	if err != nil {
//		if err == datastore.ErrNoSuchEntity {
//			return nil, nil
//		}
//
//		return nil, errors.Wrapf(err, "datastore.LoadGuild() 失败")
//	}
//
//	return d.parseGuild(id, x.Data)
//}
//
//func (d *DatastoreDbService) SaveGuild(ctx context.Context, id int64, data []byte) error {
//
//	key := datastore.IDKey(guild_kind, id, nil)
//	_, err := d.client.Put(ctx, key, &guild_entity{
//		Data: data,
//	})
//	if err != nil {
//		return errors.Wrapf(err, "datastore.SaveGuild() 失败")
//	}
//
//	return nil
//}
//
//func (d *DatastoreDbService) DeleteGuild(ctx context.Context, id int64) error {
//
//	key := datastore.IDKey(guild_kind, id, nil)
//	err := d.client.Delete(ctx, key)
//	if err != nil {
//		return errors.Wrapf(err, "datastore.DeleteGuild() 失败")
//	}
//
//	d.deleteGuildLogs(ctx, id)
//
//	return nil
//}
//
//func (d *DatastoreDbService) deleteGuildLogs(ctx context.Context, guildId int64) error {
//
//	q := datastore.NewQuery(guild_logs_kind).Filter("Guild =", guildId)
//	if err := d.client.DeleteAll(ctx, q); err != nil {
//		return errors.Wrapf(err, "datastore.deleteGuildLogs() 失败")
//	}
//
//	return nil
//}
//
//type guild_log_entity struct {
//	Guild int64
//	Type  int
//	Data  []byte `datastore:",noindex"`
//}
//
//func (d *DatastoreDbService) LoadGuildLogs(ctx context.Context, guildId int64, logType shared_proto.GuildLogType, id int64, count uint64) ([]*shared_proto.GuildLogProto, error) {
//
//	count = u64.Min(u64.Max(count, d.datas.MiscConfig().MailMinBatchCount), d.datas.MiscConfig().MailMaxBatchCount)
//
//	q := datastore.NewQuery(guild_logs_kind)
//	if id > 0 {
//		q = q.Filter("__key__ <", d.client.NamespaceKey(datastore.IDKey(guild_logs_kind, int64(id), nil)))
//	}
//
//	q = q.Filter("Guild =", guildId)
//	q = q.Filter("Type =", int(logType))
//
//	q = q.Order("-__key__").Limit(int(count))
//
//	var xs []*guild_log_entity
//	_, err := d.client.GetAll(ctx, q, &xs)
//	if err != nil {
//		return nil, errors.Wrapf(err, "datastore.LoadGuildLogs() 失败")
//	}
//
//	var datas []*shared_proto.GuildLogProto
//	for _, x := range xs {
//		proto := &shared_proto.GuildLogProto{}
//		if err := proto.Unmarshal(x.Data); err != nil {
//			return nil, errors.Wrapf(err, "datastore.LoadGuildLogs() MailProto.Unmarshal失败")
//		}
//		//proto.Id = ?
//
//		datas = append(datas, proto)
//	}
//
//	return datas, nil
//
//}
//
//func (d *DatastoreDbService) maxGuildLogId(ctx context.Context) (uint64, error) {
//
//	q := datastore.NewQuery(guild_logs_kind).Order("-__key__").Limit(1).KeysOnly()
//
//	key, err := d.getFirstKey(ctx, q)
//	if err != nil {
//		if err == datastore.ErrNoSuchEntity {
//			return 0, nil
//		}
//
//		return 0, errors.Wrapf(err, "获取最大的联盟日志id失败")
//	}
//
//	return uint64(key.ID), nil
//}
//
//func (d *DatastoreDbService) InsertGuildLog(ctx context.Context, guildId int64, proto *shared_proto.GuildLogProto) error {
//
//	id := d.guildLogIdGen.Inc()
//	proto.Id = u64.Int32(id)
//
//	key := datastore.IDKey(guild_logs_kind, int64(id), nil)
//
//	_, err := d.client.RunInTransaction(ctx, func(tx *nsds.Transaction) error {
//
//		err := tx.Get(key, &guild_log_entity{})
//		if err != datastore.ErrNoSuchEntity {
//			if err != nil {
//				return err
//			}
//
//			return errors.Errorf("联盟日志id重复, id: %d", id)
//		}
//
//		_, err = tx.Put(key, &guild_log_entity{
//			Guild: guildId,
//			Type:  int(proto.Type),
//			Data:  util.SafeMarshal(proto),
//		})
//
//		return err
//	})
//
//	if err != nil {
//		return errors.Wrapf(err, "datastore.InsertGuildLog() 失败")
//	}
//
//	return nil
//}
//
//type baizhan_entity struct {
//	HeroId int64
//	Data   []byte `datastore:",noindex"`
//	Time   int64
//}
//
//func (d *DatastoreDbService) LoadBaiZhanRecord(ctx context.Context, heroId int64, count uint64) (isql.BytesArray, error) {
//
//	// 从2个列表中取，做一个merge sort
//	aq := datastore.NewQuery(baizhan_kind).Filter("HeroId =", heroId).Order("-Time").Limit(int(count))
//
//	var records [][]byte
//	for t := d.client.Run(ctx, aq); ; {
//		var x baizhan_entity
//		_, err := t.Next(&x)
//		if nsds.IgnoreFieldMissmatch(err) != nil {
//			if err == iterator.Done {
//				break
//			}
//
//			return nil, errors.Wrapf(err, "datastore.loadBaiZhanRecord() 失败")
//		}
//
//		records = append(records, x.Data)
//	}
//
//	return records, nil
//}
//
//func (d *DatastoreDbService) InsertBaiZhanReplay(ctx context.Context, attackerId, defenderId int64, replay *shared_proto.BaiZhanReplayProto, isDefenderNpc bool, time int64) error {
//
//	replayBytes := must.Marshal(replay)
//
//	keys := []*datastore.Key{datastore.IncompleteKey(baizhan_kind, nil)}
//	entity := []*baizhan_entity{
//		&baizhan_entity{
//			HeroId: attackerId,
//			Data:   replayBytes,
//			Time:   time,
//		},
//	}
//
//	if !isDefenderNpc {
//		keys = append(keys, datastore.IncompleteKey(baizhan_kind, nil))
//		entity = append(entity, &baizhan_entity{
//			HeroId: defenderId,
//			Data:   replayBytes,
//			Time:   time,
//		})
//	}
//
//	_, err := d.client.PutMulti(ctx, keys, entity)
//	if err != nil {
//		return errors.Wrapf(err, "datastore.InsertBaiZhanReplay() 失败")
//	}
//
//	return nil
//}
//
//type farm_entity struct {
//	FarmData []byte `datastore:",noindex"`
//}
//
//var ErrUnsupport = errors.Errorf("datastore unsupport this func")
//
//func (d *DatastoreDbService) AddChatMsg(ctx context.Context, senderId int64, room []byte, proto *shared_proto.ChatMsgProto) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) RemoveChatMsg(ctx context.Context, senderId int64) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) ListHeroChatMsg(ctx context.Context, room []byte, minChatId uint64) ([]*shared_proto.ChatMsgProto, error) {
//	return nil, ErrUnsupport
//}
//
////func (d *DatastoreDbService) updateChatWindow(ctx context.Context, heroId int64, room []byte, chatMsg *shared_proto.ChatMsgProto) error {
////	return errUnsupport
////}
//
//func (d *DatastoreDbService) DeleteChatWindow(ctx context.Context, heroId int64, room []byte) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) ReadChat(ctx context.Context, heroId int64, room []byte) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) ListHeroChatWindow(ctx context.Context, heroId int64) ([]uint64, isql.BytesArray, error) {
//	return nil, nil, ErrUnsupport
//}
//
//func (d *DatastoreDbService) UpdateChatWindow(ctx context.Context, heroId int64, room []byte, target []byte, addUnread bool, sendTime int32, updateSendTime bool) error {
//	return ErrUnsupport
//}
//
///** 农场 **/
//
//type farm_cube struct {
//	HeroId          int64
//	Cube            uint64
//	StartTimeInt    int32
//	RemoveTimeInt   int32
//	ConflictTimeInt int32
//	ResId           int64
//	StealTimes      int64
//	ConflictHeroId  int64
//}
//
//func (d *DatastoreDbService) CreateFarmCube(ctx context.Context, cube *entity.FarmCube) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) SaveFarmCube(ctx context.Context, cube *entity.FarmCube) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) PlantFarmCube(ctx context.Context, heroId int64, cube cb.Cube, startTime int64, ripeTime int64, resId uint64) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) UpdateFarmCubeState(ctx context.Context, heroId int64, cube cb.Cube, conflictedTime, removeTime int64) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) LoadFarmCube(ctx context.Context, heroId int64, cb cb.Cube) (*entity.FarmCube, error) {
//	return nil, ErrUnsupport
//}
//
//func (d *DatastoreDbService) LoadFarmCubes(ctx context.Context, heroId int64) ([]*entity.FarmCube, error) {
//	return nil, ErrUnsupport
//}
//
//func (d *DatastoreDbService) LoadFarmHarvestCubes(ctx context.Context, heroId int64) ([]*entity.FarmCube, error) {
//	return nil, ErrUnsupport
//}
//
//func (d *DatastoreDbService) LoadFarmStealCubes(ctx context.Context, heroId int64, ripeTime int64, maxStealTimes uint64) ([]*entity.FarmCube, error) {
//
//	return nil, ErrUnsupport
//}
//
//func (d *DatastoreDbService) UpdateFarmStealTimes(ctx context.Context, heroId int64, cubes []cb.Cube) error {
//
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) RemoveFarmCube(ctx context.Context, heroId int64, cube cb.Cube) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) RemoveFarmLog(ctx context.Context, logTime int32) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) LoadFarmLog(ctx context.Context, heroId int64, size uint64) ([]*shared_proto.FarmStealLogProto, error) {
//	return nil, ErrUnsupport
//}
//
//func (d *DatastoreDbService) CreateFarmLog(ctx context.Context, logProto *shared_proto.FarmStealLogProto) error {
//	return nil
//}
//
//func (d *DatastoreDbService) AddFarmSteal(ctx context.Context, heroId, thiefId int64, cube cb.Cube) error {
//	return nil
//}
//
//func (d *DatastoreDbService) LoadFarmStealCount(ctx context.Context, heroId, thiefId int64, cube cb.Cube) (count uint64, err error) {
//	return 0, ErrUnsupport
//}
//
//func (d *DatastoreDbService) LoadCanStealCount(ctx context.Context, heroId, thiefId int64, ripeTime int64, maxStealTime uint64) (count uint64, err error) {
//	return 0, ErrUnsupport
//}
//
//func (d *DatastoreDbService) LoadCanStealCube(ctx context.Context, heroId, thiefId, ripeTime int64, maxStealTimes uint64) (array []*entity.FarmCube, err error) {
//	return nil, ErrUnsupport
//}
//
//func (d *DatastoreDbService) RemoveFarmSteal(ctx context.Context, heroId int64, cubes []cb.Cube) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) ResetFarmCubes(ctx context.Context, heroId int64) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) ResetConflictFarmCubes(ctx context.Context, heroId int64) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) LoadRecommendHeros(ctx context.Context, needLocation bool, location, minLevel, page, size uint64, excludeHeroId int64) (heros []*entity.Hero, err error) {
//	return []*entity.Hero{}, ErrUnsupport
//}
//
//func (d *DatastoreDbService) LoadHerosByName(ctx context.Context, text string, page, size uint64) (heros []*entity.Hero, err error) {
//	return []*entity.Hero{}, ErrUnsupport
//}
//
//func (d *DatastoreDbService) GMFarmRipe(ctx context.Context, heroId int64, startTime, ripeTime int64) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) AddMcWarRecord(ctx context.Context, mcWarId, mcId uint64, record *shared_proto.McWarFightRecordProto) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) AddMcWarHeroRecord(ctx context.Context, mcWarId, mcId uint64, heroId int64, record *shared_proto.McWarTroopAllRecordProto) error {
//	return ErrUnsupport
//}
//
//func (d *DatastoreDbService) LoadMcWarRecord(ctx context.Context, mcWarId, mcId uint64) (record *shared_proto.McWarFightRecordProto, err error) {
//	return nil, ErrUnsupport
//}
//
//func (d *DatastoreDbService) LoadMcWarHeroRecord(ctx context.Context, mcWarId, mcId uint64, heroId int64) (record *shared_proto.McWarTroopAllRecordProto, err error) {
//	return nil, ErrUnsupport
//}
//
////func (d *DatastoreDbService) AddTlog(ctx context.Context, msg []byte, createTime time.Time) error {
////	return ErrUnsupport
////}
////
////func (d *DatastoreDbService) RemoveTlog(ctx context.Context, ids []uint64) error {
////	return ErrUnsupport
////}
////
////func (d *DatastoreDbService) LoadTlog(ctx context.Context, size uint64) ([]*entity.TlogDBInfo, error) {
////	return nil, ErrUnsupport
////}
