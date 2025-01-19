package main

import (
	"fmt"
	"net"
)

// getFirstMXRecord retrieves the first MX record for the domain with caching.
func getFirstMXRecord(domain string, noCache bool) ([]string, bool, error) {
    if !noCache && config.CacheEnabled {
        if mxRecords, found := cache.GetMX(domain); found {
            return mxRecords, true, nil
        }
    }

    mxRecords, err := net.LookupMX(domain)
    if err != nil || len(mxRecords) == 0 {
        return nil, false, fmt.Errorf("No MX record found for domain %s", domain)
    }

    var hosts []string
    for _, mx := range mxRecords {
        hosts = append(hosts, mx.Host)
    }

    if !noCache && config.CacheEnabled {
        cache.SetMX(domain, hosts)
    }

    return hosts, false, nil
}
