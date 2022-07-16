package subnetcalc

import (
	"testing"
)

// Define testing variables for different functions here
var TestsConvert32BitBinaryToDecimal = []struct{
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

var TestsConvertDecimalTo32BitBinaryString = []struct{
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

// var TestsConvertBinaryStringToDottedDecimalIPv4 = []struct{
// 	description string
// 	input       string
// 	expected    string
// }{
// 	{
// 		description: "Convert binary 00001010000000000000000000000111 to ipv4 10.0.0.7",
// 		input: "00001010000000000000000000000111",
// 		expected: "10.0.0.7",
// 	},
// 	{
// 		description: "Convert binary 00001010000000000000000000000111 to ipv4 192.168.10.10",
// 		input: "11000000101010000000101000001010",
// 		expected: "192.168.10.10",
// 	},
// }

// var TestsNewIPv4AddressPool = []struct{
// 	description string
// 	input       int64
// 	expected    []IpAddress
// }{}

// var TestsNewAddressSpace = []struct{
// 	description string
// 	input       *ipsubnet.Ip
// 	expected    AddressSpace
// }{
// 	description: "",
// 	input: []int64{174325784, 174325785, 174325786, 174325787},
// 	expected: "",
// }

// --
// Testing functions start here
// --

func TestConvert32BitBinaryToDecimal(t *testing.T){
	for _, i := range TestsConvert32BitBinaryToDecimal {
		result := Convert32BitBinaryToDecimal(i.input)
		if result != i.expected {
			t.Errorf("Convert32BitBinaryToDecimal(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("Convert32BitBinaryToDecimal test OK!")
		}
	}
}
func TestConvertDecimalTo32BitBinaryString(t *testing.T){
	for _, i := range TestsConvertDecimalTo32BitBinaryString {
		result := ConvertDecimalTo32BitBinaryString(i.input)
		if result != i.expected {
			t.Errorf("Convert32BitBinaryToDecimal(%v) = %v; want %v", i.input, result, i.expected)
		} else {
			t.Logf("Convert32BitBinaryToDecimal test OK!")
		}
	}
}