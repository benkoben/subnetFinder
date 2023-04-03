package subnetcalc

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

type addressSpace struct {
	IpSubnet *ipsubnet.Ip
	subnets  []*addressSpace
	ipPool   *ipPool
	ranges   *ranges
}

type ipPool struct {
	addresses []ipAddress
}

type ipAddress struct {
	address   string
	mask      int
	available bool
}

/*
used to keep track of the first and last address of each address space
in different formats
*/
type ranges struct {
	decimin int64
	decimax int64
	binmin  string
	binmax  string
	addrmin string
	addrmax string
}

/*
Mark the IP address of a parent addressSpace as not available
*/
func (ip *ipAddress) markIpv4Address() {
	ip.available = false
}

/*
Add an address space to another address space
*/
func (parent *addressSpace) setChild(child *addressSpace) {
	if childWithinScope(parent, child) {
		// Add child to parent and mark the addresses in parent ipPool as not available
		parent.subnets = append(parent.subnets, child)
		firstIpIndex, lastIpIndex := ipPoolChunkIndexes(parent.ipPool, child.ranges.addrmin, child.ranges.addrmax)
		// fmt.Println(child.ranges.addrmin)
		for i := firstIpIndex; i <= lastIpIndex; i++ {
			parent.ipPool.addresses[i].markIpv4Address()
		}
	}
}

/*
Returns a string representing a complete network address prefix (CIDR notation)
*/
func (a *addressSpace) getCidrNotation() string {
	return fmt.Sprintf("%v/%v", a.IpSubnet.GetIPAddress(), a.IpSubnet.GetNetworkSize())
}

/*
Creates a new addressSpace representing a subnet and checks its validition in relation of it's parent's prefix
*/
func (parent *addressSpace) newSubnet(mask int) addressSpace {
	var start int
	var stop int
	var chunk int
	var subnet addressSpace
	chunk = maskHostsSize(mask)
	stop = chunk - 1

	for i := start; stop <= len(parent.ipPool.addresses); i += chunk {
		if parent.ipPool.addresses[i].available == true && parent.ipPool.addresses[stop].available == true {
			subnet.set(parent.ipPool.addresses[i].address, mask)
			parent.setChild(&subnet)
			break
		}
		stop += chunk
	}
	return subnet
}

/*
Dereference all struct attributes and print them
Used for debugging purposes
*/
func (a *addressSpace) printaddressSpace() {
	fmt.Println("IPSubnet: ", a.IpSubnet)
	fmt.Println("Pool length & content: ", len(a.ipPool.addresses), a.ipPool)
	fmt.Println("  Allocated subnets:")
	for _, s := range a.subnets {
		fmt.Println("    ", s.IpSubnet)
	}
}

/*
Use for debugging purposes. Dereference and print all of the children currently allocated to an addressSpace
*/
func (a *addressSpace) printChildren() {
	for _, s := range a.subnets {
		s.printaddressSpace()
	}
}

/*
Initialize a new addressSpace
All keywords are set with values representing a Virtual network object
*/
func (a *addressSpace) set(address string, cidr int) {
	s := ipsubnet.SubnetCalculator(address, cidr)
	var pool *ipPool
	var asr *ranges
	asr = &ranges{
		decimin: convert32BitBinaryToDecimal(ipsubnet.SubnetCalculator(s.GetIPAddressRange()[0], s.GetNetworkSize()).GetIPAddressBinary()),
		decimax: convert32BitBinaryToDecimal(ipsubnet.SubnetCalculator(s.GetIPAddressRange()[1], s.GetNetworkSize()).GetIPAddressBinary()),
		binmin:  ipsubnet.SubnetCalculator(s.GetIPAddressRange()[0], s.GetNetworkSize()).GetIPAddressBinary(),
		binmax:  ipsubnet.SubnetCalculator(s.GetIPAddressRange()[1], s.GetNetworkSize()).GetIPAddressBinary(),
		addrmin: s.GetIPAddressRange()[0],
		addrmax: s.GetIPAddressRange()[1],
	}
	pool = &ipPool{
		addresses: calculateIPv4AddressPool(asr),
	}
	a.ranges = asr
	a.IpSubnet = s
	a.ipPool = pool
}

