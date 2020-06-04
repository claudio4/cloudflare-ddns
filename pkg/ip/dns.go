package ip

import "net"

func lookupIPs(domain string) (ipv4 []net.IP, ipv6 []net.IP, err error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, nil, err
	}
	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4 = append(ipv4, ip)
		} else {
			ipv6 = append(ipv6, ip)
		}
	}
	return
}
