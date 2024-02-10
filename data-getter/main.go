package main

import "fmt"

func main() {
	response, err := GetDailyData("08.02.2024")
	if err != nil {
		panic(1)
	}

	fmt.Println(response.Body)

	response.Body.Close()
}
