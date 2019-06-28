package zkmgr

import (
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"gomicro-repo/server/config"
	"gomicro-repo/server/constants"
	"gomicro-repo/server/logger"
	"math/rand"
	"time"
)

var serverList []string

type IZk struct {
	host []string
	Conn *zk.Conn
}

func NewIZk(host []string) (*IZk,error) {
	var (
		izk *IZk
		err error
	)

	izk = &IZk{
		host: host,
	}

	if izk.Conn,err = izk.GetConnect();err!=nil{
		errInfo := fmt.Sprintf("NewIZk, fail to GetConnect, connect zk error:%v.",err)
		logger.Error(errInfo)
		return nil,errors.New(errInfo)
	}

	return izk,err
}

func (z IZk)GetConnect() (conn *zk.Conn, err error) {
	hosts := z.host
	conn, _, err = zk.Connect(hosts, 5*time.Second)
	return
}

func (z IZk)RegisterPerServer(name string) (err error) {
	// "/server"
	var (
		str string
	)

	str,err = z.Conn.Create(name, []byte("hello Go!"), 0, zk.WorldACL(zk.PermAll))
	if err!=nil{
		logger.Warn(fmt.Sprintf("RegisterPerServer, fail to Create node:%s. err:%v.",name,err))
	}else {
		logger.Info(fmt.Sprintf("RegisterPerServer, succeed to Create node:%s. retStr:%s.",name,str))
	}

	return
}

func (z IZk)RegisterEphServer(name,host string) (err error) {
	// "/server/192.168.1.100:8080"
	var (
		str string
	)

	str, err = z.Conn.Create(name+host, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err!=nil{
		logger.Error(fmt.Sprintf("RegisterEphServer, fail to Create node:%s. err:%v.",name+host,err))
	}else {
		logger.Info(fmt.Sprintf("RegisterEphServer, succeed to Create node:%s. retStr:%s.",name+host,str))
	}

	return
}

func (z IZk)GetServerList(name string) (list []string, err error) {
	// "/server"
	list, _, err = z.Conn.Children(name)
	logger.Info(fmt.Sprintf("IZk, GetServerList list:%v", list))
	return
}

func (z IZk)WatchServerList(path string) (chan []string, chan error) {
	snapshots := make(chan []string)
	errors := make(chan error)

	go func() {
		for {
			snapshot, _, events, err := z.Conn.ChildrenW(path)
			if err != nil {
				errors <- err
				return
			}
			snapshots <- snapshot
			evt := <-events
			if evt.Err != nil {
				errors <- evt.Err
				return
			}
		}
	}()

	return snapshots, errors
}

func GetServerHost() (host string, err error) {
	count := len(serverList)
	if count == 0 {
		errInfo := fmt.Sprintf("GetServerHost, server list is empty!")
		logger.Error(errInfo)
		err = errors.New(errInfo)
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	host = serverList[r.Intn(count)]
	return
}

func StartZkService() {
	var (
		izk *IZk
		err error
	)

	if izk,err = NewIZk([]string{config.GetConf().ZK.Hostsports});err!=nil{
		logger.Error(fmt.Sprintf("StartZkService, fail to NewIZk. err:%v.",err))
		return
	}

	logger.Info(fmt.Sprintf("StartZkService ppp111, get zk conn:%p.",izk.Conn))
	serverList, err = izk.GetServerList(constants.ServerName)
	if err != nil {
		logger.Error(fmt.Sprintf("StartZkService, get server list error: %s \n", err))
		return
	}

	logger.Info(fmt.Sprintf("StartZkService, first GetServerList:%v.",serverList))
	snapshots, errors := izk.WatchServerList(constants.ServerName)
	go func() {
		for {
			select {
			case serverList = <-snapshots:
				logger.Info(fmt.Sprintf("StartZkService 1111:%+v\n", serverList))
			case err := <-errors:
				logger.Info(fmt.Sprintf("StartZkService 2222:%+v\n", err))
			}
		}
	}()
}
