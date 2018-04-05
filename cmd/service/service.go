package main

import (
	"flag"
	"github.com/spf13/viper"
	log "github.com/inconshreveable/log15"
	"github.com/qizheng09/leaderelection/pkg/etcdcli"
	"sync"
	"github.com/qizheng09/leaderelection/pkg/service"
)


func main() {
	// load config
	configpath := flag.String("configpath", "", "use -configpath=<configpath>")
	viper.SetConfigName("config")
	viper.AddConfigPath(*configpath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Load config file error!", "err", err.Error())
		return
	}

	// init etcd client
	endpoints := viper.GetStringSlice("endpoints")
	cli := etcdcli.NewEtcdClient(endpoints)
	if cli == nil {
		log.Error("New etcd client error!")
		return
	}
	defer cli.Client.Close()

	// TODO: serviceinfo is which node can be uniquely identified
	serviceinfo := viper.GetString("serviceinfo")
	// All services of the network
	var services *sync.Map
	err = service.InitServices(services, cli)
	if err != nil {
		return
	}
	go service.SetServicesWatcher(services, cli)

	err = service.RegisterServices(serviceinfo, cli)
	if err != nil {
		log.Error("RegisterServices error!", "err", err.Error())
	}
}
