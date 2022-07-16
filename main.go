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

func splitAddressNetmask(s string) (string, int){
	cidr, _ := strconv.Atoi(strings.Split(s, "/")[1])
	address := strings.Split(s, "/")[0]

	return address, cidr
}

func main() {
	var addressSpaces []subnetcalc.AddressSpace
	var dummyVnet = vnet{
		addressSpaces: []string{
			"10.100.100.0/23", 
			"172.16.0.0/24", 
			"192.168.10.0/28",
		},
		name: "bigAssLandingZone",
		allocatedSubnets: []string{
			"10.100.100.0/24",
			"10.100.100.1/24",
			"10.100.100.2/24",
			"10.100.100.3/24",
			"192.168.10.0/30",
			"192.168.10.4/30",
		},
	}
	// Loopa igenom varenda addressSpace
	for _, as := range dummyVnet.addressSpaces {
		var a subnetcalc.AddressSpace
		address, cidr := splitAddressNetmask(as)
		// Populate the parent addressSpace
		a.Set(address, cidr)
		// Allocate the existing subnets as child addressSpaces
		for index, subnet := range as.allocatedSubnets {
			var s subnetcalc.AddressSpace
			address, cidr := splitAddressNetmask
			s.Set(address, cidr)
			s.PrintIpPool()
		}
		addressSpaces = append(addressSpaces, a)
	}

	fmt.Println(addressSpaces)
}
