package configure

import "errors"

var (
	errPathNotExist = errors.New("config path is not exit")
)


var JsonConfig jsonConfig
var IniConfig iniConfig