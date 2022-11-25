package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	subnetcalc "github.com/benkoben/azsubnetcalc"
)

var (
	data  subnetcalc.VirtualNetwork
	bytes []byte
	vnet  = flag.String("vnet", "", `
	[]object{
		addressSpace []object{AddressPrefixes []string},
		Subnets []string
	}
	example:
	VNET=$(az network vnet show -n hub-vnet-weeu-dev-001 -g connectivity-rg-weeu-dev-001 -o json)
	./subnetCalc -vnet $VNET
  `)
	desiredSubnets = flag.String(
		"new-subnets",
		"",
		`
	[]map[string]int

	example:
		SUBNETS='[{"aks":24}, {"dbxPriv": 28}, {"dbsPub": 28}]'
		./subnetCalc -new-subnets=$SUBNETS
  `)
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Application returned an error: ", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	/*
		Try to read from cmdline flags first
		if no argument is given to the -vnet flag then try to read from stdin
	*/
	if len(*vnet) > 0 {
		bytes = []byte(*vnet)
	} else if len(*vnet) == 0 { // Read from stdin when no cmdline flag is present
		stdIn, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		} else if len(stdIn) == 0 {
			return errors.New("No vnet could be found in either -vnet flag nor stdin")
		}
		bytes = stdIn
	}

	err := data.UnmarshalVirtualNetwork(bytes)
	if err != nil {
		return err
	}

	if len(*desiredSubnets) > 0 {
		desiredSubnets := fmt.Sprintf("{\"desiredSubnets\": %s}", *desiredSubnets)
		data.UnmarshalVirtualNetwork([]byte(desiredSubnets))
	} else {
		return errors.New("Argument -new-subnets must not be empty")
	}

	result, err := data.CalculateSubnets()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	fmt.Print(string(b))
	return nil
}
