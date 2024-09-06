package ddi

// Config represents the configuration for the UniFi API.
type Config struct {
	Host           string `env:"YAMU_HOST,notEmpty"`
	User           string `env:"YAMU_API_USER,notEmpty"`
	Key            string `env:"YAMU_API_KEY,notEmpty"`
	OpenAPITimeout int    `env:"YAMU_OPENAPI_TIMEOUT" envDefault:"60"`
	SkipTLSVerify  bool   `env:"YAMU_DDI_SKIP_TLS_VERIFY" envDefault:"true"`

	View       string `env:"VIEW" envDefault:"default"`
	DefaultTTL uint32 `env:"DEFAULT_TTL" envDefault:"0"`
}

// DNSRecord represents a DNS record in the YamuDDI API.
type DNSRecord struct {
	Name        string `json:"name"`
	Rtype       string `json:"qtype"`
	TTL         uint32 `json:"ttl"`
	TTLStrategy string `json:"ttlStrategy"`
	Rdata       any    `json:"rdata"`

	Enabled bool   `json:"enabled"`
	Source  string `json:"source"`
}

type respCode struct {
	RCode       int32  `json:"rcode"`
	Description string `json:"description"`
}

type respRRs struct {
	Data []*DNSRecord `json:"data"`
}

type DNSRecordsDel struct {
	RRs []*DNSRecord `json:"rrs"`
}
