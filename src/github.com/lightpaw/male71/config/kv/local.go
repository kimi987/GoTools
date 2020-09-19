package kv

import (
	"github.com/pkg/errors"
	"net"
	"strings"
)

func getLocalAddr() (string, error) {
	ips, err := getNetInterfaceIps()
	if err != nil {
		return "", err
	}

	return chooseBestLocalIp(ips), nil
}

func getNetInterfaceIps() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, errors.Wrap(err, "无法获得本地ip")
	}

	result := make([]string, 0, 2)

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			result = append(result, ipnet.IP.String())
		}
	}

	if len(result) == 0 {
		return nil, errors.Wrap(err, "没有从网卡上找到本地ip")
	}

	return result, nil
}

func chooseBestLocalIp(ips []string) string {

	for _, ip := range ips {
		if strings.HasPrefix(ip, "192.168") {
			return ip
		}
	}

	for _, ip := range ips {
		if strings.HasPrefix(ip, "172.16") {
			return ip
		}
	}

	for _, ip := range ips {
		if strings.HasPrefix(ip, "10.0") {
			return ip
		}
	}

	return ips[0]
}
