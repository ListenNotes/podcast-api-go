package main

import (
	"fmt"
	"strconv"

	listennotes "github.com/ListenNotes/podcast-api-go"
)

func main() {

	client := listennotes.NewClient("")

	// the test data will return the same page each time, but this is an example of getting the next_offset out fo the resulting payload
	nextOffset := fetchAndOutputPage(client, 0)
	fetchAndOutputPage(client, nextOffset)
	fetchAndOutputPage(client, nextOffset)
}

func fetchAndOutputPage(client listennotes.HTTPClient, offset int) int {
	resp, err := client.Search(map[string]string{
		"q":      "text",
		"offset": strconv.Itoa(offset),
	})
	if err != nil {
		fmt.Printf("Search failed for offset %d: %s\n", offset, err)
		return offset
	}

	fmt.Printf("Results for offset %d\n", offset)
	fmt.Printf(" Free Quota: %d\n", resp.Stats.FreeQuota)
	fmt.Printf(" Total: %.0f\n", resp.Data["total"])
	fmt.Printf(" Count: %.0f\n", resp.Data["count"])
	for _, result := range resp.Data["results"].([]interface{}) {
		if singleMap, ok := result.(map[string]interface{}); ok {
			fmt.Printf(" - %s\n", singleMap["title_original"])
		}
	}

	nextOffset, err := strconv.Atoi(fmt.Sprintf("%.0f", resp.Data["next_offset"]))
	if err != nil {
		fmt.Printf(" Failed to parse next_offset: %s\n", err)
	}

	return nextOffset
}
