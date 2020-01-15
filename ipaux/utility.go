package aux

import (
	"errors"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	aux "github.com/theovassiliou/dta-server/ipaux"
)

// ExternalIP determines the external IP adress in use
func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

// GetHostname retrieves the hostname as a string. Returns empty string in case of problems.
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Info("Unable to find hostname from OS")
		return ""
	}
	return hostname
}

// GetIPAdress determines the external IP adress in use, or empty string in case of problems.
func GetIPAdress() string {
	ipAddress, err := aux.ExternalIP()
	if err != nil {
		log.Info("Unable to find IP address from OS")
	}
	return ipAddress
}
