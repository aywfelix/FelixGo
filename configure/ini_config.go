package configure

import (
	"gopkg.in/ini.v1"
)

type iniConfig struct {
	cfg *ini.File
}

func (c iniConfig) Load(configPath string) error {
	cfg, err := ini.Load(configPath)
	if err != nil {
		return err
	}
	c.cfg = cfg
	return nil
}

func (c iniConfig) GetConfig() interface{} {
	return c.cfg
}

func (c iniConfig) GetInt(section string, key string) int {
	return c.cfg.Section(section).Key(key).MustInt()
}

func (c iniConfig) GetString(section string, key string) string {
	return c.cfg.Section(section).Key(key).MustString("")
}

func (c iniConfig) GetBool(section string, key string) bool {
	return c.cfg.Section(section).Key(key).MustBool(false)
}

func (c iniConfig) GetFloat(section string, key string) float64 {
	return c.cfg.Section(section).Key(key).MustFloat64()
}