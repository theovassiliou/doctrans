package aux

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
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
	ipAddress, err := ExternalIP()
	if err != nil {
		log.Info("Unable to find IP address from OS")
	}
	return ipAddress
}

// PublicIP returns the public IP address as seen from the internet, by quering a public API for it perceived IP address
func PublicIP() (string, error) {
	url := "https://api.ipify.org?format=text" // we are using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	log.Debugf("Getting IP address from  ipify ...\n")
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		return string(ip), nil
	}

	return "", err
}
