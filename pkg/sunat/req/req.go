package req

import (
	"io"
	"log"
	"net/http"
)

func DoRequest(client *http.Client, request *http.Request) (*http.Response, error) {
	log.Printf(">>> New Request: %s %s\n", request.Method, request.URL.String())
	for h, v := range request.Header {
		log.Printf("Header %s: %v\n", h, v)
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("error reading body: %v\n", err)
	}
	defer request.Body.Close()

	log.Println(string(body))
	res, err := client.Do(request)
	if err != nil {
		log.Printf("<<< Error in response: %v\n", err)
		return res, err
	}
	log.Printf("<<< %s\n", res.Status)

	resBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	log.Printf("%s\n", resBody)

	return res, err
}
