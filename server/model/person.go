package model

import (
	"errors"
	"fmt"
	"gomicro-repo/server/constants"
	"gomicro-repo/server/db"
	"gomicro-repo/server/entity"
	"gomicro-repo/server/logger"
	"gopkg.in/mgo.v2/bson"
)

type DAO struct {}

func (d DAO) GetUserInfo(userId string) (*entity.PersonMG, error){
	logger.Info(fmt.Sprintf("DAO, GetUserInfo enter, userId:%s.",userId))

	var (
		err error
		personMg *entity.PersonMG
	)

	collection := db.GetMs().DB(constants.Db_GoMicro).C(constants.Collection_Person)
	//err = collection.Find(bson.M{"_id": bson.ObjectIdHex(userId)}).One(personMg)
	err = collection.FindId(bson.ObjectIdHex(userId)).One(personMg)
	if err!=nil{
		errInfo := fmt.Sprintf("DAO, GetUserInfo, fail to collection.Find by userId:%v. err:%v.",userId,err)
		logger.Error(errInfo)
		return nil,errors.New(errInfo)
	}

	logger.Info(fmt.Sprintf("DAO, GetUserInfo, find person, info: %v.",personMg))
	return personMg, nil
}

func (d DAO) GetAllUserInfo(page,size int32) ([]entity.PersonMG, error){
	logger.Info(fmt.Sprintf("DAO, GetAllUserInfo enter, page:%d, size:%d.",page,size))

	var (
		err error
		personsMG []entity.PersonMG
	)

	collection := db.GetMs().DB(constants.Db_GoMicro).C(constants.Collection_Person)
	err = collection.Find(bson.M{ }).
						Skip(int((page-1)*size)).
						Limit(int(size)).
						All(&personsMG)
	if err!=nil{
		errInfo := fmt.Sprintf("DAO, batch GetAllUserInfo, fail to collection.Find with page:%d, size:%d. err:%v.",page,size,err)
		logger.Error(errInfo)
		return nil,errors.New(errInfo)
	}

	logger.Info(fmt.Sprintf("DAO, GetAllUserInfo, find person list, len(personsMG):%d \n, personsMG: %v.",len(personsMG),personsMG))

	return personsMG, nil
}
