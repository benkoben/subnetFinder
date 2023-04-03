package subnetcalc

import (
	"fmt"
	"testing"
)

// Helper functions

// The second struct (named "expectedChild") in the input also acts as the expected result
var CasenewSubnet = []struct {
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
	expected    []ipAddress
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
		expected: []ipAddress{
			{address: "10.0.0.1", mask: 32, available: true},
			{address: "10.0.0.2", mask: 32, available: true},
			{address: "10.0.0.3", mask: 32, available: true},
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
			[]SnetSubnet{},
			[]map[string]int{
				{"aks": 24},
				{"dbxPriv": 22},
				{"dbxPub": 28},
			},
		},
		expected: Output{[]Subnet{
			{"aks", "10.100.0.0/24"},
			{"dbxPub", "10.100.1.0/28"},
			{"dbxPriv", "10.100.4.0/22"},
		},
		},
		description: "Simple green field vnet 1",
	},
	{
		input: VirtualNetwork{
			SpaceCollection{[]string{"192.168.0.0/16"}},
			[]SnetSubnet{},
			[]map[string]int{
				{"subnet1": 25},
				{"subnet2": 26},
				{"subnet3": 26},
				{"subnet4": 27},
			},
		},
		expected: Output{[]Subnet{
			{"subnet1", "192.168.0.0/25"},
			{"subnet2", "192.168.0.128/26"},
			{"subnet3", "192.168.0.192/26"},
			{"subnet4", "192.168.1.0/27"},
		},
		},
		description: "simple green field vnet 2",
	},
	{
		input: VirtualNetwork{
			SpaceCollection{[]string{"192.168.0.0/23", "10.90.90.0/24"}},
			[]SnetSubnet{
				{"192.168.1.0/25"},
				{"192.168.1.128/25"},
			},
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
			{"192-168-0-X_01", "192.168.0.0/25"},
			{"192-168-0-X_02", "192.168.0.128/25"},
			{"10-90-90-X_01", "10.90.90.0/26"},
			{"10-90-90-X_02", "10.90.90.64/26"},
			{"10-90-90-X_03", "10.90.90.128/26"},
			{"10-90-90-X_04", "10.90.90.192/26"},
		},
		},
		description: "Existing VNET with pre-allocated subnets and multiple addressSpaces",
	},
}

// --
// Testing methods start here
// --

func TestCalculateSubnets(t *testing.T) {
	for caseIndex := range Cases {
		result, err := Cases[caseIndex].input.CalculateSubnets()
		if err != nil {
			t.Errorf("CalculateSubnets return an error: %v", err)
		}
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

func TestCasenewSubnet(t *testing.T) {
	for _, i := range CasenewSubnet {
		var p addressSpace // acts as the parent
		var c addressSpace // acts as the expected result
		p.set(i.input[0].address, i.input[0].size)
		c.set(i.input[1].address, i.input[1].size)
		// Set the existingSubnets for the parent
		for _, x := range i.input[0].existingSubnets {
			var existingSubnet addressSpace
			address, cidr := splitAddressNetmask(x)
			existingSubnet.set(address, cidr)
			p.setChild(&existingSubnet)
		}
		// Call the newSubnet method in order to allocate a new subnet to the parent
		newSubnet := p.newSubnet(i.input[1].size)

		// Retrieve a string that holds the expected address/size notation
		// for both c and newSubnet so that we can perform a string comparison
		expected := c.getCidrNotation()
		result := newSubnet.getCidrNotation()
		if result != expected {
			t.Errorf("newSubnet(%v) = %v; want %v", i.input[0].size, result, expected)
		} else {
			t.Logf("newSubnet test OK!")
		}
	}
}

// --
// Testing functions start here
// --

func TestConvert32BitBinaryToDecimal(t *testing.T) {
	for _, i := range CaseConvert32BitBinaryToDecimal {
		result := convert32BitBinaryToDecimal(i.input)
		if result != i.expected {
			t.Errorf("Convert32BitBinaryToDecimal(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("Convert32BitBinaryToDecimal test OK!")
		}
	}
}

func TestConvertDecimalTo32BitBinaryString(t *testing.T) {
	for _, i := range CaseConvertDecimalTo32BitBinaryString {
		result := convertDecimalTo32BitBinaryString(i.input)
		if result != i.expected {
			t.Errorf("Convert32BitBinaryToDecimal(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("Convert32BitBinaryToDecimal test OK!")
		}
	}
}

func TestConvertBinaryStringToDottedDecimalIPv4(t *testing.T) {
	for _, i := range CaseConvertBinaryStringToDottedDecimalIPv4 {
		result := convertBinaryStringToDottedDecimalIPv4(i.input)
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
		result := calculateIPv4AddressPool(i.input)
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
		result := maskHostsSize(i.input)
		if result != i.expected {
			t.Errorf("MaskHostsSize(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("MaskHostsSize test OK!")
		}
	}
}
