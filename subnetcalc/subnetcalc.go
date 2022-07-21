package subnetcalc

 // TODO: 
 //       Start a function that calculcates new addresses according to desired state

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
	"github.com/brotherpowers/ipsubnet"
)

//
// --- Global variables
//

// 
// --- Custom errors
//

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

// Mark the IP address of a parent addressSpace as not available
func (ip *IpAddress) MarkIpv4Address() {
	ip.available = false
}

// Add an address space to another address space
// will error if the child is bigger in size than the parent
func (parent *AddressSpace) SetChild(child *AddressSpace) {
	if ChildWithinScope(parent, child) {
		// Add child to parent and mark the addresses in parent ipPool as not available
		parent.subnets = append(parent.subnets, child)
		firstIpIndex, lastIpIndex := IpPoolChunkIndexes(parent.ipPool, child.ranges.addrmin, child.ranges.addrmax)
		// fmt.Println(child.ranges.addrmin)
		for i:=firstIpIndex; i<=lastIpIndex; i++ {
			parent.ipPool.addresses[i].MarkIpv4Address()
		}
	}
}

// Find the next available indexes where a certain subnet can be allocated within the available parent.IpSubnet
func (parent *AddressSpace) NewSubnet(mask int) AddressSpace {
    var start int
    var stop int
    var chunk int
    var subnet AddressSpace
    chunk = MaskHostsSize(mask)
    stop = chunk - 1
    for i:=start;stop<=len(parent.ipPool.addresses); i+=chunk {
        if parent.ipPool.addresses[i].available == true && parent.ipPool.addresses[stop].available == true {
            subnet.Set(parent.ipPool.addresses[i].address, mask)
            parent.SetChild(&subnet)
            break
        }
        stop += chunk
    }
    return subnet
}

// Dereference all struct attributes and print them
// Used for debugging purposes
func (a *AddressSpace) PrintAddressSpace() {
	// fmt.Println("IPSubnet: ", a.ipSubnet)
	// fmt.Println("Pool length & content: ", len(a.ipPool.addresses), a.ipPool)
    fmt.Println("  Allocated subnets:")
    for _, s := range a.subnets {
        fmt.Println("    ", s.ipSubnet)
    }
	// fmt.Println("ranges: ", a.ranges)
}

func (a *AddressSpace) PrintChildren(){
	for _, s := range a.subnets {
		s.PrintAddressSpace()
	}
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

// Calculates the number of hosts available for a certain netmask
// This function is used to retrieve the end index number of a addressSpace.IpPool
// when a new subnet is fitted within it.
func MaskHostsSize(mask int) int {
    hostBits := 32 - mask
    binary := fmt.Sprintf("%s", strings.Repeat("1", hostBits))
    size, err := strconv.ParseInt(binary, 2, 64) 
    if err != nil {
       fmt.Println(err) 
    }
    // 
    return int(size) + 1
} 


// See if a child addressSpace actually is a legal fit within a parent addressSpace
func ChildWithinScope(parent *AddressSpace, child *AddressSpace) bool {
	var withinScope bool
	if child.ranges.decimin >= parent.ranges.decimin && child.ranges.decimax <= parent.ranges.decimax {
		withinScope = true
	}
	return withinScope
}

// Retrieves the start and stop for first and last ips in a larger collection of IP addresses
func IpPoolChunkIndexes(pool *IpPool, firstIp, lastIp string) (int, int) {
	var start int
	var stop int
	for i := 0; i < len(pool.addresses) && start == 0 || stop == 0; i++ {
		switch {
		case pool.addresses[i].address == firstIp:
			start = i
		case pool.addresses[i].address == lastIp:
			stop = i
		}
	}
	return start, stop
}