// -
// -- Functions are defined here
// -

// Convert a 32 bit binary string representing a binaryNumber into a decimal number
func convert32BitBinaryToDecimal(thirtyTwoBits string) int64 {
	var deci int64
	if i, err := strconv.ParseInt(thirtyTwoBits, 2, 64); err != nil {
		fmt.Println(err)
	} else {
		deci = i
	}
	return deci
}

// Convert 32 bit binary number string into a 4 octet IPv4 address string
func convertBinaryStringToDottedDecimalIPv4(nonDottedBinary string) string {
	var octets []string
	chunk := 8
	start := 0
	stop := chunk
	for i := 0; i < len(nonDottedBinary)/chunk; i++ {
		octet, _ := strconv.ParseInt(nonDottedBinary[start:stop], 2, 64)
		octets = append(octets, strconv.FormatInt(octet, 10))
		start = stop
		stop = start + chunk
	}
	return strings.Join(octets[:], ".")
}

/*
binary octets for RFC1918 networks
10.0.0.0/8 (24 bit block)     = 00001010. 00000000.00000000.00000000
172.16.0.0/12 (20 bit block)  = 10110000.0001 0000.00000000.00000000
192.168.0.0/16 (16 bit block) = 11000000.10101000. 00000000.00000000
*/
func convertDecimalTo32BitBinaryString(num int64) string {
	var binaryString string
	binaryString = fmt.Sprintf("%b", num)
	if utf8.RuneCountInString(binaryString) < 32 {
		diff := 32 - utf8.RuneCountInString(binaryString)
		padding := strings.Repeat("0", diff)
		binaryString = fmt.Sprintf("%v%v", padding, binaryString)
	}
	return binaryString
}

/*
Calcute the range of IP addreses between a first and last set of IP addresses expresses as binary strings
*/
func calculateIPv4AddressPool(r *ranges) []ipAddress {
	var ipPool []ipAddress
	var decimalPool []int64
	for i := r.decimin; i <= r.decimax; i++ {
		decimalPool = append(decimalPool, i)
	}
	for _, decSnet := range decimalPool {
		binaryString := convertDecimalTo32BitBinaryString(decSnet)
		ipPool = append(ipPool, ipAddress{address: convertBinaryStringToDottedDecimalIPv4(binaryString), mask: 32, available: true})
	}
	return ipPool
}

/*
Calculates the number of hosts available for a certain netmask
This function is used to retrieve the end index number of a addressSpace.ipPool
when a new subnet is fitted within it.
*/
func maskHostsSize(mask int) int {
	hostBits := 32 - mask
	binary := fmt.Sprintf("%s", strings.Repeat("1", hostBits))
	size, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		fmt.Println(err)
	}
	//
	return int(size) + 1
}

/*
See if a child addressSpace actually is a legal fit within a parent addressSpace
*/
// func childWithinScope(parent *addressSpace, child *addressSpace) error {
// 	if child.ranges.decimin <= parent.ranges.decimin && child.ranges.decimax >= parent.ranges.decimax {
// return errors.New(fmt.Sprintf("All subnets dont fit: subnet with size %d does not fit with current specification",
// 	child.IpSubnet.GetNetworkSize(),
// ))
// 	}
// 	return nil
// }

func childWithinScope(parent, child *addressSpace) bool {
	var withinScope bool
	if child.ranges.decimin >= parent.ranges.decimin && child.ranges.decimax <= parent.ranges.decimax {
		withinScope = true
	}
	return withinScope
}

/*
Retrieves the start and stop for first and last ips in a larger collection of IP addresses
*/
func ipPoolChunkIndexes(pool *ipPool, firstIp, lastIp string) (int, int) {
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
