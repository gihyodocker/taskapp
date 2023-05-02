package db

import "fmt"

type Datasource interface {
	Driver() string
	DSN() string
}

type mysqlDatasource struct {
	username string
	password string
	host     string
	dbname   string
}

func NewMySQLDatasource(username, password, host, dbname string) Datasource {
	return &mysqlDatasource{
		username: username,
		password: password,
		host:     host,
		dbname:   dbname,
	}
}

func (d *mysqlDatasource) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", d.username, d.password, d.host, d.dbname)
}

func (d mysqlDatasource) Driver() string {
	return "mysql"
}
