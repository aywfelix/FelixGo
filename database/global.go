package database

import (
	"errors"
)

// 数据表
type TableNameType string

const (
	TB_User TableNameType = "user"
)

var (
	QueryRowNil = errors.New("query row nil")
)
