package auto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jimeh/ps4-20th-tool/find"
	"github.com/jinzhu/now"
)

var testURL = "http://localhost:4010/secret.php"
var clueURLs = []string{
	"https://twitter.com/PlayStationUK",
	"https://twitter.com/GAMEdigital",
}

type Config struct {
	CurrentSp   string `json:"current_sp"`
	SubmitDelay int    `json:"submit_delay"`
	Forms       []Form `json:"forms"`
	UserAgent   string `json:"user_agent"`
}

type Form struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	AddressLine  string `json:"address_line"`
	AddressTown  string `json:"address_town"`
	Country      string `json:"country"`
	PostCode     string `json:"post_code"`
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_number"`
}

type Msg struct {
	Type string
	Body string
}

// Do does it.
func Do(configFile string) {
	conf := loadConfig(configFile)
	notif := make(chan Msg)

	go checkSp(conf, notif)
	for _, uri := range clueURLs {
		go clueFinder(uri, conf, notif)
	}

	submitter(conf, notif)
}

func submitter(conf Config, notif chan Msg) {
	sp := ""
	clueFound := false

	for msg := range notif {
		if msg.Type == "sp" {
			sp = msg.Body
			fmt.Println("SP changed!")
		} else if msg.Type == "clue" {
			fmt.Println("Clue found!")
			clueFound = true
		}

		if sp != "" && clueFound {
			fmt.Printf("Delaying form submission by %d seconds...\n",
				conf.SubmitDelay)
			time.Sleep(time.Duration(conf.SubmitDelay) * time.Second)
			submitForms(sp, conf)
			return
		}
	}
}

func checkSp(conf Config, notif chan Msg) {
	changed := false
	for !changed {
		sp := find.GetSp()
		if sp != "" && sp != conf.CurrentSp {
			notif <- Msg{Type: "sp", Body: sp}
			changed = true
		}

		if !changed {
			fmt.Println("SP unchanged, retrying...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func clueFinder(uri string, conf Config, notif chan Msg) {
	found := false
	for !found {
		doc, err := goquery.NewDocument(uri)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find(".ProfileTweet").Each(func(i int, s *goquery.Selection) {
			text := s.Find(".ProfileTweet-text").First().Text()
			text = strings.TrimSpace(text)
			hasClue := strings.Contains(text, "Clue")
			hasHashTag := strings.Contains(text, "#20YearsOfCharacters")
			if hasClue && hasHashTag {
				timeStr, _ := s.Find(".ProfileTweet-timestamp").First().
					Attr("title")

				datetime := parseTwitterTimestamp(timeStr)
				midnight := now.BeginningOfDay()
				if datetime.After(midnight) || datetime.Equal(midnight) {
					notif <- Msg{Type: "clue"}
					found = true
				}
			}
		})
		if !found {
			fmt.Println("Clue not found, retrying...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func parseTwitterTimestamp(timeStr string) time.Time {
	layout := "15:04 - 2 Jan 2006"
	containsAm := strings.Contains(timeStr, "am")
	containsPm := strings.Contains(timeStr, "pm")
	if containsAm || containsPm {
		layout = "3:04 pm - 2 Jan 2006"
	}

	datetime, _ := time.Parse(layout, timeStr)
	return datetime
}

func submitForms(sp string, conf Config) {
	var postURL string
	if os.Getenv("TEST") == "" {
		postURL = find.GetSecretURL(sp)
	} else {
		postURL = testURL
	}

	for _, form := range conf.Forms {
		submitForm(postURL, form, conf)
	}
}

func submitForm(postURL string, form Form, conf Config) {
	formName := form.FirstName + " " + form.LastName
	fmt.Println("Submitting " + formName)

	userAgent := "Mozilla/5.0 " +
		"(Macintosh; Intel Mac OS X 10.10; rv:34.0) " +
		"Gecko/20100101 Firefox/34.0"

	if conf.UserAgent != "" {
		userAgent = conf.UserAgent
	}

	status, body, err := httpPost(
		postURL,
		map[string]string{
			"firstName":    form.FirstName,
			"lastName":     form.LastName,
			"addressLine":  form.AddressLine,
			"addressTown":  form.AddressTown,
			"Country":      form.Country,
			"postCode":     form.PostCode,
			"emailAddress": form.EmailAddress,
			"phoneNumber":  form.PhoneNumber,
			"marketing":    "0",
			"submit1994":   "Submit",
		},
		map[string]string{
			"Host":       getHost(postURL),
			"Referer":    postURL,
			"User-Agent": userAgent,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	if status == 200 && strings.Contains(body, "you have been registered") {
		fmt.Println("Successfully submitted " + formName)
	} else {
		fmt.Println("Submitting " + formName + " failed.\n" +
			"   HTTP Status: " + strconv.Itoa(status) + "\n" +
			body)
	}
}

func httpPost(uri string, body map[string]string, headers map[string]string) (int, string, error) {
	postBody := url.Values{}
	for field, value := range body {
		postBody.Add(field, value)
	}
	postBodyStr := postBody.Encode()

	r, _ := http.NewRequest("POST", uri, bytes.NewBufferString(postBodyStr))

	for header, value := range headers {
		r.Header.Add(header, value)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(postBodyStr)))

	client := &http.Client{}
	resp, err := client.Do(r)
	respBody := ""

	if err != nil {
		return 0, "", err
	} else {
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, "", err
		}
		respBody = string(content)
	}

	return resp.StatusCode, respBody, nil
}

func getHost(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	h := strings.Split(u.Host, ":")
	return h[0]
}

func loadConfig(file string) Config {
	fp, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(bytes.NewReader(fp))
	var conf Config

	if err := decoder.Decode(&conf); err != nil {
		log.Fatal(err)
	}
	return conf
}
