package subnetcalc

import (
	"fmt"
    "testing"
    "strings"
    "strconv"
)

// Helper functions

func splitAddressNetmask(s string) (string, int){
    cidr, _ := strconv.Atoi(strings.Split(s, "/")[1])
    address := strings.Split(s, "/")[0]

    return address, cidr
}

// Method testing cases



// The second struct (named "expectedChild") in the input also acts as the expected result
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

// --
// Testing methods start here
// --

func TestCaseNewSubnet(t *testing.T){
    for _, i := range CaseNewSubnet {
        var p AddressSpace // acts as the parent
        var c AddressSpace // acts as the expected result
        p.Set(i.input[0].address, i.input[0].size)
        c.Set(i.input[1].address, i.input[1].size)
        // Set the existingSubnets for the parent
        for _, x := range i.input[0].existingSubnets {
            var existingSubnet AddressSpace
            address, cidr := splitAddressNetmask(x)
            existingSubnet.Set(address, cidr)
            p.SetChild(&existingSubnet)
        }
        // Call the NewSubnet method in order to allocate a new subnet to the parent 
        newSubnet := p.NewSubnet(i.input[1].size)
        // Retrieve a string that holds the expected address/size notation
        // for both c and newSubnet so that we can perform a string comparison 
        expected := c.GetCidrNotation()
        result := newSubnet.GetCidrNotation()
        if result != expected {
			t.Errorf("NewSubnet(%v) = %v; want %v", i.input[0].size, result, expected)
        } else {
            t.Logf("NewSubnet test OK!")
        }
    }
}

// --
// Testing functions start here
// --

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
