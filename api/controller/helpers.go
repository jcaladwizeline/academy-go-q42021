package controller

import (
	"encoding/json"
	"net/http"
)

func responseHandle(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}
	js, _ := json.Marshal(data)

	w.Write(js)
}
