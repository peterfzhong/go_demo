package main

import (
	"database/sql"
	"fmt"
	//"errors"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlUtil struct {
	user      string
	password  string
	ip        string
	port      int
	db_name   string
	db 		  *sql.DB
}

func (mysql* MysqlUtil) Close()(){
	mysql.db.Close()
}

func (mysql* MysqlUtil) Init()(){
	var err error
	sql_conn_info := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", mysql.user, mysql.password, mysql.ip, mysql.port, mysql.db_name)
	fmt.Println(sql_conn_info)
	mysql.db, err = sql.Open("mysql",  sql_conn_info)

	if err != nil{
		fmt.Println("err in conn ", sql_conn_info, err)
	}

	return
}

func (mysql *MysqlUtil) Exec(sql string)(code int){
	code = 0

	res, err := mysql.db.Exec(sql)
	if err != nil{
		fmt.Println("error in exec", sql, err)
	}
	//mysql.db.

	fmt.Println(res.LastInsertId())
	fmt.Println(res.RowsAffected())
	return
}

func (mysql* MysqlUtil) Query(sql string)(rows *sql.Rows){
	rows ,err := mysql.db.Query(sql)
	if err != nil{
		fmt.Println("error in query ", sql)
		return
	}

	return
}

func Test()(){
	mysql := &MysqlUtil{"dev", "123456", "127.0.0.1", 3306, "db_staff", nil}

	mysql.Init()
	defer mysql.Close()

	sql := "insert into db_staff.db_staff set name = 'peter', addr = '深圳南山', salary = 10000"
	mysql.Exec(sql)

	rows := mysql.Query("select * from db_staff")
	for rows.Next() {
		var id int
		var username string
		var addr string
		var salary float32
		err := rows.Scan(&id, &username, &addr, &salary)

		if err != nil{
			fmt.Println("error in fetch rows")
		}
		fmt.Println(id, username, addr, salary)
	}
}
