package dnsstore

import (
	"maps"
	"sync"
)

type DNSStore struct {
	mu      sync.RWMutex
	records map[string]string
}

func NewStore() *DNSStore {
	return &DNSStore{
		records: make(map[string]string),
	}
}

func (s *DNSStore) Set(hostname, ipv4 string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[hostname] = ipv4
}

func (s *DNSStore) Lookup(hostname string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ip, ok := s.records[hostname]
	return ip, ok
}

func (s *DNSStore) Delete(hostname string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.records, hostname)
}

func (s *DNSStore) GetAll() map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]string, len(s.records))
	maps.Copy(out, s.records)
	return out
}

const Banner = `
                o                                            8              
                8                                            8              
.oPYo. .oPYo.  o8P .oPYo. o   o   o .oPYo. o    o       .oPYo8 odYo. .oPYo. 
8    8 .oooo8   8  8oooo8 Y. .P. .P .oooo8 8    8 ooooo 8    8 8' ´8 Yb..   
8    8 8    8   8  8.     ´b.d'b.d' 8    8 8    8       8    8 8   8   'Yb. 
ˋYooP8 ˋYooP8   8  ˋYooo'  ˋY' ˋY'  ˋYooP8 ˋYooP8       ˋYooP' 8   8 ˋYooP' 
:....8 :.....:::..::.....:::..::..:::.....::....8 :::::::.....:..::..:.....:
::ooP'.::::::::::::::::::::::::::::::::::::::ooP'.::::::::::::::::::::::::::
::...::::::::::::::::::::::::::::::::::::::::...::::::::::::::::::::::::::::
`
