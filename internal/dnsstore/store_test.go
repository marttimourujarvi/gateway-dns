package dnsstore

import "testing"

func TestLookup(t *testing.T) {
	s := NewStore()
	s.Set("foo.example.com", "1.2.3.4")

	ip, ok := s.Lookup("foo.example.com")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if ip != "1.2.3.4" {
		t.Errorf("got %q, want %q", ip, "1.2.3.5")
	}
}

func TestDelete(t *testing.T) {
	s := NewStore()
	// first set the record
	s.Set("foo.example.com", "1.2.3.4")

	s.Delete("foo.example.com")

	_, ok := s.Lookup("foo.example.com")

	if ok {
		t.Fatal("record should be missing")
	}

	records := s.GetAll()

	if len(records) > 0 {
		t.Fatal("there should be no records")
	}
}

func TestGetAll(t *testing.T) {
	s := NewStore()
	want := map[string]string{
		"foo.example.com": "1.2.3.4",
		"bar.example.com": "5.6.7.8",
	}
	for hostname, ipv4 := range want {
		s.Set(hostname, ipv4)
	}

	records := s.GetAll()

	if len(records) != len(want) {
		t.Fatalf("got %d records, want %d", len(records), len(want))
	}
	for hostname, ipv4 := range want {
		got, ok := records[hostname]
		if !ok {
			t.Errorf("missing record for %q", hostname)
			continue
		}
		if got != ipv4 {
			t.Errorf("for %q got %q, want %q", hostname, got, ipv4)
		}
	}
}
