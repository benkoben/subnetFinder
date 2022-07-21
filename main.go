package main

import (
	"fmt"
	"mymodule/subnetcalc"
	"strings"
	"strconv"
)

type vnet struct {
	addressSpaces []string
	name string
	allocatedSubnets []string
}

type subnet struct {
	name string
	size int
}

func splitAddressNetmask(s string) (string, int){
	cidr, _ := strconv.Atoi(strings.Split(s, "/")[1])
	address := strings.Split(s, "/")[0]

	return address, cidr
}

func main() {
	var addressSpaces []subnetcalc.AddressSpace
	var subnets []subnetcalc.AddressSpace
	var dummyVnet = vnet{
		addressSpaces: []string{
			"10.100.100.0/24", 
			"172.16.0.0/24", 
			"192.168.10.0/28",
		},
		name: "bigAssLandingZone",
		allocatedSubnets: []string{
			"10.100.100.0/25",
			// "10.100.101.0/24",
			// "10.100.102.2/24",
			// "10.100.103.3/24",
			"172.16.0.0/25",
			"172.16.0.128/25",
			"192.168.10.0/30",
			"192.168.10.4/30",
		},
	}

	var newSubnets = []subnet{
		{
			name: "functionAppSubnet",
			size: 25,
		},
		{
			name: "aksSubnet",
			size: 26,
		},
		{
			name: "dataBricksPriv",
			size: 26,
		},
		{
			name: "dataBricksPub",
			size: 26,
		},
	}


	// Loopa igenom varnewenda addressSpace
	for _, as := range dummyVnet.addressSpaces {
		var a subnetcalc.AddressSpace
		
		address, cidr := splitAddressNetmask(as)
		// Populate the parent addressSpace
		a.Set(address, cidr)
		// Allocate the existing subnets as child addressSpaces
		addressSpaces = append(addressSpaces, a)
	}

	// Allocate child addressSpaces in memory for each subnet of the vnet
	for _, subnet := range dummyVnet.allocatedSubnets {
		var s subnetcalc.AddressSpace
		address, cidr := splitAddressNetmask(subnet)
		s.Set(address, cidr)
		subnets = append(subnets, s)
	}
	// fmt.Println(subnets)
	// Set all children (subnet address spaces are added to parent address space)
	// Use indexes instead in order to modify the original addressSpace slice
	// https://stackoverflow.com/questions/20185511/range-references-instead-values
	for i, a := range addressSpaces {
		for j := range subnets {
			a.SetChild(&subnets[j])
		}
		addressSpaces[i] = a
	}

    addressSpaces[0].NewSubnet(newSubnets[0].size)
    // for _, ns := range newSubnets {
    //     fmt.Println(subnetcalc.MaskHostsSize(ns.size))
    // }

	// for _, a := range addressSpaces {
	// 	a.PrintChildren()
	// }

	fmt.Println("Done :)")
}
