package handler

import (
	"dbaas/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
	conditions := r.URL.Path
	err := db.Read("gopgx_schema.Users", conditions)
	if err != nil {
		fmt.Println(err)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	str_body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Sorry cant do")
		return
	}
	json.Unmarshal(str_body, &body)
	fmt.Println(body)
	db.Insert("gopgx_schema.Users", body)
}
