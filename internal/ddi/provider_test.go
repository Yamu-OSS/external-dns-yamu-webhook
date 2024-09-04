package ddi

import (
	"context"
	"testing"

	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
)

var p provider.Provider
var RRs *plan.Changes

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
	p, err = NewYamuDDIProvider(endpoint.DomainFilter{Filters: []string{"test.com"}}, c)
	if err != nil {
		panic(err)
	}

	RRs = &plan.Changes{
		Create: []*endpoint.Endpoint{
			{DNSName: "www.test.com", Targets: []string{"123.123.123.123"}, RecordTTL: 0, RecordType: "A"},
			{DNSName: "www.test.com", Targets: []string{"2001:db8::1"}, RecordTTL: 30, RecordType: "AAAA"},
		},
		UpdateOld: []*endpoint.Endpoint{
			{DNSName: "cname.test.com", Targets: []string{"abc.com"}, RecordTTL: 30, RecordType: "CNAME"},
		},
		UpdateNew: []*endpoint.Endpoint{
			{DNSName: "cname.test.com", Targets: []string{"abc.com"}, RecordTTL: 30, RecordType: "CNAME"},
		},
		Delete: []*endpoint.Endpoint{
			{DNSName: "www.test.com", Targets: []string{"123.123.123.123"}, RecordTTL: 0, RecordType: "A"},
			{DNSName: "www.test.com", Targets: []string{"2001:db8::1"}, RecordTTL: 30, RecordType: "AAAA"},
		}}
}

func TestApplyChanges(t *testing.T) {
	t.Skip("need a real server to test")
	err := p.ApplyChanges(context.Background(), RRs)
	if err != nil {
		t.Errorf("TestCreateHostOverride=%v, want=%v", err, nil)
	}
}

func TestRecords(t *testing.T) {
	t.Skip("need a real server to test")
	ds, err := p.Records(context.Background())
	if err != nil {
		t.Errorf("TestCreateHostOverride=%v, wantNumOfRRs!=%v", err, len(ds))
	}
}
