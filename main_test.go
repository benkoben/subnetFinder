package main

import (
	"testing"
)

var Cases = []struct {
	input    VirtualNetwork
	expected Output
    description string
}{
	{
		input: VirtualNetwork{
			SpaceCollection{[]string{"10.100.0.0/16"}},
			[]string{},
			[]map[string]int{
				{"aks": 24},
				{"dbxPriv": 22},
				{"dbxPub": 28},
			},
		},
		expected: Output{[]Subnet{
			Subnet{"aks", "10.100.0.0/24"},
			Subnet{"dbxPub", "10.100.1.0/28"},
			Subnet{"dbxPriv", "10.100.4.0/22"},
		},
		},
        description: "Simple green field vnet 1",
	},
	{
		input: VirtualNetwork{
			SpaceCollection{[]string{"192.168.0.0/16"}},
			[]string{},
			[]map[string]int{
				{"subnet1": 25},
				{"subnet2": 26},
				{"subnet3": 26},
				{"subnet4": 27},
			},
		},
		expected: Output{[]Subnet{
			Subnet{"subnet1", "192.168.0.0/25"},
			Subnet{"subnet2", "192.168.0.128/26"},
			Subnet{"subnet3", "192.168.0.192/26"},
			Subnet{"subnet4", "192.168.1.0/27"},
		},
		},
        description: "simple green field vnet 2",
	},
	{
		input: VirtualNetwork{
			SpaceCollection{[]string{"192.168.0.0/23", "10.90.90.0/24"}},
			[]string{"192.168.1.0/25", "192.168.1.128/25"},
			[]map[string]int{
				{"192-168-0-X_01": 25},
				{"192-168-0-X_02": 25},
				{"10-90-90-X_01": 26},
				{"10-90-90-X_02": 26},
				{"10-90-90-X_03": 26},
				{"10-90-90-X_04": 26},
			},
		},
		expected: Output{[]Subnet{
			Subnet{"192-168-0-X_01", "192.168.0.0/25"},
			Subnet{"192-168-0-X_02", "192.168.0.128/25"},
			Subnet{"10-90-90-X_01", "10.90.90.0/26"},
			Subnet{"10-90-90-X_02", "10.90.90.64/26"},
			Subnet{"10-90-90-X_03", "10.90.90.128/26"},
			Subnet{"10-90-90-X_04", "10.90.90.192/26"},
		},
		},
        description: "Existing VNET with pre-allocated subnets and multiple addressSpaces",
	},
}

func TestCalculateSubnets(t *testing.T) {
	for caseIndex := range Cases {
		var result Output
		result = Cases[caseIndex].input.calculateSubnets()
		for i := range result.Parameters {
			valid := false
			for j := 0; j < len(result.Parameters) && valid == false; j++ {

				expName := Cases[caseIndex].expected.Parameters[i].Name
				expPrefix := Cases[caseIndex].expected.Parameters[i].Prefix
				resultName := result.Parameters[j].Name
				resultPrefix := result.Parameters[j].Prefix

				if expName == resultName && expPrefix == resultPrefix {
					valid = true
				}
			}
			if valid == false {
				t.Errorf("calculateSubnets(%v) = %v; want %v", Cases[caseIndex].input, result.Parameters, Cases[caseIndex].expected)
			} else {
				t.Logf("%v -- test OK!", Cases[caseIndex].description)
			}
		}
	}
}
