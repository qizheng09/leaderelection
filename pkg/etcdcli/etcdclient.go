package etcdcli

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
)

const DialTimeout  = 5 * time.Second

// Client defines typed wrappers for the etcd API.
type Client struct {
	Client *clientv3.Client
}

// NewEtcdClient create a etcd client that is configured to be used with the given endpoints
func NewEtcdClient(endPoints []string) *Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endPoints,
		DialTimeout: DialTimeout,
	})
	if err != nil {
		fmt.Errorf("New etcd client error %+v", err)
	}
	return &Client{
		Client: cli,
	}
}

