package configure

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	. "github.com/aywfelix/felixgo/utils"
)

type jsonConfig struct {
	st   interface{}
	data []byte
}

func (c jsonConfig) Load(configPath string) error {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	configPath = cwd + "/" + configPath
	if isExists := File.Exists(configPath); !isExists {
		return errPathNotExist
	}
	var data []byte
	data, err = ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	c.data = data
	return nil
}

func (c jsonConfig) AssignStruct(st interface{}) error {
	if c.data == nil {
		return errors.New("call load function first")
	}
	if err := json.Unmarshal(c.data, st); err != nil {
		return err
	}
	c.st = st
	return nil
}

func (c jsonConfig) GetConfig() interface{} {
	return c.st
}
