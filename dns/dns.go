package dns

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type dnsClient struct {
	*dnsimple.Client
	accountID string
}

var (
	client      dnsClient
	DomainName  string
	DnSimpleKey string
)

func GetDnsName(loc string) (string, string, error) {

	oauthToken := DnSimpleKey
	if os.Getenv("DNSIMPLE_OAUTH_KEY") != "" {
		oauthToken = os.Getenv("DNSIMPLE_OAUTH_KEY")
	}

	if oauthToken == "" {
		return "", "", fmt.Errorf("NEED DNSIMPLE_OAUTH_KEY env var set.")
	}
	// new client
	client = dnsClient{
		Client: dnsimple.NewClient(dnsimple.NewOauthTokenCredentials(oauthToken)),
	}

	// get the current authenticated account (if you don't know who you are)
	whoamiResponse, err := client.Identity.Whoami()
	if err != nil {
		return "", "", err
	}
	accountID := fmt.Sprintf("%d", whoamiResponse.Data.Account.ID)
	client.accountID = accountID

	t, err := whatIPsCanReach()
	if err != nil {
		if t == None {
			log.Println(err)
			return "", "", err
		}
	}

	if err := updateDNS(loc, t); err != nil {
		return "", "", err
	}

	err = testDNSWorks(loc, t)
	if err != nil {
		return "", "", err
	}

	return loc, DomainName, nil
}

func updateDNS(loc string, t IPType) error {
	ip4, ip6, err := getIP46Svc()
	if err == nil {
		if t&V6 == V6 {
			err = client.setNewRecord(ip6, loc, "AAAA")
			if err != nil {
				return err
			}
			err = client.setNewRecord(ip6, "ipv6."+loc, "AAAA")
			if err != nil {
				return err
			}
			log.Println("Made AAAA records for ", loc)
		}
		if t&V4 == V4 {
			err = client.setNewRecord(ip4, loc, "A")
			if err != nil {
				return err
			}
			err = client.setNewRecord(ip4, "ipv4."+loc, "A")
			if err != nil {
				return err
			}
			log.Println("Made A records for ", loc)
		}
		return nil
	}
	return fmt.Errorf("Nope")

	// log.Println(err)
	//
	// conn, err := net.Dial("udp", "google.com:80")
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	//
	// // conn.LocalAddr().String() returns ip_address:port
	// a := conn.LocalAddr().String()
	//
	// if strings.HasPrefix(a, "[") {
	// 	a = strings.Split(a, "]:")[0][1:]
	// 	ip := net.ParseIP(a)
	// 	if ip.IsGlobalUnicast() {
	// 		err := client.setNewRecord(ip.String(), loc, "AAAA")
	// 		if err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	}
	//
	// } else {
	// 	log.Println("IP is ipv4")
	// 	return fmt.Errorf("Cannot handle ipv4 nat is bs!")
	// }
	// return fmt.Errorf("sorry cannot auto handle dns")
}

func (client *dnsClient) removeRecord(r dnsimple.ZoneRecord) error {
	_, err := client.Zones.DeleteRecord(client.accountID, r.ZoneID, r.ID)
	if err != nil {
		return err
	}
	return nil
}

func (client *dnsClient) createNewRecord(ip string, zoneName string, loc string, t string) error {
	_, err := client.Zones.CreateRecord(client.accountID, zoneName, dnsimple.ZoneRecord{
		TTL:     60,
		Name:    loc,
		Type:    t,
		Content: ip,
	})
	if err != nil {
		return fmt.Errorf("failed to create new record: %v", err)
	}
	log.Println(loc+"."+zoneName, "has been made")
	return nil
}

func (client *dnsClient) setNewRecord(ip string, loc string, t string) error {
	loc = strings.ToLower(loc)

	zones, err := client.Zones.ListZones(client.accountID, nil)
	if err != nil {
		return fmt.Errorf("failed to list zones: %v", err)
	}
	records := []dnsimple.ZoneRecord{}
	foundZone := false
	for _, z := range zones.Data {
		if z.Name != DomainName {
			continue
		}
		foundZone = true
		zone, err := client.Zones.ListRecords(client.accountID, z.Name, &dnsimple.ZoneRecordListOptions{Type: t, ListOptions: dnsimple.ListOptions{PerPage: 10000}})
		if err != nil {
			return fmt.Errorf("failed to list records for zone: %v", err)
		}
		for _, z := range zone.Data {
			if z.Type == t && z.Name == loc {
				records = append(records, z)
			}
		}
	}
	if !foundZone {
		return fmt.Errorf("failed to check zone records for: %s", DomainName)
	}
	found := false
	foundOther := false
	for _, z := range records {
		if z.Content != ip {
			err := client.removeRecord(z)
			log.Println("Removed a record for ", loc)
			if err != nil {
				log.Printf("err: failed to remove record %v :%v\n", z.Name+DomainName, err)
			}
			foundOther = true
		} else {
			log.Println("Found a record for it", loc)
			found = true
		}
		log.Println(z.Name, z.Content, z.TTL, z.Type, z.UpdatedAt)
	}
	if !found {
		return client.createNewRecord(ip, DomainName, loc, t)
	}
	if foundOther {
		return fmt.Errorf("Found other records using your dns")
	}
	return nil
}

type ipSvc struct {
	Asn      string `json:"asn"`
	AsnName  string `json:"asn_name"`
	Asnlist  string `json:"asnlist"`
	Country  string `json:"country"`
	IP       string `json:"ip"`
	Padding  string `json:"padding"`
	Protocol string `json:"protocol"`
	Subtype  string `json:"subtype"`
	Type     string `json:"type"`
	Via      string `json:"via"`
}

func getIP46Svc() (string, string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Transport: tr}
	resp, err := c.Get("https://ipv6.lookup.test-ipv6.com/ip/")
	if err != nil {
		return "", "", err
	}
	ip6 := ipSvc{}
	err = json.NewDecoder(resp.Body).Decode(&ip6)
	if err != nil {
		return "", "", err
	}

	resp, err = c.Get("https://ipv4.lookup.test-ipv6.com/ip/")
	if err != nil {
		return "", "", err
	}
	ip4 := ipSvc{}
	err = json.NewDecoder(resp.Body).Decode(&ip4)
	if err != nil {
		return "", "", err
	}
	return ip4.IP, ip6.IP, nil
}
