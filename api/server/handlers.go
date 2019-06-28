package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"gomicro-repo/api/logger"
	apiproto "gomicro-repo/api/proto"
	"gomicro-repo/api/zkmgr"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
)

func GetUserInfo(ctx *gin.Context) {
	var (
		tokenOk bool
		err     error
	)

	token := ctx.Request.Header.Get("Authorization")
	logger.Info(fmt.Sprintf("GetUserInfo, Authorization:%v.", token))

	//call user server to auth token
	if tokenOk, err = verifyToken(token); err != nil {
		//token鉴权失败，调用链出错
		logger.Error(fmt.Sprintf("GetUserInfo, fail to verifyToken. err:%v.", err))
		ctx.ProtoBuf(http.StatusBadRequest, &apiproto.GetUserInfoRsp{
			Code: 10100,
			Msg:  "token鉴权失败，调用链出错",
			Data: nil,
		})
		return
	}

	logger.Info(fmt.Sprintf("GetUserInfo, tokenOk:%v.", tokenOk))
	if !tokenOk {
		//token鉴权失败，token无效
		logger.Error(fmt.Sprintf("GetUserInfo, token is invalid."))
		ctx.ProtoBuf(http.StatusOK, &apiproto.GetUserInfoRsp{
			Code: 10200,
			Msg:  "token鉴权失败，token无效",
			Data: nil,
		})
		return
	}

	//token鉴权成功，可以拉取数据
	bodyPro, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("GetUserInfo, fail to ReadAll Request.Body, err:%v.", err))
		ctx.ProtoBuf(http.StatusOK, &apiproto.GetUserInfoRsp{
			Code: 10300,
			Msg:  "请求参数读取失败",
			Data: nil,
		})
		return
	}

	getUserInfoReq := new(apiproto.GetUserInfoReq)
	err = proto.Unmarshal(bodyPro, getUserInfoReq)
	if err != nil {
		errInfo := fmt.Sprintf("GetUserInfo, fail to Unmarshal, bodyPro type error:%v", err)
		logger.Error(errInfo)
		ctx.ProtoBuf(http.StatusOK, &apiproto.GetUserInfoRsp{
			Code: 10400,
			Msg:  "请求参数格式错误，反序列化失败",
			Data: nil,
		})
		return
	}

	logger.Info(fmt.Sprintf("GetUserInfo req params ok, userId:%v.", getUserInfoReq.UserId))

	//call gRPC server, get userInfo from db
	getUserInfoRsp, err := callGetUserInfoService(getUserInfoReq)
	if err != nil {
		logger.Error(fmt.Sprintf("GetUserInfo, fail to callGetUserInfoService, err:%v.", err))
		ctx.ProtoBuf(http.StatusOK, &apiproto.GetUserInfoRsp{
			Code: 10500,
			Msg:  "内部gRPC调用失败",
			Data: nil,
		})
		return
	}
	logger.Info(fmt.Sprintf("GetUserInfo, succeed to callGetUserInfoService. rsp: %v.", getUserInfoRsp))

	ctx.ProtoBuf(http.StatusOK, getUserInfoRsp)
	return
}

func GetAllUserInfo(ctx *gin.Context) {
	bodyPro, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("GetAllUserInfo, fail to ReadAll Request.Body, err:%v.", err))
		ctx.ProtoBuf(http.StatusOK, &apiproto.GetAllUserInfoRsp{
			Code: 10300,
			Msg:  "请求参数读取失败",
			Data: nil,
		})
		return
	}

	getAllUserInfoReq := new(apiproto.GetAllUserInfoReq)
	err = proto.Unmarshal(bodyPro, getAllUserInfoReq)
	if err != nil {
		errInfo := fmt.Sprintf("GetAllUserInfo, fail to Unmarshal, bodyPro type error:%v", err)
		logger.Error(errInfo)
		ctx.ProtoBuf(http.StatusOK, &apiproto.GetAllUserInfoRsp{
			Code: 10400,
			Msg:  "请求参数格式错误，反序列化失败",
			Data: nil,
		})
		return
	}

	logger.Info(fmt.Sprintf("GetAllUserInfo req params ok, page:%d, size:%d.", getAllUserInfoReq.Page, getAllUserInfoReq.Size))

	//call gRPC server, get allUserInfo from db
	getAllUserInfoRsp, err := callGetAllUserInfoService(getAllUserInfoReq)
	if err != nil {
		logger.Error(fmt.Sprintf("GetAllUserInfo, fail to callGetAllUserInfoService, err:%v.", err))
		ctx.ProtoBuf(http.StatusOK, &apiproto.GetAllUserInfoRsp{
			Code: 10500,
			Msg:  "内部gRPC调用失败",
			Data: nil,
		})
		return
	}
	logger.Info(fmt.Sprintf("GetAllUserInfo, succeed to callGetAllUserInfoService. rsp: %v.", getAllUserInfoRsp))

	ctx.ProtoBuf(http.StatusOK, getAllUserInfoRsp)
}

