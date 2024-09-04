package ddi

import "testing"

var client *httpClient
var addRRs map[string]*DNSRecord

func init() {
	var err error
	c := &Config{
		Host:          "https://192.168.19.117",
		User:          "admin",
		Key:           "123456",
		SkipTLSVerify: true,
		View:          "default",
		DefaultTTL:    0,
	}
	client, err = newYamuDDIClient(c)
	if err != nil {
		panic(err)
	}

	addRRs = map[string]*DNSRecord{
		"testA": {
			Name:        "www",
			Rtype:       "A",
			TTL:         0,
			TTLStrategy: strategyInherit,
			Rdata:       "123.123.123.123",
			Enabled:     true,
			Source:      source,
		},
		"testAAAA": {
			Name:        "www",
			Rtype:       "AAAA",
			TTL:         30,
			TTLStrategy: strategyRewrite,
			Rdata:       "2001:db8::1",
			Enabled:     true,
			Source:      source,
		},
		"testCNAME": {
			Name:        "cname",
			Rtype:       "CNAME",
			TTL:         30,
			TTLStrategy: strategyRewrite,
			Rdata:       "abc.com",
			Enabled:     true,
			Source:      source,
		},
	}
}

func TestCreateHostOverride(t *testing.T) {
	t.Skip("need a real server to test")
	for tName, rr := range addRRs {
		err := client.CreateHostOverride("test.com", rr)
		if err != nil {
			t.Errorf("TestCreateHostOverride=%v, test=%v", err, tName)
		}
	}
}

func TestGetHostOverrides(t *testing.T) {
	t.Skip("need a real server to test")
	rrs, err := client.GetHostOverrides("test.com")
	if err != nil {
		t.Errorf("TestCreateHostOverride=%v, wantNumOfRRs!=%v", err, len(rrs))
	}
}

func TestDeleteHostOverrideBulk(t *testing.T) {
	t.Skip("need a real server to test")
	for tName, rr := range addRRs {
		err := client.DeleteHostOverrideBulk("test.com", []*DNSRecord{rr})
		if err != nil {
			t.Errorf("TestDeleteHostOverrideBulk=%v, test=%v", err, tName)
		}
	}
}
