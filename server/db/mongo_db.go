package db

import (
	"fmt"
	"gomicro-repo/server/config"
	"gomicro-repo/server/logger"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
)

var mgoS *mgo.Session

func InitMgo() *mgo.Session {
	logger.Info("start to connect to mongoDB...")

	var(
		err error
		connstr = "mongodb://"
		role = config.GetConf().MongoDb.Role
		userpass = config.GetConf().MongoDb.Userpass
		hostsports = config.GetConf().MongoDb.Hostsports
		replicaset = config.GetConf().MongoDb.ReplicaSet
	)

	if !strings.EqualFold("", userpass) {
		connstr += userpass + "@"
	}

	connstr += fmt.Sprintf("%s/%s", hostsports, role)
	if !strings.EqualFold("", replicaset) {
		connstr += "?replicaSet=" + replicaset
	}

	//连接
	mgoS, err = mgo.DialWithTimeout(connstr, time.Second * 3)
	for {
		if err != nil {
			logger.Info(err.Error(), "retry to conn... ")
			mgoS, err = mgo.DialWithTimeout(connstr, time.Second * 1)
		} else {
			logger.Info("conn success. prepare to ping...")
			err = mgoS.Ping()
			if err != nil {
				logger.Error(err.Error(), "ping Mongo error")
			} else {
				logger.Info("ping Mongo success.")
				break
			}
		}
	}

	//设置模式
	mgoS.SetMode(mgo.Monotonic, true)
	logger.Info("succeed connect to mongodb...")

	return mgoS
}

func init() {
	//初始化mongodb
	mgoS = InitMgo()
}

func GetMs() *mgo.Session {
	if mgoS == nil{
		mgoS = InitMgo()
	}

	return mgoS
}