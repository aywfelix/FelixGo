package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var (
	errPathNotExist = errors.New("config path is not exit")
)

type IConfigure interface {
	Load(configPath string) error
	AssignStruct(st interface{}) error
	GetConfig() interface{}
}

type configure struct {
	st   interface{}
	data []byte
}

func (c configure) Load(configPath string) error {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}
	configPath = pwd + "/" + configPath
	if isExists := File.Exists(configPath); !isExists {
		return errPathNotExist
	}
	c.data, err = ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	return nil
}

func (c configure) AssignStruct(st interface{}) error {
	if c.data == nil {
		return errors.New("call load function first")
	}
	if err := json.Unmarshal(c.data, st); err != nil {
		return err
	}
	c.st = st
	return nil
}

func (c configure) GetConfig() interface{} {
	return c.st
}

//=====================================================================================

var Configure configure
