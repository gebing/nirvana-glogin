package db

import (
	log "github.com/castbox/nirvana-gcore/glog"
	"github.com/castbox/nirvana-gcore/gmongo"
	"glogin/config"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo() {
	//gmongo.Init(config.MongoUrl())
	InitAccount()
	InitVerifyCode()
	InitForBusiness()
}

func LoadOne(filter interface{}, result interface{}, tableName string) (err error) {
	doc, errFind := gmongo.FindOne(config.MongoUrl(), config.MongoDb(), tableName, filter)
	if errFind != nil {
		log.Warnw("Table LoadOne", "err", err)
		err = errFind
		return
	}
	if errDecode := doc.Decode(result); errDecode != nil {
		err = errDecode
	}
	return
}

func CheckNotExist(filter interface{}, tableName string) bool {
	count, err := gmongo.CountDocuments(config.MongoUrl(), config.MongoDb(), tableName, filter)
	if err != nil {
		log.Warnw("CheckNotExist", "err", err)
		return false
	}
	if count == 0 {
		return true
	}
	return false
}

func UpdateOne(filter interface{}, update interface{}, tableName string) (err error) {
	_, errUpdate := gmongo.UpdateOne(config.MongoUrl(), config.MongoDb(), tableName, filter, update)
	if errUpdate != nil {
		log.Warnw("UpdateOne Table error", "tableName", tableName, "errUpdate", errUpdate)
		err = errUpdate
		return
	}
	return
}

func UpdateOne_Upsert(filter interface{}, update interface{}, tableName string) (err error) {
	_, errUpdate := gmongo.UpdateOne(config.MongoUrl(), config.MongoDb(), tableName, filter, update, options.Update().SetUpsert(true))
	if errUpdate != nil {
		log.Warnw("UpdateOne Table error", "tableName", tableName, "errUpdate", errUpdate)
		err = errUpdate
		return
	}
	return
}

// 账号数量
func AccountCount(filter interface{}, tableName string) int64 {
	count, err := gmongo.CountDocuments(config.MongoUrl(), config.MongoDb(), tableName, filter)
	if err != nil {
		log.Warnw("AccountCount", "err", err)
		return -1
	}
	return count
}

func DeleteOne(filter interface{}, tableName string) (err error) {
	n, err := gmongo.DeleteOne(config.MongoUrl(), config.MongoDb(), tableName, filter)
	log.Infof("DeleteOne Table count %d", n)
	if err != nil {
		log.Infof("DeleteOne Table error, tableName %v, errDeleteOne %v", tableName, err)
	}
	return
}
