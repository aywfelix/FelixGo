package configure

import (
	"gopkg.in/ini.v1"
)

type iniConfig struct {
	handler *ini.File
}

func (c *iniConfig) Load(configPath string) error {
	var err error
	c.handler, err = ini.Load(configPath)
	if err != nil || c.handler == nil {
		return err
	}
	return nil
}

func (c *iniConfig) GetConfig() interface{} {
	return c.handler
}

func (c *iniConfig) GetInt(section string, key string) int {
	if c.handler == nil {
		return 0
	}
	return c.handler.Section(section).Key(key).MustInt(0)
}

func (c *iniConfig) GetString(section string, key string) string {
	return c.handler.Section(section).Key(key).MustString("")
}

func (c *iniConfig) GetBool(section string, key string) bool {
	return c.handler.Section(section).Key(key).MustBool(false)
}

func (c *iniConfig) GetFloat(section string, key string) float64 {
	return c.handler.Section(section).Key(key).MustFloat64(0)
}