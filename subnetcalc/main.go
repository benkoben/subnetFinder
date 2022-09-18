package subnetcalc

/*
	This file is the entrypoint for the subnetcalc module.
	Json input is unmarshalled here and required subnets are calculated.
*/

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var (
	out Output // Used to hold calculated subnets and printed to stdout as json
)

type Output struct {
	Parameters []Subnet `json:"parameters"`
}

type Subnet struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

type SpaceCollection struct {
	AddressPrefixes []string `json:"addressPrefixes"`
}

type VirtualNetwork struct {
	Space          SpaceCollection  `json:"addressSpace"`
	Subnets        []string         `json:"subnets"`
	DesiredSubnets []map[string]int `json:"desiredSubnets"`
}

/*
Convert json string from the command line into a struct
*/
func (vnet *VirtualNetwork) UnmarshalVirtualNetwork(jsonString []byte) error {
	err := json.Unmarshal(jsonString, &vnet)
	if err != nil {
		return err
	}
	return nil
}

/*
Splits an subnet prefix into two parts
string{192.168.0.0}, int{24}
*/
func splitAddressNetmask(s string) (string, int) {
	cidr, _ := strconv.Atoi(strings.Split(s, "/")[1])
	address := strings.Split(s, "/")[0]
	return address, cidr
}

/*
Save the current virtual network and calculates the possible subnets
*/
func (d *VirtualNetwork) CalculateSubnets() (Output, error) {
	var addressSpaces []addressSpace   // collection that holds all addressPrefixes present in virtual network
	var subnets []addressSpace         // collection that holds subnets that already exist in virtual network
	var calculatedSubnets = []Subnet{} // collection of subnets calculated after input processing
	for i := range d.Space.AddressPrefixes {
		var a addressSpace
		address, cidr := splitAddressNetmask(d.Space.AddressPrefixes[i])
		// Populate the parent addressSpace
		a.set(address, cidr)
		// Allocate the existing subnets as child addressSpaces
		addressSpaces = append(addressSpaces, a)
	}
	// Create addressSpace objects for each subnet
	for i := range d.Subnets {
		var s addressSpace
		address, cidr := splitAddressNetmask(d.Subnets[i])
		s.set(address, cidr)
		subnets = append(subnets, s)
	}
	// set all children (subnet address spaces are added to parent address space)
	// Use indexes instead in order to modify the original addressSpace slice
	// https://stackoverflow.com/questions/20185511/range-references-instead-values
	for i, a := range addressSpaces {
		for j := range subnets {
			/*
				No need for error checking because pre-existing subnets
				are guaranteed to fit within it's parent's scope.
			*/
			a.setChild(&subnets[j])
		}

		addressSpaces[i] = a
	}
	for i := range d.DesiredSubnets {
		for key, val := range d.DesiredSubnets[i] {
			for _, a := range addressSpaces {
				subnet := a.newSubnet(val)
				if subnet.IpSubnet != nil {
					prefix := fmt.Sprintf("%s/%d", subnet.IpSubnet.GetIPAddress(), subnet.IpSubnet.GetNetworkSize())
					s := Subnet{key, prefix}
					calculatedSubnets = append(calculatedSubnets, s)
					break
				}
			}
		}
	}

	output := Output{Parameters: calculatedSubnets}
	return output, nil
}
