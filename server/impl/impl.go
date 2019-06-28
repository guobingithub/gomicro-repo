package impl

import (
	"fmt"
	"golang.org/x/net/context"
	"gomicro-repo/server/config"
	"gomicro-repo/server/logger"
	"gomicro-repo/server/model"
	pb "gomicro-repo/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

func GetUserInfoPbResult(code int32, msg string, data *pb.Person) *pb.GetUserInfoRsp{
	resp := &pb.GetUserInfoRsp{
		Code: code,
		Msg: msg,
		Data: data,
	}

	return resp
}

func GetAllUserInfoPbResult(code int32, msg string, data *pb.AllPerson) *pb.GetAllUserInfoRsp{
	resp := &pb.GetAllUserInfoRsp{
		Code: code,
		Msg: msg,
		Data: data,
	}

	return resp
}

func PingChatPbResult(code int32, msg string, data *pb.PongRsp) *pb.PongChatRsp{
	resp := &pb.PongChatRsp{
		Code: code,
		Msg: msg,
		Data: data,
	}

	return resp
}

type DataService struct {

}

func (r *DataService) GetUserInfo(ctx context.Context, in *pb.GetUserInfoReq) (*pb.GetUserInfoRsp, error){
	userId := in.UserId
	logger.Info(fmt.Sprintf("GetUserInfo enter, input param userId:%v.",userId))

	var dao = new(model.DAO)
	dbData,err := dao.GetUserInfo(userId)
	if err!=nil{
		logger.Error(fmt.Sprintf("GetUserInfo, fail to search data from mongodb. err:%v.",err))
		return GetUserInfoPbResult(20200,"数据库操作失败！",nil),nil
	}

	var retData = new(pb.Person)
	retData.Id = string(dbData.Id)
	retData.Name = dbData.Name
	retData.Age = uint32(dbData.Age)

	logger.Info(fmt.Sprintf("GetUserInfo ok. retData=%v.",*retData))
	return GetUserInfoPbResult(20000,"GetUserInfo成功！",retData),nil
}

func (r *DataService) GetAllUserInfo(ctx context.Context, in *pb.GetAllUserInfoReq) (*pb.GetAllUserInfoRsp, error){
	logger.Info(fmt.Sprintf("GetAllUserInfo enter, input param in:%v.",in))

	if in.Page<1 || in.Size<=0{
		errInfo := fmt.Sprintf("GetAllUserInfo, input param error! page:%d, size:%d.", in.Page,in.Size)
		logger.Error(errInfo)
		return GetAllUserInfoPbResult(20100,"参数不合法！",nil),nil
	}

	var dao = new(model.DAO)
	dbData,err := dao.GetAllUserInfo(in.Page,in.Size)
	if err!=nil{
		errInfo := fmt.Sprintf("GetAllUserInfo, fail to search data from mongodb. err:%v.",err)
		logger.Error(errInfo)
		return GetAllUserInfoPbResult(20200,"数据库操作失败！",nil),nil
	}

	var (
		retData = new(pb.AllPerson)
		retDataPer = make([]*pb.Person, len(dbData))
	)
	for i,v := range dbData{
		logger.Error(fmt.Sprintf("idddddd====%v",v.Id.Hex()))
		retDataPer[i] = new(pb.Person)
		retDataPer[i].Id = v.Id.Hex()
		retDataPer[i].Name = v.Name
		retDataPer[i].Age = uint32(v.Age)
	}
	retData.Per = retDataPer

	logger.Info(fmt.Sprintf("GetAllUserInfo ok. length(retDataPer)=%d.",len(retDataPer)))
	return GetAllUserInfoPbResult(20000,"GetAllUserInfo成功！",retData),nil
}

func (r *DataService) PingChat(ctx context.Context, in *pb.PingChatReq) (*pb.PongChatRsp, error){
	logger.Info(fmt.Sprintf("PingChat enter, input param in:%v.",in.InputParam))

	var (
		retData = new(pb.PongRsp)
		fid,tid int32
		fname,tname,pong string
	)

	fid = in.InputParam.Fid
	fname = in.InputParam.Fname
	tid = in.InputParam.Tid
	tname = in.InputParam.Tname
	pong = fmt.Sprintf("I am ok, my id is %d, thank you %s!",tid,fname)
	fid,fname,tid,tname = tid,tname,fid,fname

	retData.Fid = fid
	retData.Fname = fname
	retData.Tid = tid
	retData.Tname = tname
	retData.Pong = pong

	logger.Info(fmt.Sprintf("PingChat ok. retData=%v.",retData))
	return PingChatPbResult(20000,"PingChat成功！",retData),nil
}

func registerGRpcServer(s *grpc.Server) {
	pb.RegisterApiServiceServer(s, &DataService{})
}

func StartGrpcServer() {
	addr := config.GetConf().Server.Address+config.GetConf().Server.GrpcPort
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		grpclog.Fatalf("StartGrpcServer, failed to listen: %v", err)
	}

	// 实例化grpc Server
	s := grpc.NewServer()

	//注册
	registerGRpcServer(s)

	println("StartGrpcServer, Grpc server Listen on ------" + addr)

	if err = s.Serve(listen);err!=nil{
		logger.Error(fmt.Sprintf("StartGrpcServer, fail to listen!"))
		listen.Close()
	}
}
