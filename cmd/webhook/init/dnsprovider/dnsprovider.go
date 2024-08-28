package dnsprovider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Yamu-OSS/external-dns-yamu-webhook/cmd/webhook/init/configuration"
	"github.com/Yamu-OSS/external-dns-yamu-webhook/internal/ddi"
	"github.com/caarlos0/env/v11"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/provider"

	log "github.com/sirupsen/logrus"
)

type DDIProviderFactory func(baseProvider *provider.BaseProvider, config *ddi.Config) provider.Provider

func Init(config configuration.Config) (provider.Provider, error) {
	var domainFilter endpoint.DomainFilter
	createMsg := "creating ddi provider with "

	if config.RegexDomainFilter != "" {
		createMsg += fmt.Sprintf("regexp domain filter: '%s', ", config.RegexDomainFilter)
		if config.RegexDomainExclusion != "" {
			createMsg += fmt.Sprintf("with exclusion: '%s', ", config.RegexDomainExclusion)
		}
		domainFilter = endpoint.NewRegexDomainFilter(
			regexp.MustCompile(config.RegexDomainFilter),
			regexp.MustCompile(config.RegexDomainExclusion),
		)
	} else {
		if config.DomainFilter != nil && len(config.DomainFilter) > 0 {
			createMsg += fmt.Sprintf("domain filter: '%s', ", strings.Join(config.DomainFilter, ","))
		}
		if config.ExcludeDomains != nil && len(config.ExcludeDomains) > 0 {
			createMsg += fmt.Sprintf("exclude domain filter: '%s', ", strings.Join(config.ExcludeDomains, ","))
		}
		domainFilter = endpoint.NewDomainFilterWithExclusions(config.DomainFilter, config.ExcludeDomains)
	}

	createMsg = strings.TrimSuffix(createMsg, ", ")
	if strings.HasSuffix(createMsg, "with ") {
		createMsg += "no kind of domain filters"
	}
	log.Info(createMsg)

	ddiConfig := ddi.Config{}
	if err := env.Parse(&ddiConfig); err != nil {
		return nil, fmt.Errorf("reading ddi configuration failed: %v", err)
	}

	return ddi.NewYamuDDIProvider(domainFilter, &ddiConfig)
}
