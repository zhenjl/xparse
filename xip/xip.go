// Copyright (c) 2014 Dataence, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// xip is an IP (only v4 for now) parser that expands a given IP string to all the
// IP addresses it represents.
package netx

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	stateOctet = iota
	stateRange // 2 numbers separated by dash (-), e.g., 1-10
	stateList  // Numbers separated by comma (,), e.g., 1,2,3,6
	stateCIDR
)

const (
	maxOctetValue = 0xff
)

// Parse takes a string that represents an IP address, IP range, or CIDR block and
// return a list of individual IPs. Currently only IPv4 is supported.
//
// For example:
//   10.1.1.1      -> 10.1.1.1
//   10.1.1.1,2    -> 10.1.1.1, 10.1.1.2
//   10.1.1,2.1    -> 10.1.1.1, 10.1.2.1
//   10.1.1,2.1,2  -> 10.1.1.1, 10.1.1.2 10.1.2.1, 10.1.2.2
//   10.1.1.1-2    -> 10.1.1.1, 10.1.1.2
//   10.1.1.-2     -> 10.1.1.0, 10.1.1.1, 10.1.1.2
//   10.1.1.1-10   -> 10.1.1.1, 10.1.1.2 ... 10.1.1.10
//   10.1.1.1-     -> 10.1.1.1 ... 10.1.1.254, 10.1.1.255
//   10.1.1-3.1    -> 10.1.1.1, 10.1.2.1, 10.1.3.1
//   10.1-3.1-3.1  -> 10.1.1.1, 10.1.2.1, 10.1.3.1, 10.2.1.1, 10.2.2.1, 10.2.3.1, 10.3.1.1, 10.3.2.1, 10.3.3.1
//   10.1.1        -> 10.1.1.0, 10.1.1.1 ... 10.1.1.254, 10.1.1.255
//   10.1.1-2      -> 10.1.1.0, 10.1.1.1 ... 10.1.1.255, 10.1.2.0, 10.1.2.1 ... 10.1.2.255
//   10.1-2        -> 10.1.0.0, 10.1.0,1 ... 10.2.255.254, 10..2.255.255
//   10            -> 10.0.0.0 ... 10.255.255.255
//   10.1.1.2,3,4  -> 10.1.1.1, 10.1.1.2, 10.1.1.3, 10.1.1.4
//   10.1.1,2      -> 10.1.1.0, 10.1.1.1 ... 10.1.1.255, 10.1.2.0, 10.1.2.1 ... 10.1.2.255
//   10.1.1/28     -> 10.1.1.0 ... 10.1.1.255
//   10.1.1.0/28   -> 10.1.1.0 ... 10.1.1.15
//   10.1.1.0/30   -> 10.1.1.0, 10.1.1.1, 10.1.1.2, 10.1.1.3
//   10.1.1.128/25 -> 10.1.1.128 ... 10.1.1.255
func Parse(ip string) ([]net.IP, error) {
	if strings.IndexByte(ip, ':') == -1 {
		return ParseIPv4(ip)
	}

	return nil, fmt.Errorf("parse/PasreIP: Invalid IP Address %s", ip)
}

// ParseIPv4 is called by ParseIP for IPv4 addresses. See ParseIP for more detials.
func ParseIPv4(ip string) ([]net.IP, error) {
	parts := strings.Split(ip, "/")
	if len(parts) > 2 {
		return nil, fmt.Errorf("parse/ParseIPv4: Invalid IP Address %s", ip)
	}

	ips, err := parseIPv4(parts[0])
	if err != nil {
		return nil, err
	}

	var cidr int64

	if len(parts) != 2 {
		cidr = 32
	} else {

		cidr, err = strconv.ParseInt(parts[1], 0, 8)
		if err != nil {
			return nil, err
		}

		if cidr > 32 {
			return nil, fmt.Errorf("parse/ParseIPv4: Invalid IP Address %s: Invalid CIDR notation", ip)
		}
	}

	return parseIPv4CIDR(ips, net.CIDRMask(int(cidr), 32))
}

