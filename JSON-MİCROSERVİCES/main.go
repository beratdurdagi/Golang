package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/karalakrepp/Golang/BasicMicroservices/client"
)

func main() {

	listenAddr := flag.String("listenAddr", ":3000", "The Microservice's server is running")
	svc := NewLogging(NewMetricService(&priceFetcher{}))

	flag.Parse()
	api := NewJSONAPIServer(svc, *listenAddr)

	client := client.NewClient("http://localhost:3000")

	price, err := client.FetchPrice(context.Background(), "ETH")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", price)
	api.Run()
}
