package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	sp := r.URL.Query().Get("sp")
	response := "0"

	if sp == "SA" {
		response = "http://foo.bar/hello.php"
	} else if sp == "RU" {
		response = "http://foo.bar/world.php"
	} else if sp == "TX" {
		response = "http://foo.bar/dude.php"
	} else if sp == "vK" {
		response = "http://foo.bar/wtf.php"
	} else if sp == "SBC" {
		response = "http://foo.bar/MSSDXF9394X.php"
	}

	fmt.Fprint(w, response)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/redirect.php", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
