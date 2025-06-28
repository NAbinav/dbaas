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
	db.Insert("GODB", body)
}
