package ddi

import (
	"context"
	"fmt"
	"sync"

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

	client               *httpClient
	domainFilter         endpoint.DomainFilter
	domainFilterDDIRWMux sync.RWMutex
	domainFilterDDI      []string
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

type EndpointKey struct {
	DNSName    string
	RecordType string
}

// Records returns the list of HostOverride records in YamuDDI Unbound.
func (p *Provider) Records(ctx context.Context) (endpoints []*endpoint.Endpoint, err error) {

	p.setDDIDomainFilter()
	endpoints = make([]*endpoint.Endpoint, 0)
	for _, zone := range p.getDDIDomainFilter() {
		records, err := p.client.GetHostOverrides(zone)
		if err != nil {
			return nil, err
		}

		epMap := map[EndpointKey]*endpoint.Endpoint{}
		for _, record := range records {
			dnsName := domain.HostAddDomain(record.Name, zone)
			if _, ok := epMap[EndpointKey{dnsName, record.Rtype}]; !ok {
				epMap[EndpointKey{dnsName, record.Rtype}] = &endpoint.Endpoint{
					DNSName:    dnsName,
					RecordType: record.Rtype,
					RecordTTL:  endpoint.TTL(record.TTL),
				}
			}
			epMap[EndpointKey{dnsName, record.Rtype}].Targets = append(
				epMap[EndpointKey{dnsName, record.Rtype}].Targets, fmt.Sprintf("%v", record.Rdata))
		}

		for _, ep := range epMap {
			endpoints = append(endpoints, ep)
		}
	}

	log.Infof("records: retrieving: %+v", endpoints)

	return endpoints, nil
}

// ApplyChanges applies a given set of changes in the DNS provider.
func (p *Provider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	log.Infof("apply: changes: %+v", changes)
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

// setDomainDDIFilter
func (p *Provider) setDDIDomainFilter() {
	p.domainFilterDDIRWMux.Lock()
	defer p.domainFilterDDIRWMux.Unlock()

	p.domainFilterDDI = make([]string, 0)

	for _, domain := range p.domainFilter.Filters {
		if !p.client.ZoneExist(domain) {
			continue
		}

		p.domainFilterDDI = append(p.domainFilterDDI, domain)
	}
}

// getDDIDomainFilter
func (p *Provider) getDDIDomainFilter() []string {
	p.domainFilterDDIRWMux.RLock()
	defer p.domainFilterDDIRWMux.RUnlock()

	return p.domainFilterDDI
}

// convertDnsRecord converts the endpoint to DNSRecord.
func (p *Provider) convertDnsRecord(req []*endpoint.Endpoint) (map[string][]*DNSRecord, error) {
	rd := make(map[string][]*DNSRecord, 0)
	for _, ep := range req {
		if !arrays.Contains(supportTypes, ep.RecordType) {
			log.Infof("RecordType %s is not supported", ep.RecordType)
			continue
		}
		pre, suff := domain.SplitSuffixToDomain(ep.DNSName, p.getDDIDomainFilter())
		if suff == "" {
			log.Infof("Does not match zone: %v", ep.DNSName)
			continue
		}

		if _, ok := rd[suff]; !ok {
			rd[suff] = make([]*DNSRecord, 0)
		}

		for _, target := range ep.Targets {
			dnsr := &DNSRecord{
				Name:        pre,
				Rtype:       ep.RecordType,
				TTL:         uint32(ep.RecordTTL),
				TTLStrategy: strategyRewrite,
				Rdata:       target,

				Enabled: true,
				Source:  source,
			}
			if p.client.Config.DefaultTTL == 0 && ep.RecordTTL == 0 {
				dnsr.TTLStrategy = strategyInherit
			}
			if p.client.Config.DefaultTTL != 0 && ep.RecordTTL == 0 {
				// if the TTL is not set and the default TTL is not 0, use the default TTL
				dnsr.TTL = p.client.Config.DefaultTTL
			}
			rd[suff] = append(rd[suff], dnsr)
		}
	}

	return rd, nil
}

// GetDomainFilter returns the domain filter for the provider.
func (p *Provider) GetDomainFilter() endpoint.DomainFilter {
	return p.domainFilter
}
