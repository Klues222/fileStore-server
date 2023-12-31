package handler

import (
	"encoding/json"
	"filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt = "*#890"
)

/*
登陆接口
*/
func SignupHander(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}
	enc_password := util.Sha1([]byte(password + pwd_salt))
	suc := db.UserSignUp(username, enc_password)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}

}

// SignInHandler: 登陆接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// data, err := ioutil.ReadFile("./static/view/signin.html")
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		// w.Write(data)
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPassword := util.Sha1([]byte(password + pwd_salt))
	//1.校验用户名及密码
	pwdChecked := db.UserSignin(username, encPassword)
	if !pwdChecked {
		esp1 := util.RespMsg{
			Code: 0,
			Msg:  "FAILED",
		}
		b1, _ := json.Marshal(esp1)
		io.WriteString(w, string(b1))
		return
	}
	//2.生成访问凭证（token）
	token := GenToken(username)
	upRes := db.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}
	//3。登陆成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	//http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
	w.Header().Set("Authorization", token)
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	b, _ := json.Marshal(resp)
	io.WriteString(w, string(b))
	//w.Write(resp.JSONBytes())
}

// 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	token := r.Header.Get("Authorization")
	print(token)
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")
	//2.验证token是否有效
	//isValidToken := IsTokenValid(token)
	//if !isValidToken {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	//3.查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//4.组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	//40位字符: md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// token是否有效
func IsTokenValid(token string) bool {
	//TODO:判断token的时效性，是否过期
	if len(token) != 40 {
		return false
	}
	//TODO:从数据表tbl_user_token查询username对应的token信息
	//TODO:对比两个token是否一致
	return true
}
