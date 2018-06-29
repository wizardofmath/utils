package auth

import (
	"fmt"

	"github.com/hibooboo2/utils/auth/github"
	"github.com/hibooboo2/utils/dns"
)

func AuthedDNS(user string) (string, string, error) {
	if dns.IsUser(user) {
		return dns.GetDnsName(user)
	}
	if !github.IsUser(user) {
		return "", "", fmt.Errorf("failed to auth with github as %v", user)
	}
	return dns.GetDnsName(user)
}
