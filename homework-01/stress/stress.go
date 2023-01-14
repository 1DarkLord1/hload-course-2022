package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"main/utils"
	"io/ioutil"
	"encoding/json"
)

type CreateResponse struct {
	Longurl string
	Tinyurl string
}

func createPutRequest(longurl string) (int, string) {
	bodyJSON := []byte(fmt.Sprintf(`{"longurl": "%s"}`, longurl))

	request, err := http.NewRequest("PUT", "http://localhost:8080/create", bytes.NewBuffer(bodyJSON))

	if err != nil {
		panic(err)
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	res, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	bodyJson, _ := ioutil.ReadAll(res.Body)
	var body CreateResponse
	json.Unmarshal(bodyJson, &body)

	return res.StatusCode, body.Tinyurl
}

func urlGetBadRequest(tinyurl string) int {
	client := &http.Client{}
	res, err := client.Get("http://localhost:8080/" + tinyurl)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	return res.StatusCode
}

func urlGetGoodRequest(tinyurl string) (int,string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get("http://localhost:8080/"+tinyurl)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	loc, err := res.Location()

	if err != nil {
		panic(err)
	}

	return res.StatusCode, loc.Host
}

func main() {
	countUrls := 1000

	var longurls []string
	var tinyurls []string

	for i := 0; i < countUrls; i++ {
		strI, err := utils.UintToString(uint32(i))

		if err != nil {
			panic(err)
		}

		longurls = append(longurls, "google" + strI + ".com")
	}

	countCreate := 2000

	for i := 0; i < countCreate; i++ {
		code, tinyurl := createPutRequest("https://" + longurls[i % countUrls] + "/")

		if code != 200 {
			panic("Wrong status code " + strconv.Itoa(code))
		}

		if i < countUrls {
			tinyurls = append(tinyurls, tinyurl)
		}
	}

	fmt.Println("Create requests sent")

	for i := 0; i < countUrls; i++ {
		tinyurl := tinyurls[i] 
		code, longurl := urlGetGoodRequest(tinyurl)

		if code != 302 {
			panic("Wrong status code " + strconv.Itoa(code))
		}

		if longurl != longurls[i % countUrls] {
			print(longurls[i])
			panic("Longurl and tinyurl mismatch " + longurl + " " + longurls[i])
		}
	}

	fmt.Println("Url get 302 requests sent")

	for i := 0; i < countUrls; i++ {
		oldUrlUint, err := utils.StringToUint(tinyurls[i])

		if err != nil {
			panic(err)
		}

		badTinyurl, err := utils.UintToString(oldUrlUint + uint32(countUrls))

		if err != nil {
			panic(err)
		}

		code := urlGetBadRequest(badTinyurl)

		if code != 404 {
			panic("Wrong status code " + strconv.Itoa(code))
		}

	}

	fmt.Println("Url get 404 requests sent")
}
