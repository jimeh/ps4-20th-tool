package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {
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

func secretHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Congratulations, you have been registered!")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/redirect.php", redirectHandler)
	http.HandleFunc("/secret.php", secretHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
