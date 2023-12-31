package handler

import (
	"io"
	"net/http"
	"os"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("./static/view/signin.html")
	if err != nil {
		io.WriteString(w, "internet server error")

	}
	io.WriteString(w, string(data))
	
}
