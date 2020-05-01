package main

import (
	"encoding/json"
	"fmt"
	pub "github.com/go-ap/activitypub"
	"log"
	"net/http"
)

func main() {

	url := "http://localhost:4000/actors"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("%v\n", err)
	}

	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	fmt.Println("Headers:", resp.Header)
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Not StatusOK, exiting...")
	}

	var people []pub.Person

	if err := json.NewDecoder(resp.Body).Decode(&people); err != nil {
		log.Fatal("%v\n", err)
	}

	// peopleStr, _ := json.MarshalIndent(people, "", "    ")
	// println("Body:", peopleStr)

	for _, person := range people {
		fmt.Printf("Body: %+v\n", person.Inbox.GetID())
	}

}
