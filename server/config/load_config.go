package config

import (
	"fmt"
	"github.com/jinzhu/configor"
	"gomicro-repo/server/logger"
	mflag "gomicro-repo/server/my-flag"
	"os"
)

var conf AppConfig

type AppConfig struct {
	MongoDb struct {
		Hostsports string
		Dbname     string
		Userpass   string
		ReplicaSet string
		Role       string
	}

	ZK struct{
		Hostsports string
	}

	Server struct{
		Address string
		GrpcPort string
	}
}

func LoadConfig(file string) {
	conf = AppConfig{}
	err := configor.Load(&conf, file)
	if err != nil {
		logger.Error("Failed to find configuration ", file)
		os.Exit(1)
	}
}

func init() {
	confPath := mflag.ConfPath
	LoadConfig(confPath)
	for k, v := range os.Args {
		logger.Info(fmt.Sprintf("k:%d, v:%s", k, v))
	}

	logger.Info("\r\n")
	logger.Info("loading config...... \r\n")
	logger.Info("config path:", confPath)
	logger.Info("mongo addr:", conf.MongoDb.Hostsports)
	logger.Info("zookeeper addr:", conf.ZK.Hostsports)
	logger.Info("server addr:", conf.Server.Address)
	logger.Info("server grpcport:", conf.Server.GrpcPort)
}

func GetConf() AppConfig {
	return conf
}
