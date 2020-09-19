package nsds
//
//import (
//	"cloud.google.com/go/datastore"
//	"context"
//	"github.com/lightpaw/logrus"
//	"github.com/pkg/errors"
//	"os"
//	"time"
//)
//
//func NewNamespaceClient(emulatorHost, projectID, namespace string) (*NamespaceClient, error) {
//
//	logrus.Infof("EmulatorHost: %s, ProjectID: %s, Namespace: %s", emulatorHost, projectID, namespace)
//
//	if host := os.Getenv("DATASTORE_EMULATOR_HOST"); len(host) > 0 {
//		logrus.Info("使用环境变量中的datastore模拟器地址，", host)
//	} else {
//		if len(emulatorHost) > 0 {
//			logrus.Info("使用配置中的变量中的datastore模拟器地址，", emulatorHost)
//			os.Setenv("DATASTORE_EMULATOR_HOST", emulatorHost)
//		} else {
//			logrus.Info("使用正式环境的datastore")
//		}
//	}
//
//	timeoutContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	client, err := datastore.NewClient(timeoutContext, projectID)
//	if err != nil {
//		return nil, errors.Wrapf(err, "datastore.NewClient 错误")
//	}
//
//	// 测试连接可用性
//
//	t, err := client.NewTransaction(timeoutContext)
//	if err != nil {
//		return nil, errors.Wrapf(err, "datastore connect client 错误, host: %s", emulatorHost)
//	}
//	t.Rollback()
//
//	logrus.Info("创建client成功")
//
//	return &NamespaceClient{
//		client:    client,
//		namespace: namespace,
//	}, nil
//}
//
//type NamespaceClient struct {
//	client *datastore.Client
//
//	namespace string
//}
//
//func (c *NamespaceClient) newTransaction(tx *datastore.Transaction) *Transaction {
//	return &Transaction{
//		c:  c,
//		tx: tx,
//	}
//}
//
//func (c *NamespaceClient) NamespaceKey(key *datastore.Key) *datastore.Key {
//	key.Namespace = c.namespace
//	return key
//}
//
//func (c *NamespaceClient) NamespaceKeys(keys []*datastore.Key) []*datastore.Key {
//	for _, key := range keys {
//		key.Namespace = c.namespace
//	}
//	return keys
//}
//
//func (c *NamespaceClient) NamespaceQuery(q *datastore.Query) *datastore.Query {
//	return q.Namespace(c.namespace)
//}
//
//func (c *NamespaceClient) AllocateIDs(ctx context.Context, keys []*datastore.Key) ([]*datastore.Key, error) {
//	return c.client.AllocateIDs(ctx, c.NamespaceKeys(keys))
//}
//
//func (c *NamespaceClient) Close() error {
//	return c.client.Close()
//}
//
//func (c *NamespaceClient) Count(ctx context.Context, q *datastore.Query) (int, error) {
//	return c.client.Count(ctx, c.NamespaceQuery(q))
//}
//
//func (c *NamespaceClient) Delete(ctx context.Context, key *datastore.Key) error {
//	return c.client.Delete(ctx, c.NamespaceKey(key))
//}
//
//func (c *NamespaceClient) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
//	return c.client.DeleteMulti(ctx, c.NamespaceKeys(keys))
//}
//
//func (c *NamespaceClient) DeleteAll(ctx context.Context, q *datastore.Query) error {
//	keys, err := c.GetAll(ctx, q.KeysOnly(), nil)
//	if err != nil {
//		return err
//	}
//
//	return c.DeleteMulti(ctx, keys)
//}
//
//func (c *NamespaceClient) Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
//	return IgnoreFieldMissmatch(c.client.Get(ctx, c.NamespaceKey(key), dst))
//}
//
//func (c *NamespaceClient) GetAll(ctx context.Context, q *datastore.Query, dst interface{}) ([]*datastore.Key, error) {
//	keys, err := c.client.GetAll(ctx, c.NamespaceQuery(q), dst)
//	return keys, IgnoreFieldMissmatch(err)
//}
//
//func (c *NamespaceClient) GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
//	return IgnoreFieldMissmatch(c.client.GetMulti(ctx, c.NamespaceKeys(keys), dst))
//}
//
//func (c *NamespaceClient) NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (*Transaction, error) {
//	tx, err := c.client.NewTransaction(ctx, opts...)
//	if err != nil {
//		return nil, err
//	}
//
//	return c.newTransaction(tx), nil
//}
//
//func isRetryError(err error) bool {
//	return false // 判断出是否需要重试
//}
//
//func (c *NamespaceClient) Put(ctx context.Context, key *datastore.Key, src interface{}) (returnKey *datastore.Key, err error) {
//
//	DoBackoffRetry(func() (isRetry bool) {
//		returnKey, err = c.client.Put(ctx, c.NamespaceKey(key), src)
//		return isRetryError(err)
//	})
//
//	return
//}
//func (c *NamespaceClient) PutMulti(ctx context.Context, keys []*datastore.Key, src interface{}) (returnKeys []*datastore.Key, err error) {
//	DoBackoffRetry(func() (isRetry bool) {
//		returnKeys, err = c.client.PutMulti(ctx, c.NamespaceKeys(keys), src)
//		return isRetryError(err)
//	})
//	return
//}
//
//func (c *NamespaceClient) Run(ctx context.Context, q *datastore.Query) *datastore.Iterator {
//	return c.client.Run(ctx, c.NamespaceQuery(q))
//}
//
//func (c *NamespaceClient) RunInTransaction(ctx context.Context, f func(tx *Transaction) error, opts ...datastore.TransactionOption) (*datastore.Commit, error) {
//	return c.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
//		return f(c.newTransaction(tx))
//	}, opts...)
//}
//
//func IgnoreFieldMissmatch(err error) error {
//	if err != nil {
//		switch err.(type) {
//		case *datastore.ErrFieldMismatch:
//			// 忽略不存在的字段
//			logrus.WithError(err).Warnf("存在不匹配的字段")
//			return nil
//		default:
//			return err
//		}
//	}
//	return nil
//}
//
//type Transaction struct {
//	c  *NamespaceClient
//	tx *datastore.Transaction
//}
//
//func (c *Transaction) Tx() *datastore.Transaction {
//	return c.tx
//}
//
//func (c *Transaction) Delete(key *datastore.Key) error {
//	return c.tx.Delete(c.c.NamespaceKey(key))
//}
//
//func (c *Transaction) DeleteMulti(keys []*datastore.Key) error {
//	return c.tx.DeleteMulti(c.c.NamespaceKeys(keys))
//}
//
//func (c *Transaction) Get(key *datastore.Key, dst interface{}) error {
//	return IgnoreFieldMissmatch(c.tx.Get(c.c.NamespaceKey(key), dst))
//}
//
//func (c *Transaction) GetMulti(keys []*datastore.Key, dst interface{}) error {
//	return IgnoreFieldMissmatch(c.tx.GetMulti(c.c.NamespaceKeys(keys), dst))
//}
//
//func (c *Transaction) Put(key *datastore.Key, src interface{}) (*datastore.PendingKey, error) {
//	return c.tx.Put(c.c.NamespaceKey(key), src)
//}
//func (c *Transaction) PutMulti(keys []*datastore.Key, src interface{}) ([]*datastore.PendingKey, error) {
//	return c.tx.PutMulti(c.c.NamespaceKeys(keys), src)
//}
