package subnetcalc

 // TODO: 
 //       Start a function that marks IP adresses that are taken
 // 	  Create a method that adds child address spaces(subnets) to a parent
//		     - Error contol must ensure that the child actually fits within a parent
 //       Start a function that calculcates new addresses according to desired state
 //		  Reword the CalculateIPv4AddressPool so that is contains less loops (?)

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
	"github.com/brotherpowers/ipsubnet"
)

// 
// --- Custom errors
//
// type ChildOutsideOfScopeError struct {
	
// }

// func (c *ChildOutsideOfScopeError)

// var ChildOutsideOfScopeError = errors.New("")

// -
// --- Structs and methods are defined here
// -

type AddressSpace struct {
	ipSubnet *ipsubnet.Ip
	subnets []*AddressSpace
	ipPool *IpPool
	ranges *ranges
}

type IpPool struct {
	addresses []IpAddress
}

type IpAddress struct {
	address string
	mask int
	available bool
}

// used to keep track of the first and last address of each address space
// in different formats
type ranges struct {
	decimin int64
	decimax int64
	binmin string
	binmax string
	addrmin string
	addrmax string
}

//
// --- Methods start here
// 
func (ip *IpAddress) MarkIpv4Address() {
	ip.available = false
}

// Add an address space to another address space
// will error if the child is bigger in size than the parent
func (parent *AddressSpace) SetChild(child *AddressSpace) {
	// TODO: Learn about how Error interfaces work and implement proper error handling
	switch {
	case child.ranges.decimin < parent.ranges.decimin || child.ranges.decimax > parent.ranges.decimax:
		fmt.Println(fmt.Sprintf("Child range %v-v% is outside the parent scope %v-%v", child.ranges.addrmin, child.ranges.addrmax, parent.ranges.decimin,parent.ranges.decimax ))
	default:
		// Add child to parent
		parent.subnets = append(parent.subnets, child)
	}
}

// Dereference all struct attributes and print them
// Used for debugging purposes
func (a* AddressSpace) PrintAddressSpace() {
	fmt.Println("IPSubnet: ", a.ipSubnet)
	fmt.Println("Pool: ", a.ipPool)
	fmt.Println("subnets: ", a.subnets)
	fmt.Println("ranges: ", a.ranges)
}

func (a *AddressSpace) Set(address string, cidr int) {
	s := ipsubnet.SubnetCalculator(address, cidr)
	var pool *IpPool
	var asr *ranges
	asr = &ranges{
		decimin: Convert32BitBinaryToDecimal(ipsubnet.SubnetCalculator(s.GetIPAddressRange()[0], s.GetNetworkSize()).GetIPAddressBinary()),
		decimax: Convert32BitBinaryToDecimal(ipsubnet.SubnetCalculator(s.GetIPAddressRange()[1], s.GetNetworkSize()).GetIPAddressBinary()),
		binmin: ipsubnet.SubnetCalculator(s.GetIPAddressRange()[0], s.GetNetworkSize()).GetIPAddressBinary(),
		binmax: ipsubnet.SubnetCalculator(s.GetIPAddressRange()[1], s.GetNetworkSize()).GetIPAddressBinary(),
		addrmin: s.GetIPAddressRange()[0],
		addrmax: s.GetIPAddressRange()[1],
	}
	pool = &IpPool{
		addresses: CalculateIPv4AddressPool(asr),
	}
	a.ranges = asr
	a.ipSubnet = s
	a.ipPool = pool
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
func CalculateIPv4AddressPool(r *ranges) []IpAddress {
	var ipPool []IpAddress
	var decimalPool []int64
	for i := r.decimin; i <= r.decimax; i++ {
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
