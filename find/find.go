package find

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var sourceURL = "http://ps20.software.eu.playstation.com/"
var redirectURL = "http://ps20.software.eu.playstation.com/redirect.php"
var httpHeaders = map[string]string{
	"User-Agent": "Mozilla/5.0 " +
		"(Macintosh; Intel Mac OS X 10.10; rv:34.0) " +
		"Gecko/20100101 Firefox/34.0",
}

func httpGet(url string) (string, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)

	for header, value := range httpHeaders {
		req.Header.Add(header, value)
	}

	response, err := client.Do(req)
	body := ""

	if err != nil {
		return "", err
	} else {
		defer response.Body.Close()
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		body = string(content)
	}

	return body, nil
}

func getPageSource() string {
	response, err := httpGet(sourceURL)

	if err != nil {
		log.Fatal(err)
	}

	return response
}

func GetSp() string {
	sp := ""
	r, _ := regexp.Compile("config\\.sp = \"(.+)\"")

	found := r.FindStringSubmatch(getPageSource())
	if len(found) > 1 {
		sp = found[1]
	}

	return sp
}

func GetSecretURL(sp string) string {
	if sp == "" {
		sp = GetSp()
	}

	secretURL, err := httpGet(redirectURL + "?sp=" + sp)

	if err != nil {
		log.Fatal(err)
	}

	return secretURL
}

func makeLookupURL(sp string) string {
	return redirectURL + "?sp=" + sp
}

// Source outputs page source to STDOUT.
func Source() {
	fmt.Println(getPageSource())
}

// Sp outputs SP details to STDOUT.
func Sp(sp string) {
	if sp == "" {
		sp = GetSp()
	}

	fmt.Println("SP Code: " + sp)
}

// RedirectURL outputs Redirect URL details to STDOUT.
func RedirectURL(sp string) {
	if sp == "" {
		sp = GetSp()
	}

	fmt.Println("Redirect URL: " + makeLookupURL(sp))
}

// Secret output secret URL details to STDOUT.
func Secret(sp string) {
	if sp == "" {
		sp = GetSp()
	}

	fmt.Println("Secret URL: " + GetSecretURL(sp))
}

// Details outputs a summary of SP, Redirect, and Secret URL to STDOUT.
func Details() {
	sp := GetSp()
	Sp(sp)
	RedirectURL(sp)
	Secret(sp)
}
