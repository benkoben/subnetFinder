package main

// TODO:
// - Rewrite the main module so that it can be used as an API

import (
	"fmt"
	"mymodule/subnetcalc"
    "encoding/json"
    "log"
    "io"
    "os"
    "flag"
    "strings"
	"strconv"
)

var (
    addressSpaces []subnetcalc.AddressSpace // collection that holds all addressPrefixes present in virtual network
    subnets []subnetcalc.AddressSpace       // collection that holds subnets that already exist in virtual network
    data *VirtualNetwork                    //object that holds input data (i.e. the virtual network)
    calculatedSubnets []Subnet
    out Output
    vnet = flag.String("vnet", "", `
      []object{
          AddressSpace []object{AddressPrefixes []string},
          Subnets []string
      }
      example:
      VNET=$(az network vnet show -n hub-vnet-weeu-dev-001 -g connectivity-rg-weeu-dev-001 -o json)
      ./subnetCalc -vnet $VNET
    `)
    desiredSubnets   = flag.String(
                        "new-subnets",
                        "",
                        `
    []map[string]int

    example:
        SUBNETS='[{"aks":24}, {"dbxPriv": 28}, {"dbsPub": 28}]'
        ./subnetCalc -new-subnets=$SUBNETS

    `)
)


type Output struct {
    Parameters []Subnet `json:"parameters"`
}

type Subnet struct {
    Name string `json:"name"`
    Prefix string `json:"prefix"`
}

type SpaceCollection struct {
     AddressPrefixes []string `json:"addressPrefixes"`
}

type VirtualNetwork struct {
    AddressSpace   SpaceCollection    `json:"addressSpace"`
    Subnets         []string          `json:"subnets"`
    DesiredSubnets  []map[string]int  `json:"desiredSubnets"`
    Location        string            `json:"location"`
}

func (vnet *VirtualNetwork) unmarshalVirtualNetwork(jsonString []byte){
    jsonErr := json.Unmarshal(jsonString, &data)
    if jsonErr != nil {
        fmt.Printf("Input error %v", jsonErr)
    }
}

func splitAddressNetmask(s string) (string, int){
	cidr, _ := strconv.Atoi(strings.Split(s, "/")[1])
	address := strings.Split(s, "/")[0]

	return address, cidr
}

func main() {
    flag.Parse()

    if len(*vnet) > 0 {
        log.Println("Reading from cmdline flag")
        data.unmarshalVirtualNetwork([]byte(*vnet))
    } else {
        log.Println("Reading from stdin")
        bytes, inputErr := io.ReadAll(os.Stdin)
        if inputErr != nil {
            fmt.Printf("Input error %v", inputErr)
        }
        data.unmarshalVirtualNetwork(bytes)
    }

    if len(*desiredSubnets) > 0 {
        desiredSubnets := fmt.Sprintf("{\"desiredSubnets\": %s}", *desiredSubnets)
        data.unmarshalVirtualNetwork([]byte(desiredSubnets))
    }
    log.Printf("json addressPrefixes: %#v\n", data.AddressSpace.AddressPrefixes)
    log.Printf("json subnets: %#v\n", data.Subnets)
    log.Printf("json subnets: %#v\n", data.DesiredSubnets)

    for i := range data.AddressSpace.AddressPrefixes {
		var a subnetcalc.AddressSpace
		address, cidr := splitAddressNetmask(data.AddressSpace.AddressPrefixes[i])
		// Populate the parent addressSpace
		a.Set(address, cidr)
		// Allocate the existing subnets as child addressSpaces
		addressSpaces = append(addressSpaces, a)
    }
    // Create addressSpace objects for each subnet
    // These will be added to their corresponding parent addressSpace 
    for i := range data.Subnets {
		var s subnetcalc.AddressSpace
		address, cidr := splitAddressNetmask(data.Subnets[i])
		s.Set(address, cidr)
		subnets = append(subnets, s)
	}
	// Set all children (subnet address spaces are added to parent address space)
	// Use indexes instead in order to modify the original addressSpace slice
	// https://stackoverflow.com/questions/20185511/range-references-instead-values
	for i, a := range addressSpaces {
		for j := range subnets {
			a.SetChild(&subnets[j])
		}
        // Save addressSpace to collection
		addressSpaces[i] = a
        for z :=  range data.DesiredSubnets {
            for key, val := range data.DesiredSubnets[z] {
                subnet := addressSpaces[i].NewSubnet(val)
                prefix := fmt.Sprintf("%s/%d", subnet.IpSubnet.GetIPAddress(), subnet.IpSubnet.GetNetworkSize())  
                s := Subnet{key, prefix}
                calculatedSubnets = append(calculatedSubnets, s) 
            }
        }
	}
    output := Output{Parameters: calculatedSubnets}
    b, err := json.MarshalIndent(output, "", "  ")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Print(string(b))
}
