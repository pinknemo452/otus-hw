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
	stats, err := getUsers(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return stats, nil
}

func getUsers(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	stats := make(DomainStat)

	for scanner.Scan() {
		line := scanner.Text()

		err := countUserDomain(fastjson.GetString([]byte(line), "Email"), domain, stats)
		if err != nil {
			return nil, err
		}
	}
	return stats, nil
}

func countUserDomain(email string, domain string, stat DomainStat) error {

	matched, err := regexp.Match("\\."+domain, []byte(email))
	if err != nil {
		return nil
	}

	if matched {
		emailDomain := strings.ToLower(strings.SplitN(email, "@", 2)[1])
		stat[emailDomain]++
	}

	return nil
}
