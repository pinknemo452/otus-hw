package hw10programoptimization

import (
	"github.com/valyala/fastjson"

	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainRegexp, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}
	stats, err := getUsers(r, domainRegexp)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return stats, nil
}

func getUsers(r io.Reader, domainRegexp *regexp.Regexp) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	stats := make(DomainStat)

	for scanner.Scan() {
		line := scanner.Text()

		err := countUserDomain(fastjson.GetString([]byte(line), "Email"), domainRegexp, stats)
		if err != nil {
			return nil, err
		}
	}
	return stats, nil
}

func countUserDomain(email string, domainRegexp *regexp.Regexp, stat DomainStat) error {
	matched := domainRegexp.Match([]byte(email))
	if matched {
		emailDomain := strings.ToLower(strings.SplitN(email, "@", 2)[1])
		stat[emailDomain]++
	}

	return nil
}
