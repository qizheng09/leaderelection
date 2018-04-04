package main

import (
	"github.com/spf13/viper"
	"flag"
	log "github.com/inconshreveable/log15"
	"github.com/qizheng09/leaderelection/pkg/etcdcli"
)

func main() {

	// load config
	configpath := flag.String("configpath", "", "use -configpath=<configpath>")
	viper.SetConfigName("config")
	viper.AddConfigPath(configpath)
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

}
