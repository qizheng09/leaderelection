package etcdcli

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"golang.org/x/net/context"
	log "github.com/inconshreveable/log15"
	"sync"
)

const (
	dialTimeout = 5 * time.Second
	requestTimeout = 3 * time.Second
)

// Client defines typed wrappers for the etcd API.
type Client struct {
	Client *clientv3.Client
}

// NewEtcdClient create a etcd client that is configured to be used with the given endpoints
func NewEtcdClient(endPoints []string) *Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endPoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		fmt.Errorf("New etcd client error %+v", err)
	}
	return &Client{
		Client: cli,
	}
}

// ConfirmKey used to confirm is the key already exists
func (c *Client) ConfirmKey(key string, value string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := c.Client.Get(ctx, key)
	cancel()
	if err != nil {
		log.Error("ConfirmKey error", "error", err.Error())
	}

	if resp.Kvs == nil {
		return false
	}
	return true
}

// CreateKeyWithKeepaliveTTL create key with ttl and keep alive
func (c *Client) CreateKeyWithKeepaliveTTL(key string, ttl int64, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	respLease, err := c.Client.Grant(ctx, ttl)
	cancel()
	if err != nil {
		log.Error("Create release error", "error", err.Error())
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = c.Client.Put(ctx, key, value, clientv3.WithLease(respLease.ID))
	cancel()
	if err != nil {
		log.Error("CreateKeyWithTTL error", "error", err.Error())
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = c.Client.KeepAlive(ctx, respLease.ID)
	cancel()
	return err
}

func consumerService(c chan *clientv3.WatchResponse, services *sync.Map, serviceChange func(services *sync.Map, e*clientv3.Event)) {
	select {
	case curResp, ok := <-c:
		if ok {
			for _, e := range curResp.Events {
				serviceChange(services, e)
			}
		}
	}

}

// AddServicesWatcher add service Watcher
func (c *Client) AddServicesWatcher(key string, services *sync.Map, serviceChange func(services *sync.Map, e*clientv3.Event)) {
	for {
		consumerChan := make(chan *clientv3.WatchResponse, 5000)
		go consumerService(consumerChan, services, serviceChange)
		respReceiverChan := c.Client.Watch(context.Background(), key, clientv3.WithPrefix())
		for wresp := range respReceiverChan {
			if wresp.Canceled {
				log.Info("Info: watch cancled!")
			}
			log.Info("Info: receive change!")
			consumerChan <- &wresp
		}
	}
}

// GetWithPrefix get from etcd with prefix
func (c *Client) GetWithPrefix(key string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := c.Client.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	return resp, err
}
