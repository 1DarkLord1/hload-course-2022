package main

import (
	"bytes"
	"fmt"
	"net/http"
)

type CreateResponse struct {
	Longurl string
	Tinyurl string
}

func createPutRequest(longurl string) {
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

	if res.StatusCode != 200 {
		panic(res.StatusCode)
	}
}

func urlGetRequest(tinyurl string, expectedStatus int) {
	res, err := http.Get("http://localhost:8080/" + tinyurl)

	if err != nil {
		panic(err)
	}

	if res.StatusCode != expectedStatus {
		panic(res.StatusCode)
	}
}

func main() {
	longurls := []string{
		"https://google.com/", 
		"https://math-cs.spbu.ru/", 
		"https://github.com/",
		"https://go.dev/",
	}

	tinyurls := []string {
		"7000000", 
		"8000000", 
		"9000000",
		"a000000",
	}

	badTinyurls := []string {
		"((()))#",
		"%$%&*()",
		"!@#$",
		"!!!!@@@@@@",
	}

	countCreate := 10000
	countGet200 := 1000000
	countGet404 := 1000000

	for i := 0; i < countCreate; i++ {
		createPutRequest(longurls[i % len(longurls)])
	}

	fmt.Println("Create requests sent")

	for i := 0; i < countGet200; i++ {
		urlGetRequest(tinyurls[i % len(tinyurls)], 200)
	}

	fmt.Println("Url get 302 requests sent")

	for i := 0; i < countGet404; i++ {
		urlGetRequest(badTinyurls[i % len(badTinyurls)], 400)
	}

	fmt.Println("Url get 404 requests sent")
}
