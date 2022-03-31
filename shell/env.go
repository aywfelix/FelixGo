package shell

import (
	"os"
	"strings"
)

type env struct{}

var Env env

func (e env) All() []string {
	return os.Environ()
}

func (e env) Map() map[string]string {
	m := make(map[string]string)
	i := 0
	for _, s := range os.Environ() {
		i = strings.IndexByte(s, '=')
		m[s[0:i]] = s[i+1:]
	}
	return m
}

func (e env) Get(key string, def ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok && len(def) > 0 {
		return def[0]
	}
	return v
}

func (e env) Set(key, value string) error {
	return os.Setenv(key, value)
}

func (e env) SetMap(m map[string]string) error {
	for k, v := range m {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (e env) Contains(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

func (e env) Remove(key ...string) error {
	var err error
	for _, v := range key {
		err = os.Unsetenv(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e env) Build(m map[string]string) []string {
	array := make([]string, len(m))
	index := 0
	for k, v := range m {
		array[index] = k + "=" + v
		index++
	}
	return array
}
