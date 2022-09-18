package subnetFinder

import (
	"fmt"
	"testing"
)

var CaseNewSubnet = []struct {
	description string
	input       []struct {
		name            string
		existingSubnets []string
		size            int
		address         string
	}
}{
	{
		description: "hello",
		input: []struct {
			name            string
			existingSubnets []string
			size            int
			address         string
		}{
			{
				name:            "parent",
				existingSubnets: []string{"10.0.0.0/24", "10.0.1.128/25"},
				size:            23,
				address:         "10.0.0.0",
			},
			{
				name:            "expectedChild",
				existingSubnets: []string{},
				size:            25,
				address:         "10.0.1.0",
			},
		},
	},
}

// Function Testing cases
var CaseChildWithinScope = []struct {
	description string
	input       []struct {
		address string
		cidr    int
	}
	expected bool
}{
	{
		description: "Should pass",
		input: []struct {
			address string
			cidr    int
		}{
			{
				address: "10.0.0.0",
				cidr:    24,
			},
			{
				address: "10.0.0.128",
				cidr:    25,
			},
		},
		expected: true,
	},
	{
		description: "Should fail",
		input: []struct {
			address string
			cidr    int
		}{
			{
				address: "10.0.0.0",
				cidr:    24,
			},
			{
				address: "10.0.1.0",
				cidr:    25,
			},
		},
		expected: false,
	},
}

// Define testing variables for different functions here
var CaseConvert32BitBinaryToDecimal = []struct {
	description string
	input       string
	expected    int64
}{
	{
		description: "10.0.0.7 - Convert a single binary string (00001010000000000000000000000111) to an int64 number (167772167)",
		input:       "00001010000000000000000000000111",
		expected:    167772167,
	},
	{
		description: "192.168.10.10 - Convert a single binary string (11000000101010000000101000001010) to an int64 number (3232238090)",
		input:       "11000000101010000000101000001010",
		expected:    3232238090,
	},
}

var CaseMaskHostsSize = []struct {
	description string
	input       int
	expected    int
}{
	{
		description: "",
		input:       24,
		expected:    256,
	},
	{
		description: "",
		input:       23,
		expected:    512,
	},
	{
		description: "",
		input:       22,
		expected:    1024,
	},
	{
		description: "",
		input:       8,
		expected:    16777216,
	},
}

var CaseConvertDecimalTo32BitBinaryString = []struct {
	description string
	input       int64
	expected    string
}{
	{
		description: "10.0.0.7 - Convert number 167772167 into 00001010000000000000000000000111",
		input:       167772167,
		expected:    "00001010000000000000000000000111",
	},
	{
		description: "192.168.10.10 - Convert number 167772167 into 00001010000000000000000000000111",
		input:       3232238090,
		expected:    "11000000101010000000101000001010",
	},
}

var CaseConvertBinaryStringToDottedDecimalIPv4 = []struct {
	description string
	input       string
	expected    string
}{
	{
		description: "10.0.0.7 - Convert binary number 11000000101010000000101000001010 into ipv4 10.0.0.7",
		input:       "00001010000000000000000000000111",
		expected:    "10.0.0.7",
	},
	{
		description: "192.168.10.10 - Convert binary number 11000000101010000000101000001010 into ipv4 192.168.10.10",
		input:       "11000000101010000000101000001010",
		expected:    "192.168.10.10",
	},
}

var CaseCalculateIPv4AddressPool = []struct {
	description string
	input       *ranges
	expected    []IpAddress
}{
	{
		description: "Create IPv4Pool for 10.0.0.0/30",
		input: &ranges{
			decimin: 167772161,
			decimax: 167772163,
			binmin:  "00001010000000000000000000000001",
			binmax:  "00001010000000000000000000000011",
			addrmin: "10.0.0.1",
			addrmax: "10.0.0.3",
		},
		expected: []IpAddress{
			IpAddress{address: "10.0.0.1", mask: 32, available: true},
			IpAddress{address: "10.0.0.2", mask: 32, available: true},
			IpAddress{address: "10.0.0.3", mask: 32, available: true},
		},
	},
}

var Cases = []struct {
	input       VirtualNetwork
	expected    Output
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

func TestConvert32BitBinaryToDecimal(t *testing.T) {
	for _, i := range CaseConvert32BitBinaryToDecimal {
		result := Convert32BitBinaryToDecimal(i.input)
		if result != i.expected {
			t.Errorf("Convert32BitBinaryToDecimal(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("Convert32BitBinaryToDecimal test OK!")
		}
	}
}

func TestConvertDecimalTo32BitBinaryString(t *testing.T) {
	for _, i := range CaseConvertDecimalTo32BitBinaryString {
		result := ConvertDecimalTo32BitBinaryString(i.input)
		if result != i.expected {
			t.Errorf("Convert32BitBinaryToDecimal(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("Convert32BitBinaryToDecimal test OK!")
		}
	}
}

func TestConvertBinaryStringToDottedDecimalIPv4(t *testing.T) {
	for _, i := range CaseConvertBinaryStringToDottedDecimalIPv4 {
		result := ConvertBinaryStringToDottedDecimalIPv4(i.input)
		if result != i.expected {
			t.Errorf("ConvertBinaryStringToDottedDecimalIPv4(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("ConvertBinaryStringToDottedDecimalIPv4 test OK!")
		}
	}
}

func TestCalculateIPv4AddressPool(t *testing.T) {
	var errMsg string
	for _, i := range CaseCalculateIPv4AddressPool {
		result := CalculateIPv4AddressPool(i.input)
		errMsg = fmt.Sprintf("CalculateIPv4AddressPool(%v) = %v; want %v", i.input, result, i.expected)
		if result == nil {
			t.Errorf(errMsg)
		} else {
			for j := range result {
				if result[j] != i.expected[j] {
					t.Errorf(errMsg)
				}
			}
			t.Logf("CalculateIPv4AddressPool test OK!")
		}
	}
}

func TestMaskHostsSize(t *testing.T) {
	for _, i := range CaseMaskHostsSize {
		result := MaskHostsSize(i.input)
		if result != i.expected {
			t.Errorf("MaskHostsSize(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("MaskHostsSize test OK!")
		}
	}
}

func TestChildWithinScoppe(t *testing.T) {
	for _, i := range CaseChildWithinScope {
		var p AddressSpace
		var c AddressSpace
		p.Set(i.input[0].address, i.input[0].cidr)
		c.Set(i.input[1].address, i.input[1].cidr)

		result := ChildWithinScope(&p, &c)
		if result != i.expected {
			t.Errorf("ChildWithinScope(%v/%v, %v/%v) = %v; want %v", i.input[0].address, i.input[0].cidr, i.input[1].address, i.input[1].cidr, result, i.expected)
		} else {
			t.Logf("ChildWithinScope test OK!")
		}
	}
}
