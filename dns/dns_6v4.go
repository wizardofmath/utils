package dns

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func testDNSWorks(loc string, t IPType) error {
	s := &http.Server{Addr: ":9090", Handler: nil}
	go s.ListenAndServe()
	time.Sleep(time.Second)
	if t > None {
		_, err := http.Get(fmt.Sprintf("http://%s.jhrb.us:9090", loc))
		if err != nil {
			return err
		}
	}

	if t&V6 == V6 {
		_, err := http.Get(fmt.Sprintf("http://ipv6.%s.jhrb.us:9090", loc))
		if err != nil {
			return err
		}
	}
	if t&V4 == V4 {
		_, err := http.Get(fmt.Sprintf("http://ipv4.%s.jhrb.us:9090", loc))
		if err != nil {
			return err
		}
	}

	return nil
}

type IPType int

const (
	None IPType = iota
	V4
	V6
	BOTH = V4 | V6
)

func whatIPsCanReach() (IPType, error) {
	s := &http.Server{Addr: ":9090", Handler: nil}
	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	ip4, ip6, err := getIP46Svc()
	if err != nil {
		return None, err
	}
	time.Sleep(time.Second)
	_, errIP6 := http.Get(fmt.Sprintf("http://[%s]:9090", ip6))
	ip := None
	if errIP6 == nil {
		ip |= V6
	}
	_, errIP4 := http.Get(fmt.Sprintf("http://%s:9090", ip4))
	if errIP4 == nil {
		ip |= V4
	}
	switch {
	case errIP4 == nil && errIP6 == nil:
		return ip, nil
	case errIP4 != nil && errIP6 != nil:
		return ip, fmt.Errorf("failed ipv6 and ipv4 connections")
	case errIP4 != nil && errIP6 == nil:
		return ip, fmt.Errorf("failed ipv4 connection")
	case errIP6 != nil && errIP4 == nil:
		return ip, fmt.Errorf("failed ipv6 connection")
	default:
		return ip, fmt.Errorf("illogical failure!")
	}
}
