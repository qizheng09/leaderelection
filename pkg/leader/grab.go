package leader

import (
	"github.com/qizheng09/leaderelection/pkg/etcdcli"
	"time"
	log "github.com/inconshreveable/log15"
)

const (
	// KEY1 is the main key node to grab
	KEY1 = "chain/leader1"
	// TTL for KEY1
	TTL = 600
	T  = time.Second * 3
)



// GrabKey is a demo for grab distribute key
func GrabKey(serviceinfo string, client *etcdcli.Client) {
	ticker := time.NewTicker(T)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			confirmResult := client.ConfirmKey(KEY1, serviceinfo)
			if !confirmResult {
				err := client.CreateKeyWithKeepaliveTTL(KEY1, TTL, serviceinfo)
				if err != nil {
					log.Info("Grab leader error!")
					continue
				}

			} else {
				log.Info("Key chain/leader1 is already exists!")
				continue
			}
		}
	}
}