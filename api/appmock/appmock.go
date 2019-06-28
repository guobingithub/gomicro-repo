package main

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gomicro-repo/api/logger"
	apiproto "gomicro-repo/api/proto"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	logger.Info(fmt.Sprintf("gomicro-api appmock start..."))

	var (
		err  error
		isOK = true
	)

	for {
		_, err = pingChat()
		if err != nil {
			logger.Error(fmt.Sprintf("fail to pingChat, err:%v.", err))
			isOK = false
		} else {
			logger.Info("ok to pingChat...\n")
		}

		_, err = getUserInfo()
		if err != nil {
			logger.Error(fmt.Sprintf("fail to getUserInfo, err:%v.", err))
			isOK = false
		} else {
			logger.Info("ok to getUserInfo...\n")
		}

		_, err = getAllUserInfo()
		if err != nil {
			logger.Error(fmt.Sprintf("fail to getAllUserInfo, err:%v.", err))
			isOK = false
		} else {
			logger.Info("ok to getAllUserInfo...\n")
		}

		logger.Info(fmt.Sprintf("gomicro-api appmock, api req result:%v.\n", isOK))

		time.Sleep(25 * time.Second)
	}
}

func pingChat() (rsp apiproto.PongChatRsp, err error) {
	logger.Error(fmt.Sprintf("pingChat enter..."))
	http.DefaultClient.Timeout = time.Second * 3
	pingUrlv1 := "http://127.0.0.1:8000/ping"

	pingChatReq := new(apiproto.PingChatReq)
	pingChatReq.InputParam = &apiproto.PingReq{
		Fid:   100,
		Fname: "张三",
		Tid:   101,
		Tname: "李四",
		Ping:  "hello gopher101, i am 100. Are you ok?",
	}

	// 对数据进行序列化
	pingChatReqProto, err := proto.Marshal(pingChatReq)
	if err != nil {
		log.Fatalln("pingChat, Mashal pingChatReq error:", err)
	}

	logger.Error(fmt.Sprintf("pingChat, Mashal ok. pingChatReqProto:%s", string(pingChatReqProto)))
	request, err := http.NewRequest(http.MethodPost, pingUrlv1, bytes.NewBuffer(pingChatReqProto))
	if err != nil {
		return
	}

	request.Header.Add("Authorization", "tok12345")
	request.Header.Add("Content-Type", "application/x-protobuf")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	rspBodyPro, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = proto.Unmarshal(rspBodyPro, &rsp)
	if err != nil {
		log.Fatalln("pingChat, UnMashal rspBodyPro error:", err)
		return
	}

	logger.Info(fmt.Sprintf("pingChat, rsp info:%v.", rsp))
	if rsp.Code != 200 {
		logger.Error(fmt.Sprintf("pingChat, get rsp code error. code:%d, msg:%s.", rsp.Code, rsp.Msg))
		return
	}

	logger.Info(fmt.Sprintf("Succeed to pingChat: %v.", rsp.Data))
	return
}

func getUserInfo() (rsp apiproto.GetUserInfoRsp, err error) {
	logger.Error(fmt.Sprintf("getUserInfo enter..."))
	http.DefaultClient.Timeout = time.Second * 3
	pingUrlv1 := "http://127.0.0.1:8000/get/userInfo"

	getUserInfoReq := new(apiproto.GetUserInfoReq)
	getUserInfoReq.UserId = "5d1461456ddf2c929979c67d"

	// 对数据进行序列化
	getUserInfoReqProto, err := proto.Marshal(getUserInfoReq)
	if err != nil {
		log.Fatalln("getUserInfo, Mashal getUserInfoReq error:", err)
	}

	logger.Error(fmt.Sprintf("getUserInfo, Mashal ok. getUserInfoReqProto:%s", string(getUserInfoReqProto)))
	request, err := http.NewRequest(http.MethodPost, pingUrlv1, bytes.NewBuffer(getUserInfoReqProto))
	if err != nil {
		return
	}

	request.Header.Add("Authorization", "tok12345")
	request.Header.Add("Content-Type", "application/x-protobuf")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	rspBodyPro, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = proto.Unmarshal(rspBodyPro, &rsp)
	if err != nil {
		log.Fatalln("getUserInfo, UnMashal rspBodyPro error:", err)
		return
	}

	logger.Info(fmt.Sprintf("getUserInfo, rsp info:%v.", rsp))
	if rsp.Code != 200 {
		logger.Error(fmt.Sprintf("getUserInfo, get rsp code error. code:%d, msg:%s.", rsp.Code, rsp.Msg))
		return
	}

	logger.Info(fmt.Sprintf("Succeed to getUserInfo: %v.", rsp.Data))
	return
}

func getAllUserInfo() (rsp apiproto.GetAllUserInfoRsp, err error) {
	logger.Error(fmt.Sprintf("getAllUserInfo enter..."))
	http.DefaultClient.Timeout = time.Second * 3
	pingUrlv1 := "http://127.0.0.1:8000/get/all/userInfo"

	getAllUserInfoReq := new(apiproto.GetAllUserInfoReq)
	getAllUserInfoReq.Page = 1
	getAllUserInfoReq.Size = 5

	// 对数据进行序列化
	getAllUserInfoReqProto, err := proto.Marshal(getAllUserInfoReq)
	if err != nil {
		log.Fatalln("getAllUserInfo, Mashal getAllUserInfoReq error:", err)
	}

	logger.Error(fmt.Sprintf("getAllUserInfo, Mashal ok. getAllUserInfoReqProto:%s", string(getAllUserInfoReqProto)))
	request, err := http.NewRequest(http.MethodPost, pingUrlv1, bytes.NewBuffer(getAllUserInfoReqProto))
	if err != nil {
		return
	}

	request.Header.Add("Authorization", "tok12345")
	request.Header.Add("Content-Type", "application/x-protobuf")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	rspBodyPro, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = proto.Unmarshal(rspBodyPro, &rsp)
	if err != nil {
		log.Fatalln("getAllUserInfo, UnMashal rspBodyPro error:", err)
		return
	}

	logger.Info(fmt.Sprintf("getAllUserInfo, rsp info:%v.", rsp))
	if rsp.Code != 200 {
		logger.Error(fmt.Sprintf("getAllUserInfo, get rsp code error. code:%d, msg:%s.", rsp.Code, rsp.Msg))
		return
	}

	logger.Info(fmt.Sprintf("Succeed to getAllUserInfo: %v.", rsp.Data))
	return
}
