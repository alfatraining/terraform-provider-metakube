package main

import (
	"context"
	"flag"
	"fmt"

	"gitlab.com/furkhat/terraform-provider-metakube/client"
)

func main() {
	token := flag.String("token", "", "MataKube Access Token")
	flag.Parse()
	if *token == "" {
		flag.Usage()
		return
	}

	client := client.NewClient(client.WithBearerToken(*token))

	dcs, err := client.Datacenters.List(context.Background())
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("found %v datacenters\n", len(dcs))
	for _, dc := range dcs {
		fmt.Printf(dc.Metadata.Name)
		if dc.Seed {
			fmt.Println("(seed)")
		} else {
			fmt.Println()
		}
	}
}
