package database

import (
	"fmt"
	protos "github.com/felix/felixgo/database/proto"
	. "github.com/felix/felixgo/db"
	"github.com/golang/protobuf/proto"
)

type UserHandler struct {
	protos.User
}

func (r *UserHandler) PrimaryKey() (interface{}, interface{}) {
	return r.UserId, nil
}

func (r *UserHandler) CreateTable(mysql *Mysql) error {
	sql := `CREATE TABLE IF NOT EXISTS %s (
  Id bigint(20) NOT NULL ,
  AccountId varchar(45) DEFAULT '',
  Name varchar(255) DEFAULT NULL,
  Level varchar(255) DEFAULT NULL,
  BaseData varchar(255) DEFAULT NULL,
  DetailData varchar(255) DEFAULT NULL,
  CreateTime bigint(32) DEFAULT NULL,
  PRIMARY KEY (Id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`
	sql = fmt.Sprintf(sql, TB_User)
	_, err := mysql.Exec(sql)
	return err
}

func (r *UserHandler) Select(mysql *Mysql, keys ...interface{}) (bool, error) {
	sql := `select Id, AccountId, Name, Level, BaseData, CreateTime from %s where Id=? `
	sql = fmt.Sprintf(sql, TB_User)
	row := mysql.QueryRow(sql, keys...)
	if row == nil {
		return true, QueryRowNil
	}

	tmpBaseData := make([]byte, 0)
	tmpDetailData := make([]byte, 0)
	err := row.Scan(&r.UserId, &r.AccountId, &r.Name, &r.Level, &tmpBaseData, &tmpDetailData, &r.CreateTime)
	if err != nil {
		return false, err
	}

	err = proto.Unmarshal(tmpBaseData, r.BaseData)
	if err != nil {
		return false, err
	}

	err = proto.Unmarshal(tmpDetailData, r.DetailData)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserHandler) Insert(mysql *Mysql, args ...interface{}) (bool, error) {
	sql := `insert into %s values (?, ?, ?, ?, ?, ?, ?)`
	Id := args[0].(int)
	AccountId := args[1].(string)
	Name := args[2].(string)
	Level := args[3].(int)
	tmpBaseData := new(protos.UserBaseData)
	BaseData, _ := proto.Marshal(tmpBaseData)
	tmpDetailData := new(protos.UserDetailData)
	DetailData, _ := proto.Marshal(tmpDetailData)
	sql = fmt.Sprintf(sql, TB_User)
	_, err := mysql.Exec(sql, Id, AccountId, Name, Level, BaseData, DetailData)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserHandler) Updata(mysql *Mysql, args ...interface{}) (bool, error) {
	sql := `update %s set AccountId=?, Name=?, Level=?, BaseData=?, DetailData=? where id=?`
	sql = fmt.Sprintf(sql, TB_User)
	Id := args[0].(int)
	AccountId := args[1].(string)
	Name := args[2].(string)
	Level := args[3].(int)
	tmpBaseData := new(protos.UserBaseData)
	BaseData, _ := proto.Marshal(tmpBaseData)
	tmpDetailData := new(protos.UserDetailData)
	DetailData, _ := proto.Marshal(tmpDetailData)
	_, err := mysql.Exec(sql, Id, AccountId, Name, Level, BaseData, DetailData)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserHandler) Delete(mysql *Mysql, keys ...interface{}) (bool, error) {
	sql := `delete from %s where id=?`
	sql = fmt.Sprintf(sql, TB_User)
	_, err := mysql.Exec(sql, keys...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserHandler) FromBytes(blob []byte) error {
	err := proto.Unmarshal(blob, r)
	return err
}

func (r *UserHandler) ToBytes() ([]byte, error) {
	blob, err := proto.Marshal(r)
	if err != nil {
		// TODO: log
		return nil, err
	}
	return blob, nil
}