func parseIPv4(ip string) ([]net.IP, error) {
	var (
		octets [4][]byte    // Octet 1, 2, 3, 4 of the IP address
		state  = stateOctet // Current state of the parser
		range1 = 0          // Start of the IP range, -1 means no range
		value  = 0          // Value of the current octet
		oi     = 0          // Octet index
		comma  = false      // Did we just see a comma
	)

	ip = strings.TrimSpace(ip)

	for _, b := range ip {
		//glog.Debugf("b=%q, oi=%d, state=%d, value=%d, range1=%d, comma=%t", b, oi, state, value, range1, comma)
		switch {
		case b == '.' || b == '/' || b == ',':
			// For the case of 10.1.1,2,.1, where comma is right before the dot
			if comma && b == '.' {
				continue
			}

			switch state {
			case stateOctet:
				if oi >= 3 && b != ',' {
					// Should never see dot when we are in octet 4
					return nil, fmt.Errorf("ip/parseIPv4(1): Invalid IP address %s", ip)
				}

				octets[oi] = append(octets[oi], byte(value))

			case stateRange:
				for j := range1; j <= value; j++ {
					octets[oi] = append(octets[oi], byte(j))
				}

			default:
				return nil, fmt.Errorf("ip/parseIPv4(2): Invalid IP address %s", ip)
			}

			if b == '/' {
				state = stateCIDR
			}

			value = 0
			state = stateOctet

			if b == ',' {
				comma = true
			} else {
				comma = false
				oi++
			}

		case b == '-':
			// If we hit -, that means this is a range. We save the value as the
			// start of the range, and wait for the end
			range1 = value
			value = 0
			state = stateRange
			comma = false

		case b >= '0' && b <= '9':
			value = value*10 + int(b-'0')
			if value > maxOctetValue {
				return nil, fmt.Errorf("ip/parseIPv4(4): Invalid IP address %s: octet value larger than than %x", ip, maxOctetValue)
			}

			comma = false

		default:
			return nil, fmt.Errorf("ip/parseIPv4(5): Invalid IP address %s: invalid character %b", ip, b)
		}
	}

	// End of the ip string
	switch state {
	case stateOctet, stateRange:
		if state == stateOctet {
			octets[oi] = append(octets[oi], byte(value))
		} else {
			// In case the IP is 10.1.1.1-
			if value == 0 && ip[len(ip)-1] == '-' {
				value = maxOctetValue
			}

			for j := range1; j <= value; j++ {
				octets[oi] = append(octets[oi], byte(j))
			}
		}

		for i := oi + 1; i < 4; i++ {
			for j := 0; j < 256; j++ {
				octets[i] = append(octets[i], byte(j))
			}
		}

	case stateCIDR:

	default:
		return nil, fmt.Errorf("ip/parseIPv4(6): Invalid IP address %s", ip)
	}

	var ips []net.IP

	// Create a list of IPs
	for _, o1 := range octets[0] {
		for _, o2 := range octets[1] {
			for _, o3 := range octets[2] {
				for _, o4 := range octets[3] {
					ips = append(ips, net.IPv4(o1, o2, o3, o4))
				}
			}
		}
	}

	return ips, nil
}

func parseIPv4CIDR(ips []net.IP, mask net.IPMask) ([]net.IP, error) {
	var newips []net.IP

	m := make(map[[4]byte]struct{})

	for _, ip := range ips {
		// The CIDR network
		ipnet := &net.IPNet{IP: ip.Mask(mask), Mask: mask}

		for nip := ip.Mask(mask); ipnet.Contains(nip); incIP(nip) {
			var tmp [4]byte
			copy(tmp[:], nip.To4())
			m[tmp] = struct{}{}
		}
	}

	for k := range m {
		ip2 := net.IPv4(k[0], k[1], k[2], k[3])
		newips = append(newips, ip2)
	}

	return newips, nil
}

func incIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}
