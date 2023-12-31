package db

import (
	"filestore-server/db/mysql"
	"fmt"
)

func UserSignUp(username string, passwd string) bool {
	stmt, err := mysql.DBConn().Prepare("insert into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("fail to insert, err : " + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("fail to insert, err : " + err.Error())
		return false
	}
	//校验是否重复注册
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

// UserSignin:登陆验证
func UserSignin(username string, encpwd string) bool {
	stmt, err := mysql.DBConn().Prepare("select * from tbl_user where user_name = ? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found: " + username)
		return false
	}

	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// 刷新用户登陆的token
func UpdateToken(username string, token string) bool {
	stmt, err := mysql.DBConn().Prepare("replace into tbl_user_token (`user_name`,`user_token`) values (?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mysql.DBConn().Prepare("select user_name,signup_at from tbl_user where user_name = ? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
