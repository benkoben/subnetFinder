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
    data *VirtualNetwork                    // object that holds Unmarshalled and unmodified input data (i.e. the virtual network)
    out Output                              // Used to hold calculated subnets and printed to stdout as json 
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
    Name   string   `json:"name"`
    Prefix string   `json:"prefix"`
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
        log.Printf("Input error %v", jsonErr)
    }
}

func splitAddressNetmask(s string) (string, int){
	cidr, _ := strconv.Atoi(strings.Split(s, "/")[1])
	address := strings.Split(s, "/")[0]
	return address, cidr
}

func (d *VirtualNetwork) calculateSubnets() Output {    
    var  addressSpaces []subnetcalc.AddressSpace    // collection that holds all addressPrefixes present in virtual network
    var subnets []subnetcalc.AddressSpace           // collection that holds subnets that already exist in virtual network
    var calculatedSubnets = []Subnet{}              // collection of subnets calculated after input processing
    for i := range d.AddressSpace.AddressPrefixes {
		var a subnetcalc.AddressSpace
		address, cidr := splitAddressNetmask(d.AddressSpace.AddressPrefixes[i])
		// Populate the parent addressSpace
		a.Set(address, cidr)
		// Allocate the existing subnets as child addressSpaces
		addressSpaces = append(addressSpaces, a)
    }
    // Create addressSpace objects for each subnet
    // These will be added to their corresponding parent addressSpace 
    for i := range d.Subnets {
		var s subnetcalc.AddressSpace
		address, cidr := splitAddressNetmask(d.Subnets[i])
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
        for z :=  range d.DesiredSubnets {
            for key, val := range d.DesiredSubnets[z] {
                subnet := addressSpaces[i].NewSubnet(val)
                prefix := fmt.Sprintf("%s/%d", subnet.IpSubnet.GetIPAddress(), subnet.IpSubnet.GetNetworkSize())  
                s := Subnet{key, prefix}
                calculatedSubnets = append(calculatedSubnets, s) 
            }
        }
	}
    output := Output{Parameters: calculatedSubnets}
    return output
}

func main() {
    flag.Parse()

    if len(*vnet) > 0 { // when reading vnet details from cmdline flag
        data.unmarshalVirtualNetwork([]byte(*vnet))
    } else { // when reading vnet details from stdin
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
