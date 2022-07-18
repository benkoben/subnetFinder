package subnetcalc

import (
	"testing"
	"fmt"
)

// Define testing variables for different functions here
var CaseConvert32BitBinaryToDecimal = []struct{
	description string
	input       string
	expected    int64
}{
	{
		description: "10.0.0.7 - Convert a single binary string (00001010000000000000000000000111) to an int64 number (167772167)",
		input: "00001010000000000000000000000111",
		expected: 167772167,
	},
	{
		description: "192.168.10.10 - Convert a single binary string (11000000101010000000101000001010) to an int64 number (3232238090)",
		input: "11000000101010000000101000001010",
		expected: 3232238090,
	},
}

var CaseConvertDecimalTo32BitBinaryString = []struct{
	description string
	input       int64
	expected    string
}{
	{	
		description: "10.0.0.7 - Convert number 167772167 into 00001010000000000000000000000111",
		input: 167772167,
		expected: "00001010000000000000000000000111",
	},
	{	
		description: "192.168.10.10 - Convert number 167772167 into 00001010000000000000000000000111",
		input: 3232238090,
		expected: "11000000101010000000101000001010",
	},
}

var CaseConvertBinaryStringToDottedDecimalIPv4 = []struct{
	description string
	input       string
	expected    string
}{
	{	
		description: "10.0.0.7 - Convert binary number 11000000101010000000101000001010 into ipv4 10.0.0.7",
		input: "00001010000000000000000000000111",
		expected: "10.0.0.7",
	},
	{	
		description: "192.168.10.10 - Convert binary number 11000000101010000000101000001010 into ipv4 192.168.10.10",
		input: "11000000101010000000101000001010",
		expected: "192.168.10.10",
	},
}

var CaseCalculateIPv4AddressPool = []struct{
	description string
	input       *ranges
	expected    []IpAddress
}{
	{	
		description: "Create IPv4Pool for 10.0.0.0/30",
		input: &ranges{
			decimin: 167772161,
			decimax: 167772163,
			binmin: "00001010000000000000000000000001",
			binmax: "00001010000000000000000000000011",
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
// Testing functions start here
// --

func TestConvert32BitBinaryToDecimal(t *testing.T){
	for _, i := range CaseConvert32BitBinaryToDecimal {
		result := Convert32BitBinaryToDecimal(i.input)
		if result != i.expected {
			t.Errorf("Convert32BitBinaryToDecimal(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("Convert32BitBinaryToDecimal test OK!")
		}
	}
}

func TestConvertDecimalTo32BitBinaryString(t *testing.T){
	for _, i := range CaseConvertDecimalTo32BitBinaryString {
		result := ConvertDecimalTo32BitBinaryString(i.input)
		if result != i.expected {
			t.Errorf("Convert32BitBinaryToDecimal(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("Convert32BitBinaryToDecimal test OK!")
		}
	}
}

func TestConvertBinaryStringToDottedDecimalIPv4(t *testing.T){
	for _, i := range CaseConvertBinaryStringToDottedDecimalIPv4 {
		result := ConvertBinaryStringToDottedDecimalIPv4(i.input)
		if result != i.expected {
			t.Errorf("ConvertBinaryStringToDottedDecimalIPv4(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("ConvertBinaryStringToDottedDecimalIPv4 test OK!")
		}
	}
}

func TestCalculateIPv4AddressPool(t *testing.T){
	var errMsg string
	for _, i := range CaseCalculateIPv4AddressPool {
		result := CalculateIPv4AddressPool(i.input)
		errMsg = fmt.Sprintf("CalculateIPv4AddressPool(%v) = %v; want %v", i.input, result, i.expected)
		if result == nil {
			t.Errorf(errMsg)
		} else {
			for j, _ := range result {
				if result[j] != i.expected[j] {
					t.Errorf(errMsg)
				}
			}
			t.Logf("CalculateIPv4AddressPool test OK!")
		}
	}
}