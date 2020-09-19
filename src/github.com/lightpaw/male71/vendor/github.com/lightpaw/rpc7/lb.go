package rpc7

import (
	"google.golang.org/grpc"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	"time"
	"github.com/lightpaw/logrus"
	"fmt"
	"strings"
)

func NewEtcdClient(cli *clientv3.Client, target string) (*Client, error) {

	r := &etcdnaming.GRPCResolver{Client: cli}
	b := grpc.RoundRobin(r)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, removeLastSlash(target), grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "grpclb fail: %v", target)
	}

	return &Client{
		conn: conn,
		rsc:  NewRpcServiceClient(conn),
	}, nil
}

func removeLastSlash(target string) string {
	if len(target) > 0 && strings.HasSuffix(target, "/") {
		return target[0:len(target)-1]
	}
	return target
}

func RegisterAddr(cli *clientv3.Client, target, host string, port, ttl uint32, stopNotifier chan struct{}) (loopNotifier chan struct{}) {

	addr := fmt.Sprintf("%s:%d", host, port)

	key := removeLastSlash(target) + "/" + addr
	value := fmt.Sprintf(`{"Addr":"%s"}`, addr)

	interval := time.Duration(ttl) * time.Second
	if interval < 5*time.Second {
		interval = 5 * time.Second
	}

	// 先做一次更新，
	if err := updateAddr(cli, key, value, int64(ttl)); err != nil {
		logrus.WithError(err).Error("Refresh Rpc Addr fail")
	}

	loopNotifier = make(chan struct{})
	go func() {
		defer close(loopNotifier)
		ticker := time.NewTicker((interval - time.Second) / 2)

		for {
			select {
			case <-ticker.C:
				if err := updateAddr(cli, key, value, int64(ttl)); err != nil {
					logrus.WithError(err).Error("Refresh Rpc Addr fail")
				}
			case <-stopNotifier:
				// 执行一下删除操作
				deleteKey(cli, key)
				return
			}
		}
	}()

	return loopNotifier
}

func updateAddr(cli *clientv3.Client, key, value string, ttl int64) error {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := cli.Grant(ctx, ttl)
	if err != nil {
		return errors.Wrapf(err, "clientv3.Grant(ttl) fail")
	}

	if _, err := cli.Put(ctx, key, value, clientv3.WithLease(resp.ID)); err != nil {
		return errors.Wrapf(err, "clientv3.Put() fail")
	}

	return nil
}

func deleteKey(cli *clientv3.Client, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := cli.Delete(ctx, key); err != nil {
		logrus.WithError(err).Error("clientv3.Delete(key) fail")
	}
}
