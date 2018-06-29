package dns

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func testDNSWorks(loc string, t IPType) error {
	if t == None {
		return fmt.Errorf("cannot connect to none ip.")
	}
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	s := &http.Server{Addr: ":9090", Handler: m}
	go s.ListenAndServe()
	defer s.Close()
	backOff := time.Millisecond * 50
	timer := time.NewTimer(time.Second * 20)
	var err error
	for {
		select {
		case <-timer.C:
			return fmt.Errorf("failed to verify connection before timeout to: %s.%s", loc, DomainName)
		default:
		}
		time.Sleep(backOff)
		switch t {
		case BOTH:
			_, err = http.Get(fmt.Sprintf("http://%s.%s:9090", loc, DomainName))
			if err != nil {
				log.Println("err: failed check dns domain: None", err)
			}
		case V6:
			_, err = http.Get(fmt.Sprintf("http://ipv6.%s.%s:9090", loc, DomainName))
			if err != nil {
				log.Println("err: failed check dns domain: V6", err)
			}
		case V4:
			_, err = http.Get(fmt.Sprintf("http://ipv4.%s.%s:9090", loc, DomainName))
			if err != nil {
				log.Println("err: failed check dns domain: V4", err)
			}
		default:
			return fmt.Errorf("Cannot handler other ip types.")
		}
		if err == nil {
			return nil
		}
		backOff *= 2
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
