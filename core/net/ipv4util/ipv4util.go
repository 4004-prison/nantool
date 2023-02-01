package ipv4util

import (
	"errors"
	"fmt"
	"nantool/core/convert"
	"net"
	"strconv"
	"strings"
)

const (
	IP_SPLIT          = "." // IP_SPLIT use split the ip v4 address.
	IP_SPLIT_TAG      = "-" // IP_SPLIT_TAG use to split the ip v4 interval
	IP_MASK_SPLIT_TAG = "/" // IP_MASK_SPLIT_TAG use to split subnet ip and mask bit
	MAX_MASK_BIT      = 32  // MAX_MASK_BIT the maximum value of the mask bit
)

// IsIPv4 Check whether the string is ipv4.
func IsIPv4(ip string) bool {
	return net.ParseIP(ip) != nil
}

// List get the ipv4 address set.
// ipRange is the ip segment. The format is 'x.x.x.x-x.x.x.x' or 'x.x.x.x/maskbit'
// isAll use 'x.x.x.x/maskbit'.It's true if you want to get all ip in network segment.Otherwise, it's false
func List(ipRange string, isAll bool) ([]string, error) {
	if strings.Contains(ipRange, IP_SPLIT_TAG) {
		ipArray := strings.Split(ipRange, IP_SPLIT_TAG)
		return listRange(ipArray[0], ipArray[1])
	}

	if strings.Contains(ipRange, IP_MASK_SPLIT_TAG) {
		ipArray := strings.Split(ipRange, IP_MASK_SPLIT_TAG)
		mask, err := strconv.ParseInt(ipArray[1], 10, 64)
		if err != nil {
			return nil, errors.New("'ipRange' does not conform to the format specification")
		}
		return list(ipArray[0], int(mask), isAll)
	}

	return nil, errors.New("'ipRange' does not conform to the format specification")
}

// GetMaskByMaskBit Get the mask by the mask bit.
func GetMaskByMaskBit(maskBit int) string {
	return Int64ToIPv4(GetMaskI64ByMaskBit(maskBit))
}

// GetMaskI64ByMaskBit Gets the subnet mask in integer form by mask bit.
func GetMaskI64ByMaskBit(maskBit int) int64 {
	return 0xffffffff >> (32 - maskBit) << (32 - maskBit)
}

// Int64ToIPv4 The integer ip convert to a string
func Int64ToIPv4(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip>>24&0xff, ip>>16&0xff, ip>>8&0xff, ip&0xff)
}

// IPv4ToInt64 The ip convert to int64
func IPv4ToInt64(ip string) (int64, error) {
	if !IsIPv4(ip) {
		return 0, errors.New("IPv4 format error")
	}
	ipBytes, err := convert.StrArrToIntArr(strings.Split(ip, IP_SPLIT))
	if err != nil {
		return 0, errors.New("IPv4 format error")
	}
	return int64(ipBytes[0])<<24 + int64(ipBytes[1])<<16 + int64(ipBytes[2])<<8 + int64(ipBytes[3]), nil
}

func listRange(start string, end string) ([]string, error) {
	sIPI64, err := IPv4ToInt64(start)
	if err != nil {
		return nil, err
	}
	eIPI64, err := IPv4ToInt64(end)
	if err != nil {
		return nil, err
	}
	return listRangeI64(sIPI64, eIPI64), nil
}

func listRangeI64(start int64, end int64) []string {
	ips := make([]string, 0)
	for i := start; i <= end; i++ {
		ips = append(ips, Int64ToIPv4(i))
	}
	return ips
}

func list(ip string, maskBit int, isAll bool) ([]string, error) {
	if maskBit == MAX_MASK_BIT {
		ips := make([]string, 0)
		return append(ips, ip), nil
	}
	first, last, err := getFirstIPAndLastIpInt64(ip, maskBit)
	if err != nil {
		return nil, err
	}
	if isAll {
		return listRangeI64(first, last), nil
	}
	return listRangeI64(first+1, last-1), nil
}

func getFirstIPAndLastIpInt64(ip string, maskBit int) (int64, int64, error) {
	ipI64, err := IPv4ToInt64(ip)
	if err != nil {
		return 0, 0, err
	}
	mask := GetMaskI64ByMaskBit(maskBit)

	firstIP := ipI64 & mask
	return firstIP, firstIP + (0xffffffff ^ mask), nil
}
