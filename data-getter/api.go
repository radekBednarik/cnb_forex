package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetDailyData(date string) (string, error) {
	endpoint := fmt.Sprintf("https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt?date=%s", date)

	res, err := http.Get(endpoint)

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
