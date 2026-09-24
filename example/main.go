package main

import (
	"fmt"
	"log"

	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	response, err := listennotes.NewClient("").FetchPodcastRegions(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.ToJSON())
}
