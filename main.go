package main

import (
	"encoding/json"
	"fmt"
	"log"

	"EVE-GO/api"
)

//Item Registry for Item Name <-> Item Type ID translation
//Use bbolt to store the SDE?

func main() {

	client := api.NewClient()

	var data []api.MarketOrder
	err := client.Market().RegionOrders(10000043).Get(&data)

	if err != nil {
		log.Panic(err)
	}

	for _, order := range data {
		darta, _ := json.Marshal(order)
		fmt.Println(string(darta))
	}
}
