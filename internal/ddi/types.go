package ddi

// Config represents the configuration for the UniFi API.
type Config struct {
	Host          string `env:"YAMU_HOST,notEmpty"`
	User          string `env:"YAMU_API_USER,notEmpty"`
	Key           string `env:"YAMU_API_KEY,notEmpty"`
	SkipTLSVerify bool   `env:"YAMU_DDI_SKIP_TLS_VERIFY" envDefault:"true"`
}

// DNSRecord represents a DNS record in the YamuDDI Unbound API.
type DNSRecord struct {
	Uuid        string `json:"uuid"`
	Enabled     string `json:"enabled"`
	Hostname    string `json:"hostname"`
	Domain      string `json:"domain"`
	Rr          string `json:"rr"`
	Server      string `json:"server"`
	Description string `json:"description,omitempty"`
	Mx          string `json:"mx,omitempty"`
	MxPrio      string `json:"mxprio,omitempty"`
}

// unboundRecordsList is the main item returned from the YamuDDI Unbound API
// since it has some decorators we just throw this struct away
type unboundRecordsList struct {
	RowCount int         `json:"rowCount"`
	Total    int         `json:"total"`
	Current  int         `json:"current"`
	Rows     []DNSRecord `json:"Rows"`
}

// Specific format for POST against the YamuDDI Unbound API
type unboundAddHostOverride struct {
	Host DNSRecord `json:"host"`
}
