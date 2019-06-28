package main

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"gomicro-repo/server/config"
	"gomicro-repo/server/constants"
	_ "gomicro-repo/server/db"
	"gomicro-repo/server/impl"
	"gomicro-repo/server/logger"
	"gomicro-repo/server/zkmgr"
)

func main() {
	logger.Info(fmt.Sprintf("gomicro-repo server start ..."))

	zkConn := registerServer()
	defer func() {
		if zkConn != nil{
			zkConn.Close()
		}
	}()

	impl.StartGrpcServer()
}

func registerServer() (conn *zk.Conn){
	var (
		izk *zkmgr.IZk
		err error
	)

	if izk,err = zkmgr.NewIZk([]string{config.GetConf().ZK.Hostsports});err!=nil{
		logger.Error(fmt.Sprintf("registerServer, fail to NewIZk. err:%v.",err))
		return
	}

	conn = izk.Conn
	logger.Info(fmt.Sprintf("registerServer, get zk conn:%p.",conn))

	//注册zk节点
	err = izk.RegisterPerServer(constants.ServerName)
	if err != nil {
		logger.Warn(fmt.Sprintf("registerServer, fail to RegisterPerServer node error: %v.", err))
	}

	err = izk.RegisterEphServer(constants.ServerName,constants.RootPath+config.GetConf().Server.Address+config.GetConf().Server.GrpcPort)
	if err != nil {
		logger.Error(fmt.Sprintf("registerServer, fail to RegisterEphServer node error: %v.", err))
		return
	}

	logger.Info(fmt.Sprintf("registerServer, succeed to RegistServer node."))
	return
}
