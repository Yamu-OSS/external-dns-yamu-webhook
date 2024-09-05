package domain

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/idna"
)

type Domain string

func NewDomain(s string) Domain {
	return Domain(s)
}

func (dm Domain) ToFQDN() Domain {
	str := dm.ToString()
	if !strings.HasSuffix(str, ".") {
		str += "."
	}

	return Domain(str)
}

func (dm Domain) ToDomain() Domain {
	return Domain(strings.TrimSuffix(string(dm), "."))
}

func (dm Domain) MustToUnicode() Domain {
	res, err := idna.ToUnicode(string(dm))
	if err != nil {
		log.Error("MustToUnicode", "err", err, "domain", dm)
		return dm
	}
	dm = Domain(res)
	return dm
}

func (dm Domain) ToString() string {
	return string(dm)
}

func HasSuffix(s, suffix string) bool {
	s = Domain(s).MustToUnicode().ToFQDN().ToString()
	suffix = Domain(suffix).MustToUnicode().ToFQDN().ToString()
	return hasSuffix(s, suffix)
}

func TrimSuffix(s, suffix string) string {
	tmps := Domain(s).MustToUnicode().ToFQDN().ToString()
	suffix = Domain(suffix).MustToUnicode().ToFQDN().ToString()
	if hasSuffix(tmps, suffix) {
		return tmps[:len(tmps)-len(suffix)]
	}
	return s
}

// SplitSuffixToDomain 最长匹配截断，且返回的是域名，其后缀不带点； 除非suffix是根。
// 如：www.baidu.com -> www, baidu.com
// 如：www.baidu.com -> www.baidu.com, .
func SplitSuffixToDomain(s string, suffixs []string) (pre, suff string) {
	pre, suff = splitSuffixByLongestMatch(s, suffixs)
	if suff != "." {
		suff = Domain(suff).ToDomain().ToString()
	}
	return Domain(pre).ToDomain().ToString(), suff
}

// splitSuffixByLongestMatch 最长匹配截断
func splitSuffixByLongestMatch(s string, suffixs []string) (pre, suff string) {
	sufflen := 0
	s = Domain(s).MustToUnicode().ToFQDN().ToString()
	pre = s
	for _, suffix := range suffixs {
		suffix = Domain(suffix).MustToUnicode().ToFQDN().ToString()
		if hasSuffix(s, suffix) && len(suffix) > sufflen {
			pre = s[:len(s)-len(suffix)]
			suff = suffix
			sufflen = len(suffix)
		}
	}

	return pre, suff
}

// hasSuffix 是否包含后缀 域名大小写不敏感
func hasSuffix(s, suffix string) bool {
	if suffix == "" {
		return false
	}
	if suffix == "." && s != "" {
		return true
	}

	pre := s[:len(s)-len(suffix)]
	return len(s) >= len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix) &&
		(strings.HasSuffix(pre, ".") || pre == "")
}

func HostAddDomain(host, domain string) string {
	host = Domain(host).MustToUnicode().ToDomain().ToString()
	domain = Domain(domain).MustToUnicode().ToDomain().ToString()

	if domain == "" {
		return host
	}
	return host + "." + domain
}
