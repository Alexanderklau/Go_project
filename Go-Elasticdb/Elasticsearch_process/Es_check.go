package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

func main() {
	client, err := elastic.NewClient(elastic.SetURL("http://10.0.6.247:9200"))
	if err != nil {
		fmt.Println(err)
	}
	info, code, err := client.Ping("http://10.0.6.247:9200").Do(context.Background())
	fmt.Println(info)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	clients, err := elastic.NewClusterHealthService()

}
