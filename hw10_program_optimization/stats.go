package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/jmequai/otus-golang/hw10_program_optimization/model"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	stat := make(DomainStat)

	if domain == "" {
		return stat, errors.New("domain must not be empty")
	}

	suffix := "." + domain

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()
		user := &model.User{}

		if err := user.UnmarshalJSON(line); err != nil {
			return stat, err
		}

		if strings.HasSuffix(user.Email, suffix) {
			parts := strings.SplitN(user.Email, "@", 3)

			if len(parts) != 2 {
				return stat, fmt.Errorf("invalid email %q", user.Email)
			}

			key := strings.ToLower(parts[1])

			stat[key]++
		}
	}

	return stat, nil
}
