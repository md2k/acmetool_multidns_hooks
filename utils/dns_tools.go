package utils

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
	"golang.org/x/net/publicsuffix"
)

var DNSTimeout = 10 * time.Second

var RecursiveNameservers = []string{
	"google-public-dns-a.google.com:53",
	"google-public-dns-b.google.com:53",
}

// DNS01Record returns a DNS record which will fulfill the `dns-01` challenge
func DNS01Record(domain string) (fqdn string) {
	fqdn = fmt.Sprintf("_acme-challenge.%s.", domain)
	return
}

// FindZoneByFqdn determines the zone of the given fqdn
func FindZoneByFqdn(fqdn string) (string, error) {

	// Query the authoritative nameserver for a hopefully non-existing SOA record,
	// in the authority section of the reply it will have the SOA of the
	// containing zone. rfc2308 has this to say on the subject:
	//   Name servers authoritative for a zone MUST include the SOA record of
	//   the zone in the authority section of the response when reporting an
	//   NXDOMAIN or indicating that no data (NODATA) of the requested type exists
	in, err := dnsQuery(fqdn, dns.TypeSOA, RecursiveNameservers, true)
	if err != nil {
		return "", err
	}
	if in.Rcode != dns.RcodeNameError {
		if in.Rcode != dns.RcodeSuccess {
			return "", fmt.Errorf("The NS returned %s for %s", dns.RcodeToString[in.Rcode], fqdn)
		}
		// We have a success, so one of the answers has to be a SOA RR
		for _, ans := range in.Answer {
			if soa, ok := ans.(*dns.SOA); ok {
				return checkIfTLD(fqdn, soa)
			}
		}
		// Or it is NODATA, fall through to NXDOMAIN
	}
	// Search the authority section for our precious SOA RR
	for _, ns := range in.Ns {
		if soa, ok := ns.(*dns.SOA); ok {
			return checkIfTLD(fqdn, soa)
		}
	}
	return "", fmt.Errorf("The NS did not return the expected SOA record in the authority section")
}

// dnsQuery will query a nameserver, iterating through the supplied servers as it retries
// The nameserver should include a port, to facilitate testing where we talk to a mock dns server.
func dnsQuery(fqdn string, rtype uint16, nameservers []string, recursive bool) (in *dns.Msg, err error) {
	m := new(dns.Msg)
	m.SetQuestion(fqdn, rtype)
	m.SetEdns0(4096, false)

	if !recursive {
		m.RecursionDesired = false
	}

	// Will retry the request based on the number of servers (n+1)
	for i := 1; i <= len(nameservers)+1; i++ {
		ns := nameservers[i%len(nameservers)]
		udp := &dns.Client{Net: "udp", Timeout: DNSTimeout}
		in, _, err = udp.Exchange(m, ns)

		if err == dns.ErrTruncated {
			tcp := &dns.Client{Net: "tcp", Timeout: DNSTimeout}
			// If the TCP request suceeds, the err will reset to nil
			in, _, err = tcp.Exchange(m, ns)
		}

		if err == nil {
			break
		}
	}
	return
}

// lookupNameservers returns the authoritative nameservers for the given fqdn.
func lookupNameservers(fqdn string) ([]string, error) {
	var authoritativeNss []string

	zone, err := FindZoneByFqdn(fqdn)
	if err != nil {
		return nil, err
	}

	r, err := dnsQuery(zone, dns.TypeNS, RecursiveNameservers, true)
	if err != nil {
		return nil, err
	}

	for _, rr := range r.Answer {
		if ns, ok := rr.(*dns.NS); ok {
			authoritativeNss = append(authoritativeNss, strings.ToLower(ns.Ns))
		}
	}

	if len(authoritativeNss) > 0 {
		return authoritativeNss, nil
	}
	return nil, fmt.Errorf("Could not determine authoritative nameservers")
}

// checkDNSPropagation checks if the expected TXT record has been propagated to all authoritative nameservers.
func CheckDNSPropagation(fqdn, value string) (bool, error) {
	// Initial attempt to resolve at the recursive NS
	r, err := dnsQuery(fqdn, dns.TypeTXT, RecursiveNameservers, true)
	if err != nil {
		return false, err
	}
	if r.Rcode == dns.RcodeSuccess {
		// If we see a CNAME here then use the alias
		for _, rr := range r.Answer {
			if cn, ok := rr.(*dns.CNAME); ok {
				if cn.Hdr.Name == fqdn {
					fqdn = cn.Target
					break
				}
			}
		}
	}

	authoritativeNss, err := lookupNameservers(fqdn)
	if err != nil {
		return false, err
	}

	return checkAuthoritativeNss(fqdn, value, authoritativeNss)
}

// checkAuthoritativeNss queries each of the given nameservers for the expected TXT record.
func checkAuthoritativeNss(fqdn, value string, nameservers []string) (bool, error) {
	for _, ns := range nameservers {
		r, err := dnsQuery(fqdn, dns.TypeTXT, []string{net.JoinHostPort(ns, "53")}, false)
		if err != nil {
			return false, err
		}

		if r.Rcode != dns.RcodeSuccess {
			return false, fmt.Errorf("NS %s returned %s for %s", ns, dns.RcodeToString[r.Rcode], fqdn)
		}

		var found bool
		for _, rr := range r.Answer {
			if txt, ok := rr.(*dns.TXT); ok {
				if strings.Join(txt.Txt, "") == value {
					found = true
					break
				}
			}
		}

		if !found {
			return false, fmt.Errorf("NS %s did not return the expected TXT record", ns)
		}
	}

	return true, nil
}

func checkIfTLD(fqdn string, soa *dns.SOA) (string, error) {
	zone := soa.Hdr.Name
	// If we ended up on one of the TLDs, it means the domain did not exist.
	publicsuffix, _ := publicsuffix.PublicSuffix(UnFqdn(zone))
	if publicsuffix == UnFqdn(zone) {
		return "", fmt.Errorf("Could not determine zone authoritatively")
	}
	return zone, nil
}

// ToFqdn converts the name into a fqdn appending a trailing dot.
func ToFqdn(name string) string {
	n := len(name)
	if n == 0 || name[n-1] == '.' {
		return name
	}
	return name + "."
}

// UnFqdn converts the fqdn into a name removing the trailing dot.
func UnFqdn(name string) string {
	n := len(name)
	if n != 0 && name[n-1] == '.' {
		return name[:n-1]
	}
	return name
}
