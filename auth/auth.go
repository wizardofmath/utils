package auth

import (
	"github.com/hibooboo2/utils/auth/github"
	"github.com/hibooboo2/utils/dns"
)

func AuthedDNS() (string, error) {
	u, err := github.GetUser()
	if err != nil {
		return "", err
	}
	return dns.GetDnsName(u)
}
