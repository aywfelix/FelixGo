package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	MAX_OPEN_CONNS = 256
	MAX_CONN_TIME  = 24 * 3600
	MAX_IDLE_CONNS = 32
	MAX_IDLE_TIME  = 5 * 60
)

type IMysql interface {
	Connect(user, password, ip string, port int, db string) error
	IsConnected() bool
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close()
}

type Mysql struct {
	handler *sql.DB
}

func NewMysql() *Mysql {
	return &Mysql{
		handler: nil,
	}
}

func (m *Mysql) Connect(user, password, ip string, port int, db string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, ip, port, db)
	handler, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	m.handler = handler
	m.handler.SetMaxOpenConns(MAX_OPEN_CONNS)
	m.handler.SetConnMaxLifetime(time.Second * MAX_CONN_TIME)
	m.handler.SetMaxIdleConns(MAX_IDLE_CONNS)
	m.handler.SetConnMaxIdleTime(time.Second * MAX_IDLE_TIME)
	return nil
}

func (m *Mysql) IsConnected() bool {
	if err := m.handler.Ping(); err != nil {
		fmt.Println("mysql connection error:", err.Error())
		return false
	}
	return true
}

func (m *Mysql) Close() {
	m.handler.Close()
	m.handler = nil
}

// insert/update/delete
func (m *Mysql) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := m.handler.Exec(query, args...)
	if err != nil {
		// fmt.Sprintf("exec sql error, sql:%s, error:%s\n", query, err.Error())
		return nil, err
	}
	return result, nil
}

// lines query
func (m *Mysql) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := m.handler.Query(query, args...)
	if err != nil {
		// fmt.Sprintf("query sql error, sql:%s, error:%s\n", query, err.Error())
		return rows, err
	}
	return rows, nil
}

// line query
func (m *Mysql) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.handler.QueryRow(query, args...)
}

func (m *Mysql) Stats() *sql.DBStats {
	if m.handler == nil {
		return nil
	}

	stats := m.handler.Stats()
	return &stats
}

// mysql prepared statement:m.handler.Prepare(sql)
