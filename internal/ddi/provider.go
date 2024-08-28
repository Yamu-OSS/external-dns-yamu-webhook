package ddi

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
)

// Provider type for interfacing with YamuDDI
type Provider struct {
	provider.BaseProvider

	client       *httpClient
	domainFilter endpoint.DomainFilter
}

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
	log.Debugf("records: retrieving records from YamuDDI")

	// TODO: 返回集群中由 exterNal-dns 管理的所有记录
	// records, err := p.client.GetHostOverrides()
	// if err != nil {
	// 	return nil, err
	// }

	// var endpoints []*endpoint.Endpoint
	// for _, record := range records {
	// 	ep := &endpoint.Endpoint{
	// 		DNSName:    JoinUnboundFQDN(record.Hostname, record.Domain),
	// 		RecordType: PruneUnboundType(record.Rr),
	// 		Targets:    endpoint.NewTargets(record.Server),
	// 	}

	// 	if !p.domainFilter.Match(ep.DNSName) {
	// 		continue
	// 	}

	// 	endpoints = append(endpoints, ep)
	// }

	log.Debugf("records: retrieved: %+v", endpoints)

	return endpoints, nil
}

// ApplyChanges applies a given set of changes in the DNS provider.
func (p *Provider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	log.Debugf("apply: changes: %+v", changes)
	// TODO: 创建更新删除，集群中由 exterNal-dns 管理的所有记录
	// for _, endpoint := range append(changes.UpdateOld, changes.Delete...) {
	// 	if err := p.client.DeleteHostOverride(endpoint); err != nil {
	// 		return err
	// 	}
	// }

	// for _, endpoint := range append(changes.Create, changes.UpdateNew...) {
	// 	if _, err := p.client.CreateHostOverride(endpoint); err != nil {
	// 		return err
	// 	}
	// }

	// p.client.ReconfigureUnbound()
	log.Debugf("apply: changes applied")
	return nil
}

// GetDomainFilter returns the domain filter for the provider.
func (p *Provider) GetDomainFilter() endpoint.DomainFilter {
	return p.domainFilter
}
