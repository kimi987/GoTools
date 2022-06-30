package mysqldb

import (
	"database/sql"
	"errors"
	"fmt"

	"dgServer/log"

	_ "github.com/go-sql-driver/mysql"
)

type Dbs struct {
	DbDetails map[string]*Db_detail
}

type Db_detail struct {
	Name        string
	ConnAddr    string
	MaxOpenConn int
	MaxIdleConn int
	Db          *sql.DB
}

//Run 执行Module模块
func (this *Dbs) Run(closeSig chan bool) {
	for _, v := range this.DbDetails {
		d, err := sql.Open("mysql", v.ConnAddr)

		if err != nil {
			log.Error("open DB fail %s error = %v", v.ConnAddr, err)
			continue
		}
		v.Db = d

		v.Db.SetMaxOpenConns(v.MaxOpenConn)
		v.Db.SetMaxIdleConns(v.MaxIdleConn)
		err = v.Db.Ping()

		if err != nil {
			log.Error("connect DB fail %s error = %v", v.ConnAddr, err)
			continue
		}
		fmt.Println("connection to db[%s] success", v.Name)
	}
	<-closeSig
}

//执行Module模块
func (this *Dbs) OnDestroy() {
	for _, v := range this.DbDetails {
		if v.Db != nil {
			v.Db.Close()
		}
	}
}

//通过query 来查询数据
func (this *Db_detail) QuerySelect(query string, values ...interface{}) (result_map map[int]map[string]string, err error) {
	Db := this.Db
	if Db == nil {
		log.Error("query select error mysql is not connected")
		return make(map[int]map[string]string), errors.New("query error")
	}

	query = fmt.Sprintf(query, values...)
	rows, err := Db.Query(query)
	defer rows.Close()

	if err != nil {
		log.Error("query %s error = %v", query, err)
		return make(map[int]map[string]string), errors.New("query error")
	}

	result := make(map[int]map[string]string)

	colums, _ := rows.Columns()

	valuesargs := make([]sql.RawBytes, len(colums))
	scanargs := make([]interface{}, len(valuesargs))

	for i := range valuesargs {
		scanargs[i] = &valuesargs[i]
	}

	var n = 1

	for rows.Next() {
		result[n] = make(map[string]string)

		err := rows.Scan(scanargs...)

		if err != nil {
			log.Error("scan rows error = %v", query, err)
			return make(map[int]map[string]string), errors.New("query error")
		}

		for i, v := range valuesargs {
			result[n][colums[i]] = string(v)
		}

		n++
	}
	return result, nil
}

//通过query 来增加数据
func (this *Db_detail) QueryInsert(query string, values ...interface{}) (id int, err error) {
	Db := this.Db
	if Db == nil {
		return 0, errors.New("mysql is not connected")
	}
	query = fmt.Sprintf(query, values...)
	res, err := Db.Exec(query)
	if err != nil {

		return 0, err
	}

	index, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(index), nil

}

//通过query 来更新数据
func (this *Db_detail) QueryUpdate(query string, values ...interface{}) (num int, err error) {

	Db := this.Db
	if Db == nil {
		return 0, errors.New("mysql is not connected")
	}
	query = fmt.Sprintf(query, values...)
	res, err := Db.Exec(query)
	if err != nil {
		return 0, err
	}

	af_num, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	return int(af_num), nil
}

//通过query 更新事务的数据
func (this *Db_detail) QueryUpdateTrans(query string, values ...interface{}) (num int, err error) {

	Db := this.Db
	if Db == nil {
		return 0, errors.New("mysql is not connected")
	}
	tx, err := Db.Begin()
	defer tx.Rollback()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(values...)

	af_num, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	return int(af_num), nil
}

//通过query 来删除数据
func (this *Db_detail) QueryDelete(query string, values ...interface{}) (num int, err error) {

	Db := this.Db
	if Db == nil {
		return 0, errors.New("mysql is not connected")
	}
	query = fmt.Sprintf(query, values...)
	res, err := Db.Exec(query)

	if err != nil {
		return 0, err
	}
	af_num, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	return int(af_num), nil
}
