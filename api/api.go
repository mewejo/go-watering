package api

import (
	"log"
	"net/http"
)

func StartApi() {
	http.HandleFunc("/api/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
