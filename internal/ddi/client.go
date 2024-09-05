package ddi

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	log "github.com/sirupsen/logrus"
)

var (
	apiDnsPrefix = "/openapi/dns"
	apiRRCreate  = "zone/auth/rr/view/%s/zone/%s"
	apiRRDel     = apiRRCreate
	apiRRGet     = "zone/auth/rr/all/view/%s/zone/%s?source=%s"
	apiZoneGet   = "zone/auth/view/%s/zone/%s"
)

// httpClient is the DNS provider client.
type httpClient struct {
	*Config
	*http.Client
	baseURL *url.URL
}

// newYamuDDIClient creates a new DNS provider client.
func newYamuDDIClient(config *Config) (*httpClient, error) {
	u, err := url.Parse(config.Host)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	u.Path = path.Join(u.Path, apiDnsPrefix)

	// Create the HTTP client
	client := &httpClient{
		Config: config,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: config.SkipTLSVerify},
			},
		},
		baseURL: u,
	}

	return client, nil
}

// doRequest makes an HTTP request to the Yamu firewall.
func (c *httpClient) doRequest(method, path string, body io.Reader) (*http.Response, error) {
	p, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(p)
	log.Debugf("doRequest: making %s request to %s", method, u)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("doRequest: response code from %s request to %s: %d", method, u, resp.StatusCode)

	if resp.StatusCode == http.StatusBadRequest {
		defer resp.Body.Close()
		var code respCode
		if err = json.NewDecoder(resp.Body).Decode(&code); err != nil {
			return nil, err
		}

		if code.RCode != 0 {
			return nil, fmt.Errorf("%s", code.Description)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("doRequest: %s request to %s was not successful: %d", method, u, resp.StatusCode)
	}

	return resp, nil
}

// GetHostOverrides retrieves the list of records from the YamuDDI API.
func (c *httpClient) GetHostOverrides(zone string) ([]*DNSRecord, error) {
	p := path.Join(c.baseURL.Path, fmt.Sprintf(apiRRGet, c.View, zone, source))
	resp, err := c.doRequest(
		http.MethodGet,
		p,
		nil,
	)
	if err != nil {
		log.Errorf("method: %s, path: %s", http.MethodGet, p)
		return nil, err
	}
	defer resp.Body.Close()

	var records respRRs
	if err = json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, err
	}

	log.Debugf("gethost: retrieved records: %+v", len(records.Data))

	return records.Data, nil
}

// CreateHostOverride creates a new DNS A or AAAA or CNAME record in the YamuDDI API.
func (c *httpClient) CreateHostOverride(zone string, rr *DNSRecord) error {
	log.Debugf("create recored. zone: %s, rr-counts: 1", zone)
	jsonBody, err := json.Marshal([]*DNSRecord{rr})
	if err != nil {
		return err
	}
	return c.createHostOverride(zone, jsonBody)
}

// createHostOverride
func (c *httpClient) createHostOverride(zone string, jsonBody []byte) error {
	p := path.Join(c.baseURL.Path, fmt.Sprintf(apiRRCreate, c.View, zone))
	_, err := c.doRequest(
		http.MethodPost,
		p,
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		log.Errorf("method: %s, path: %s, body: %s", http.MethodPost, p, string(jsonBody))
		return err
	}

	return nil
}

// DeleteHostOverrideBulk deletes DNS records from the YamuDDI API.
func (c *httpClient) DeleteHostOverrideBulk(zone string, rrs []*DNSRecord) error {
	log.Debugf("create recored. zone: %s, rr-counts: %d", zone, len(rrs))
	jsonBody, err := json.Marshal(DNSRecordsDel{
		RRs: rrs,
	})
	if err != nil {
		return err
	}

	p := path.Join(c.baseURL.Path, fmt.Sprintf(apiRRDel, c.View, zone))
	_, err = c.doRequest(
		http.MethodDelete,
		p,
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		log.Errorf("method: %s, path: %s, body: %s", http.MethodDelete, p, string(jsonBody))
		return err
	}

	return nil
}

// ZoneExist checks if a zone exists in the DDI filter list.
func (c *httpClient) ZoneExist(domain string) bool {
	p := path.Join(c.baseURL.Path, fmt.Sprintf(apiZoneGet, c.View, domain))
	resp, err := c.doRequest(
		http.MethodGet,
		p,
		nil,
	)
	if err != nil {
		log.Errorf("Failed to get zone: %s", err)
		return false
	}
	defer resp.Body.Close()

	var code respCode
	if err = json.NewDecoder(resp.Body).Decode(&code); err != nil {
		log.Errorf("Failed to get zone: %s", err)
		return false
	}

	if code.RCode != 0 {
		log.Errorf("Failed to get zone: %s", code.Description)
		return false
	}
	return true
}

// setHeaders sets the headers for the HTTP request.
func (c *httpClient) setHeaders(req *http.Request) {
	// Add basic auth header
	yamuDDIAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.Config.User, c.Config.Key)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", yamuDDIAuth))
	req.Header.Add("Accept", "application/json")
	if req.Method != http.MethodGet {
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
	}
	// Log the request URL
	log.Debugf("headers: Requesting %s", req.URL)
}
