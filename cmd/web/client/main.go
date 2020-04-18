package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	url := "http://localhost:4000/actors"

	/*
		fmt.Println("URL:>", url)

		var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
		req, err := http.NewRequest("GET", url)
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
	*/

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("%v\n", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Not OK, exiting...")
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("%v\n", err)
	}
	fmt.Println("response Body:", string(body))
}