func PingChat(ctx *gin.Context) {
	bodyPro, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("PingChat, fail to ReadAll Request.Body, err:%v.", err))
		ctx.ProtoBuf(http.StatusOK, &apiproto.PongChatRsp{
			Code: 10300,
			Msg:  "请求参数读取失败",
			Data: nil,
		})
		return
	}

	pingChatReq := new(apiproto.PingChatReq)
	err = proto.Unmarshal(bodyPro, pingChatReq)
	if err != nil {
		errInfo := fmt.Sprintf("PingChat, fail to Unmarshal, bodyPro type error:%v", err)
		logger.Error(errInfo)
		ctx.ProtoBuf(http.StatusOK, &apiproto.PongChatRsp{
			Code: 10400,
			Msg:  "请求参数格式错误，反序列化失败",
			Data: nil,
		})
		return
	}

	logger.Info(fmt.Sprintf("PingChat req params ok, InputParam: %v.", pingChatReq.InputParam))

	//call gRPC server, get PingChat response from server
	pongChatRsp, err := callPingChatService(pingChatReq)
	if err != nil {
		logger.Error(fmt.Sprintf("PingChat, fail to callPingChatService, err:%v.", err))
		ctx.ProtoBuf(http.StatusOK, &apiproto.PongChatRsp{
			Code: 10500,
			Msg:  "内部gRPC调用失败",
			Data: nil,
		})
		return
	}
	logger.Info(fmt.Sprintf("PingChat, succeed to callPingChatService. rsp: %v.", pongChatRsp))

	ctx.ProtoBuf(http.StatusOK, pongChatRsp)
}

//call restful api from user center
func verifyToken(token string) (bool, error) {
	logger.Info(fmt.Sprintf("verifyToken enter, token: %v.", token))

	return true, nil
}

func callGetUserInfoService(getUserInfoReq *apiproto.GetUserInfoReq) (*apiproto.GetUserInfoRsp, error) {
	logger.Info(fmt.Sprintf("callGetUserInfoService enter, input param getUserInfoReq: %v.", getUserInfoReq))

	//获取地址
	serverHost, err := zkmgr.GetServerHost()
	if err != nil {
		errInfo := fmt.Sprintf("callGetUserInfoService, get server host fail: %s \n", err)
		logger.Error(errInfo)
		return nil, errors.New(errInfo)
	}

	//serverHost := "127.0.0.1:8899"
	logger.Info(fmt.Sprintf("=========callGetUserInfoService, client succeed to connect host: " + serverHost))

	conn, err := grpc.Dial(serverHost, grpc.WithInsecure())
	if err != nil {
		errInfo := fmt.Sprintf("callGetUserInfoService, fail to grpc.Dial with err:%v.", err)
		logger.Warn(errInfo)
		return nil, errors.New(errInfo)
	}

	defer conn.Close()

	client := apiproto.NewApiServiceClient(conn)
	ret, err := client.GetUserInfo(context.Background(), getUserInfoReq)
	if err != nil {
		errInfo := fmt.Sprintf("callGetUserInfoService, fail to grpc call GetUserInfo with err:%v.", err)
		logger.Warn(errInfo)
		return nil, errors.New(errInfo)
	}

	return ret, nil
}

func callGetAllUserInfoService(getAllUserInfoReq *apiproto.GetAllUserInfoReq) (*apiproto.GetAllUserInfoRsp, error) {
	logger.Info(fmt.Sprintf("callGetAllUserInfoService enter, input param getAllUserInfoReq: %v.", getAllUserInfoReq))

	//获取地址
	serverHost, err := zkmgr.GetServerHost()
	if err != nil {
		errInfo := fmt.Sprintf("callGetAllUserInfoService, get server host fail: %s \n", err)
		logger.Error(errInfo)
		return nil, errors.New(errInfo)
	}

	//serverHost := "127.0.0.1:8899"
	logger.Info(fmt.Sprintf("=========callGetAllUserInfoService, client succeed to connect host: " + serverHost))

	conn, err := grpc.Dial(serverHost, grpc.WithInsecure())
	if err != nil {
		errInfo := fmt.Sprintf("callGetAllUserInfoService, fail to grpc.Dial with err:%v.", err)
		logger.Warn(errInfo)
		return nil, errors.New(errInfo)
	}

	defer conn.Close()

	client := apiproto.NewApiServiceClient(conn)
	ret, err := client.GetAllUserInfo(context.Background(), getAllUserInfoReq)
	if err != nil {
		errInfo := fmt.Sprintf("callGetAllUserInfoService, fail to grpc call GetAllUserInfo with err:%v.", err)
		logger.Warn(errInfo)
		return nil, errors.New(errInfo)
	}

	return ret, nil
}

func callPingChatService(pingChatReq *apiproto.PingChatReq) (*apiproto.PongChatRsp, error) {
	logger.Info(fmt.Sprintf("callPingChatService enter, input param pingChatReq: %v.", pingChatReq))

	//获取地址
	serverHost, err := zkmgr.GetServerHost()
	if err != nil {
		errInfo := fmt.Sprintf("callPingChatService, get server host fail: %s.", err)
		logger.Error(errInfo)
		return nil, errors.New(errInfo)
	}

	//serverHost := "127.0.0.1:8899"
	logger.Info(fmt.Sprintf("=========callPingChatService, client succeed to connect host: " + serverHost))

	conn, err := grpc.Dial(serverHost, grpc.WithInsecure())
	if err != nil {
		errInfo := fmt.Sprintf("callPingChatService, fail to grpc.Dial with err:%v.", err)
		logger.Warn(errInfo)
		return nil, errors.New(errInfo)
	}

	defer conn.Close()

	client := apiproto.NewApiServiceClient(conn)
	ret, err := client.PingChat(context.Background(), pingChatReq)
	if err != nil {
		errInfo := fmt.Sprintf("callPingChatService, fail to grpc call PingChat with err:%v.", err)
		logger.Warn(errInfo)
		return nil, errors.New(errInfo)
	}

	return ret, nil
}
