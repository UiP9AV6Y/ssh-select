package provider

import (
	"github.com/miekg/dns"
	"os"
	"strings"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
)

var (
	// the DNS module does not expose this as constant
	fullyQualifiedIndicator = "."
)

type ZoneFileProvider struct {
	file string
}

func (p *ZoneFileProvider) String() string {
	return p.file
}

func (p *ZoneFileProvider) Parse() ([]remote.Host, error) {
	var fd *os.File
	var err error

	fd, err = os.Open(p.file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	parser := dns.NewZoneParser(fd, "", p.file)
	return p.parseFile(parser)
}

func (p *ZoneFileProvider) parseFile(parser *dns.ZoneParser) ([]remote.Host, error) {
	var hosts []remote.Host
	entries, err := p.parseZone(parser)

	if err != nil {
		return hosts, err
	}

	hosts = make([]remote.Host, len(entries))

	for i, entry := range entries {
		fqdn := strings.TrimRight(entry, fullyQualifiedIndicator)

		hosts[i] = remote.NewSimpleHost(fqdn)
	}

	return hosts, nil
}

func (p *ZoneFileProvider) parseZone(parser *dns.ZoneParser) ([]string, error) {
	var err error
	hosts := []string{}
	sink := func(host string) {
		hosts = append(hosts, host)
	}

	for rr, ok := parser.Next(); ok; rr, ok = parser.Next() {
		if err = p.parseRR(rr, sink); err != nil {
			break
		}
	}

	return hosts, err
}

func (p *ZoneFileProvider) parseRR(rr dns.RR, host func(string)) error {
	switch v := rr.(type) {
	case *dns.PTR:
		host(v.Ptr)
		// the zone $ORIGIN is usually not a plain address,
		// therefor it makes no sense to provide it as host
		// (e.g. 1.0.10.in-addr.arpa.)
		//host(v.Header().Name)
	case *dns.A:
		if v.A != nil {
			host(v.A.String())
		}
		host(v.Header().Name)
	case *dns.AAAA:
		if v.AAAA != nil {
			host(v.AAAA.String())
		}
		host(v.Header().Name)
	case *dns.NS:
		host(v.Ns)
		host(v.Header().Name)
	case *dns.MX:
		host(v.Mx)
		host(v.Header().Name)
	case *dns.CNAME:
		host(v.Target)
		host(v.Header().Name)
	}

	return nil
}

func NewZoneFileProvider(file string) *ZoneFileProvider {
	provider := &ZoneFileProvider{
		file: file,
	}

	return provider
}
