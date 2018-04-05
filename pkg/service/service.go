package service

import (
	"github.com/qizheng09/leaderelection/pkg/etcdcli"
	"sync"
	"path"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"

)

const (
	// KEY is the key node to add
	KEY = "chain/services"
	// TTL for KEY
	TTL = 600
)

// InitServices start to get all services from etcd
func InitServices(services *sync.Map, client *etcdcli.Client) error {
	resp, err := client.GetWithPrefix(KEY)
	if err != nil {
		return err
	}
	for k, v := range resp.Kvs {
		if _, ok := services.Load(k); !ok {
			services.Store(k, v)
		}
	}
	return nil
}

func serviceChange(services *sync.Map, e *clientv3.Event) {
	switch e.Type {
	case mvccpb.PUT:
		if _, ok := services.Load(e.Kv.Key); !ok {
			services.Store(e.Kv.Key, e.Kv.Value)
		}
	case mvccpb.DELETE:
		if _, ok := services.Load(e.Kv.Key); ok {
			services.Delete(e.Kv.Key)
		}
	}
}

// SetServicesWatcher set to watch services
func SetServicesWatcher(services *sync.Map, client *etcdcli.Client) {
	go client.AddServicesWatcher(KEY, services, serviceChange)
}

// RegisterServices register to network
func RegisterServices(serviceinfo string, client *etcdcli.Client) error {
	key := path.Join(KEY, serviceinfo)
	return client.CreateKeyWithKeepaliveTTL(key, TTL, serviceinfo)
}