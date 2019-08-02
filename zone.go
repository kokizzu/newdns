package newdns

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
)

// Zone describes a single authoritative DNS zone.
type Zone struct {
	// The FQDN of the zone e.g. "example.com.".
	Name string

	// The FQDN of the master mame server responsible for this zone. The FQDN
	// must be returned as A and AAAA record by the parent zone.
	MasterNameServer string

	// A list of FQDNs to all authoritative name servers for this zone. The
	// FQDNs must be returned as A and AAAA records by the parent zone.
	AllNameServers []string

	// The email address of the administrator e.g. "hostmaster@example.com".
	//
	// Default: "hostmaster@NAME".
	AdminEmail string

	// The refresh interval.
	//
	// Default: 6h.
	Refresh time.Duration

	// The retry interval for the zone.
	//
	// Default: 1h.
	Retry time.Duration

	// The expiration interval of the zone.
	//
	// Default: 72h.
	Expire time.Duration

	// The TTl for the SOA record.
	//
	// Default: 15m.
	SOATTL time.Duration

	// The TTl for NS records.
	//
	// Default: 48h.
	NSTTL time.Duration

	// The minimum TTL for all records.
	//
	// Default: 5min.
	MinTTL time.Duration

	// The handler that responds to requests for this zone.
	Handler func(typ Type, name string) ([]Record, error)
}

// Validate will validate the zone and ensure the documented defaults.
func (z *Zone) Validate() error {
	// check name
	if !dns.IsFqdn(z.Name) {
		return fmt.Errorf("name not fully qualified")
	}

	// check master name server
	if !dns.IsFqdn(z.MasterNameServer) {
		return fmt.Errorf("master server not full qualified")
	}

	// check name servers
	for _, ns := range z.AllNameServers {
		if !dns.IsFqdn(ns) {
			return fmt.Errorf("additional name server not fully qualified")
		}
	}

	// set default admin email
	if z.AdminEmail == "" {
		z.AdminEmail = fmt.Sprintf("hostmaster@%s", z.Name)
	}

	// set default refresh
	if z.Refresh == 0 {
		z.Refresh = 6 * time.Hour
	}

	// set default retry
	if z.Refresh == 0 {
		z.Retry = time.Hour
	}

	// set default expire
	if z.Expire == 0 {
		z.Expire = 72 * time.Hour
	}

	// set default SOA TTL
	if z.SOATTL == 0 {
		z.SOATTL = 15 * time.Minute
	}

	// set default NS TTL
	if z.NSTTL == 0 {
		z.NSTTL = 48 * time.Hour
	}

	// set default min TTL
	if z.MinTTL == 0 {
		z.MinTTL = 5 * time.Minute
	}

	// check retry
	if z.Retry >= z.Refresh {
		return fmt.Errorf("retry must be less than refresh")
	}

	// check expire
	if z.Expire < z.Refresh+z.Retry {
		return fmt.Errorf("expire must be bigger than the sum of refresh and retry")
	}

	return nil
}