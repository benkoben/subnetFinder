package subnetcalc

 // TODO: 
 //       Start a function that marks IP adresses that are taken
 //       Start a function that calculcates new addresses according to desired state
 //		  Reword the CalculateIPv4AddressPool so that is contains less loops (?)

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
	"github.com/brotherpowers/ipsubnet"
)

// -
// --- Structs and methods are defined here
// -

type AddressSpace struct {
	ipSubnet *ipsubnet.Ip
	subnets []*ipsubnet.Ip
	ipPool IpPool
}

func (a* AddressSpace) PrintIpPool() {
	fmt.Println(a.ipPool)
}

func (a *AddressSpace) Set(address string, cidr int) {
	// Calculate multiple addresss
	var pool IpPool
	s := ipsubnet.SubnetCalculator(address, cidr)
	subCidrFirstBin := ipsubnet.SubnetCalculator(s.GetIPAddressRange()[0], s.GetNetworkSize()).GetIPAddressBinary()
	subCidrLastBin := ipsubnet.SubnetCalculator(s.GetIPAddressRange()[1], s.GetNetworkSize()).GetIPAddressBinary()
	pool.addresses = CalculateIPv4AddressPool(subCidrFirstBin, subCidrLastBin)
	a.ipSubnet = s
	a.ipPool = pool
}

type IpPool struct {
	addresses []IpAddress
}

type IpAddress struct {
	address string
	mask int
	available bool
}

func (ip *IpAddress) MarkIpv4Address() {
	ip.available = false
}

// -
// -- Functions are defined here
// -

// Convert a 32 bit binary string representing a binaryNumber into a decimal number
func Convert32BitBinaryToDecimal(thirtyTwoBits string) int64{
	var deci int64
	if i, err := strconv.ParseInt(thirtyTwoBits, 2, 64); err != nil {
		fmt.Println(err)
	} else {
		deci = i
	}
	return deci
}

// Convert 32 bit binary number string into a 4 octet IPv4 address string
func ConvertBinaryStringToDottedDecimalIPv4(nonDottedBinary string)string{
	var octets []string
	chunk := 8
	start := 0
	stop := chunk
	for i:= 0; i < len(nonDottedBinary) / chunk; i++ {
		octet, _ := strconv.ParseInt(nonDottedBinary[start:stop], 2, 64)
		octets = append(octets, strconv.FormatInt(octet, 10))
		start = stop
		stop = start + chunk
	}
	return strings.Join(octets[:], ".")
}

// binary octets for RFC1918 networks
// 10.0.0.0/8 (24 bit block)     = 00001010. 00000000.00000000.00000000
// 172.16.0.0/12 (20 bit block)  = 10110000.0001 0000.00000000.00000000
// 192.168.0.0/16 (16 bit block) = 11000000.10101000. 00000000.00000000
func ConvertDecimalTo32BitBinaryString(num int64) string {
	var binaryString string
	binaryString = fmt.Sprintf("%b", num)
	if utf8.RuneCountInString(binaryString) < 32 {
		diff := 32 - utf8.RuneCountInString(binaryString)
		padding := strings.Repeat("0", diff)
		binaryString = fmt.Sprintf("%v%v", padding, binaryString)
	}
	return binaryString
}

// Calcute the range of IP addreses between a first and last set of IP addresses expresses as binary strings
func CalculateIPv4AddressPool(firstAddress, lastAddress string) []IpAddress {
	var ipPool []IpAddress
	var decimalPool []int64
	min := Convert32BitBinaryToDecimal(firstAddress)
	max := Convert32BitBinaryToDecimal(lastAddress)
	for i := min; i <= max; i++ {
		decimalPool = append(decimalPool, i)
	}
	for _, decSnet := range decimalPool {
		binaryString := ConvertDecimalTo32BitBinaryString(decSnet)
		ipPool = append(ipPool, IpAddress{address: ConvertBinaryStringToDottedDecimalIPv4(binaryString), mask: 32, available: true})
	}
	return ipPool
}

// -
// --- Program execution defined here
// -
