package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/valyala/fastjson" //nolint:depguard
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainRegexp, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}
	stats, err := scanInputReader(r, domainRegexp)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return stats, nil
}

func scanInputReader(r io.Reader, domainRegexp *regexp.Regexp) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	stats := make(DomainStat)

	for scanner.Scan() {
		line := scanner.Text()
		email := fastjson.GetString([]byte(line), "Email")
		if email == "" {
			return nil, fmt.Errorf("no email field in json")
		}
		countUserDomain(email, domainRegexp, stats)
	}
	return stats, nil
}

func countUserDomain(email string, domainRegexp *regexp.Regexp, stat DomainStat) {
	matched := domainRegexp.Match([]byte(email))
	if matched {
		emailDomain := strings.ToLower(strings.SplitN(email, "@", 2)[1])
		stat[emailDomain]++
	}
}
