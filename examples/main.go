package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/benkoben/subnetFinder@rework-module"
)

var (
	data subnetFinder.VirtualNetwork
)

func main() {
	flag.Parse()
	// when reading vnet details from cmdline flag
	if len(*vnet) > 0 {
		data.unmarshalVirtualNetwork([]byte(*vnet))
		// when reading vnet details from stdin
	} else {
		bytes, inputErr := io.ReadAll(os.Stdin)
		if inputErr != nil {
			log.Printf("Input error %v", inputErr)
		}
		data.unmarshalVirtualNetwork(bytes)
	}

	if len(*desiredSubnets) > 0 {
		desiredSubnets := fmt.Sprintf("{\"desiredSubnets\": %s}", *desiredSubnets)
		data.unmarshalVirtualNetwork([]byte(desiredSubnets))
	}

	b, err := json.MarshalIndent(data.calculateSubnets(), "", "  ")
	if err != nil {
		log.Println(err)
	}
	fmt.Print(string(b))
}
