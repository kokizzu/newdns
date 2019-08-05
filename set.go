package newdns

import (
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
)

// Set is a set of records.
type Set struct {
	// The FQDN of the set.
	Name string

	// The type of the record.
	Type Type

	// The records in the set.
	Records []Record

	// The TTL of the set.
	//
	// Default: 5m.
	TTL time.Duration
}

// Validate will validate the set and ensure defaults.
func (s *Set) Validate() error {
	// check name
	if !IsDomain(s.Name, true) {
		return errors.Errorf("invalid name: %s", s.Name)
	}

	// check type
	if !s.Type.valid() {
		return errors.Errorf("invalid type: %d", s.Type)
	}

	// check records
	if len(s.Records) == 0 {
		return errors.Errorf("missing records")
	}

	// check CNAME records
	if s.Type == CNAME && len(s.Records) > 1 {
		return errors.Errorf("multiple CNAME records")
	}

	// validate records
	for _, record := range s.Records {
		err := record.Validate(s.Type)
		if err != nil {
			return errors.Wrap(err, "invalid record")
		}
	}

	// sort records
	sort.Slice(s.Records, func(i, j int) bool {
		// sort by data if TXT
		if s.Type == TXT {
			return s.Records[i].Data[0] < s.Records[j].Data[0]
		}

		// sort by priority and address if MX
		if s.Type == MX {
			a := fmt.Sprintf("%05d %s", s.Records[i].Priority, s.Records[i].Address)
			b := fmt.Sprintf("%05d %s", s.Records[j].Priority, s.Records[j].Address)
			return a < b
		}

		// otherwise by address
		return s.Records[i].Address < s.Records[j].Address
	})

	// set default ttl
	if s.TTL == 0 {
		s.TTL = 5 * time.Minute
	}

	return nil
}
