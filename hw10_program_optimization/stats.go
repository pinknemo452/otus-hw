package hw10programoptimization

import (
	"bufio"
	"encoding/json"
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
		var user User
		if err := json.Unmarshal([]byte(line), &user); err != nil {
			return nil, err
		}

		err := countUserDomain(user, domain, stats)
		if err != nil {
			return nil, err
		}
	}
	return stats, nil
}

func countUserDomain(u User, domain string, stat DomainStat) error {

	matched, err := regexp.Match("\\."+domain, []byte(u.Email))
	if err != nil {
		return nil
	}

	if matched {
		emailDomain := strings.ToLower(strings.SplitN(u.Email, "@", 2)[1])
		stat[emailDomain]++
	}

	return nil
}
