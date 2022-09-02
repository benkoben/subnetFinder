package main

import (
	"testing"
)


var Cases = []struct {
	input    VirtualNetwork
	expected Output
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
	},
}

func TestCalculateSubnets(t *testing.T) {
	for caseIndex := range Cases {
        var result Output
		result = Cases[caseIndex].input.calculateSubnets()
        for i := range result.Parameters {
            found := false
            for j:=0;j<len(result.Parameters) && found == false ;j++ {

                expName := Cases[caseIndex].expected.Parameters[i].Name
                expPrefix := Cases[caseIndex].expected.Parameters[i].Prefix
                resultName := result.Parameters[j].Name 
                resultPrefix := result.Parameters[j].Prefix

                if expName == resultName && expPrefix == resultPrefix {
                    found = true
                }
            }
            if found == false {
			    t.Errorf("calculateSubnets(%v) = %v; want %v", Cases[caseIndex].input, result.Parameters, Cases[caseIndex].expected)
            } else  {
			    t.Logf("%v calculateSubnets test OK!",  Cases[caseIndex].expected.Parameters[i].Name)
            }
        }
	}
}
