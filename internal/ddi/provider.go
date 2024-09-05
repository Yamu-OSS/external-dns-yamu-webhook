package ddi

import (
	"context"
	"fmt"

	"github.com/Yamu-OSS/external-dns-yamu-webhook/pkg/arrays"
	"github.com/Yamu-OSS/external-dns-yamu-webhook/pkg/domain"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
)

// Provider type for interfacing with YamuDDI
type Provider struct {
	provider.BaseProvider

	client          *httpClient
	domainFilter    endpoint.DomainFilter
	domainFilterDDI []string
}

var (
	source          = "external-dns-yamu"
	strategyInherit = "inherit"
	strategyRewrite = "rewrite"
	supportTypes    = []string{"A", "AAAA", "CNAME"}
)

// NewYamuDDIProvider initializes a new DNSProvider.
func NewYamuDDIProvider(domainFilter endpoint.DomainFilter, config *Config) (provider.Provider, error) {
	c, err := newYamuDDIClient(config)

	if err != nil {
		return nil, fmt.Errorf("provider: failed to create the YamuDDI client: %w", err)
	}

	p := &Provider{
		client:       c,
		domainFilter: domainFilter,
	}

	return p, nil
}

// Records returns the list of HostOverride records in YamuDDI Unbound.
func (p *Provider) Records(ctx context.Context) (endpoints []*endpoint.Endpoint, err error) {
	log.Debugf("records: retrieving: %+v", endpoints)

	p.setDDIDomainFilter()

	endpoints = make([]*endpoint.Endpoint, 0)
	for _, zone := range p.domainFilterDDI {
		records, err := p.client.GetHostOverrides(zone)
		if err != nil {
			return nil, err
		}

		for _, record := range records {
			ep := &endpoint.Endpoint{
				DNSName:    domain.HostAddDomain(record.Name, zone),
				RecordType: record.Rtype,
				Targets:    endpoint.NewTargets(fmt.Sprintf("%v", record.Rdata)),
				RecordTTL:  endpoint.TTL(record.TTL),
			}
			endpoints = append(endpoints, ep)
		}
	}
	log.Infof("records: retrieved records from YamuDDI")

	return endpoints, nil
}

// ApplyChanges applies a given set of changes in the DNS provider.
func (p *Provider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	log.Debugf("apply: changes: %+v", changes)
	p.setDDIDomainFilter()

	dels := append(changes.UpdateOld, changes.Delete...)
	dsD, err := p.convertDnsRecord(dels)
	if err != nil {
		return err
	}

	for zone, rrs := range dsD {
		if err := p.client.DeleteHostOverrideBulk(zone, rrs); err != nil {
			return err
		}
	}

	creates := append(changes.Create, changes.UpdateNew...)
	dsA, err := p.convertDnsRecord(creates)
	if err != nil {
		return err
	}

	for zone, rrs := range dsA {
		for _, rr := range rrs {
			if err := p.client.CreateHostOverride(zone, rr); err != nil {
				return err
			}
		}
	}
	log.Infof("apply: changes applied")
	return nil
}

// setDomainDDIFilter sets the domain filter for the provider.
func (p *Provider) setDDIDomainFilter() error {
	for _, domain := range p.domainFilter.Filters {
		if p.client.ZoneExist(domain) {
			p.domainFilterDDI = append(p.domainFilterDDI, domain)
		}
	}
	return nil
}

// convertDnsRecord converts the endpoint to DNSRecord.
func (p *Provider) convertDnsRecord(req []*endpoint.Endpoint) (map[string][]*DNSRecord, error) {
	rd := make(map[string][]*DNSRecord, 0)
	for _, ep := range req {
		if !arrays.Contains(supportTypes, ep.RecordType) {
			log.Infof("RecordType %s is not supported", ep.RecordType)
			continue
		}
		pre, suff := domain.SplitSuffixToDomain(ep.DNSName, p.domainFilterDDI)
		if suff == "" {
			log.Infof("Does not match zone: %v", ep.DNSName)
			continue
		}

		if _, ok := rd[suff]; !ok {
			rd[suff] = make([]*DNSRecord, 0)
		}

		dnsr := &DNSRecord{
			Name:        pre,
			Rtype:       ep.RecordType,
			TTL:         uint32(ep.RecordTTL),
			TTLStrategy: strategyRewrite,
			Rdata:       ep.Targets[0],

			Enabled: true,
			Source:  source,
		}
		if p.client.Config.DefaultTTL == 0 && ep.RecordTTL == 0 {
			dnsr.TTLStrategy = strategyInherit
		}
		rd[suff] = append(rd[suff], dnsr)
	}

	return rd, nil
}
