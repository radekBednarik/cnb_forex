package api

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetDailyData(date string) (string, error) {
	endpoint := fmt.Sprintf("https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt?date=%s", date)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("Error when creating new request for endpoint %s\n", endpoint)
		return "", err
	}

	req.Header.Add("Accept", "text/html")
	req.Header.Add("Accept-Encoding", "gzip,deflate,br")
	req.Header.Add("Accept-Language", "cs,sk;q=0.8,en-US;q=0.5;en;q=0.3")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Host", "www.cnb.cz")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:122.0) Gecko/20100101 Firefox/122.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("HTTP method 'GET' failed with error:\n%v\n", err)
		return "", err
	}

	defer res.Body.Close()

	// for now, just handle 200 OK
	if res.StatusCode == 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Failed to read response body of the endpoint: %s\n", endpoint)
			return "", err
		}

		return string(body), nil
	}
	eMsg := fmt.Sprintf("HTTP method 'GET' to endpoint %s returned status code:\n%d\n", endpoint, res.StatusCode)
	log.Print(eMsg)
	return "", errors.New(eMsg)
}
