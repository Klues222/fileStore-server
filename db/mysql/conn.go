package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/fileserver?charset=utf8")
	if err != nil {
		fmt.Println("fail to connect to mysql, err : " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("fail to connect to mysql, err : " + err.Error())
		os.Exit(1)
	}
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

// 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}
