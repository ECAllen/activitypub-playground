package main

import (
	"bytes"
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

	peopleson, _ := json.MarshalIndent(people, "", "    ")
	fmt.Println("Body:", string(peopleson))

	type Actorid struct {
		Id string `json:"id"`
	}

	for _, person := range people {
		fmt.Printf("Body: %+v\n", person.Inbox.GetID())

		/*
		   NOTE
		   {
		     "@context": "https://www.w3.org/ns/activitystreams",
		     "type": "Note",
		     "name": "A Word of Warning",
		     "content": "Looks like it is going to rain today. Bring an umbrella!"
		   }
		*/

		var note pub.Note
		note.Type = pub.NoteType
		name := pub.NaturalLanguageValues{{Ref: pub.NilLangRef, Value: "Subject of the note."}}
		note.Name = name
		content := pub.NaturalLanguageValues{{Ref: pub.NilLangRef, Value: "Body of the note."}}
		note.Content = content
		fmt.Printf("%v\n", note)

		inbox := fmt.Sprintf("%s", person.Inbox.GetID())

		noteson, err := json.Marshal(note)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Post(inbox, "Content-Type of application/ld+json", bytes.NewBuffer(noteson))
		if err != nil {
			log.Fatalln(err)
		}

		var result Actorid
		json.NewDecoder(resp.Body).Decode(&result)
		resultson, _ := json.MarshalIndent(result, "", "    ")
		fmt.Println("Body:", string(resultson))
	}
}
