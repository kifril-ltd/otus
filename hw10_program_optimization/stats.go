package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/valyala/fastjson"
)

const (
	numWorkers = 2
	batchSize  = 500
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	stat := make(DomainStat)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		email := fastjson.GetString(scanner.Bytes(), "Email")

		if email != "" {
			splitEmail := strings.SplitN(email, "@", 2)
			if len(splitEmail) == 2 && strings.HasSuffix(splitEmail[1], "."+domain) {
				resultDomain := strings.ToLower(splitEmail[1])
				stat[resultDomain]++
			}
		}
	}
	return stat, nil
}

// GetDomainStatConcurrent Конкурентная обработка данных по батчам. Не получилось попасть в ограничение по памяти.
func GetDomainStatConcurrent(r io.Reader, domain string) (DomainStat, error) {
	usersCh := reader(r)

	domainCountersCh := make([]<-chan DomainStat, numWorkers)
	for i := 0; i < numWorkers; i++ {
		domainCountersCh[i] = domainCounter(usersCh, domain)
	}

	var mu sync.Mutex
	stat := make(DomainStat)

	for s := range statMerger(domainCountersCh...) {
		mu.Lock()
		for d, count := range s {
			stat[d] += count
		}
		mu.Unlock()
	}
	return stat, nil
}

func reader(r io.Reader) <-chan []string {
	out := make(chan []string)
	scanner := bufio.NewScanner(r)

	go func() {
		defer close(out)

		usersBatch := make([]string, 0, batchSize)
		for scanner.Scan() {
			if len(usersBatch) == batchSize {
				out <- usersBatch
				usersBatch = []string{}
			}
			usersBatch = append(usersBatch, scanner.Text())
		}

		if len(usersBatch) != 0 {
			out <- usersBatch
		}
	}()

	return out
}

func domainCounter(users <-chan []string, domain string) <-chan DomainStat {
	out := make(chan DomainStat)

	go func() {
		defer close(out)
		for userBatch := range users {
			stat := make(DomainStat)
			for _, user := range userBatch {
				email := fastjson.GetString([]byte(user), "Email")

				if email != "" {
					splitEmail := strings.SplitN(email, "@", 2)
					if len(splitEmail) == 2 && strings.HasSuffix(splitEmail[1], "."+domain) {
						resultDomain := strings.ToLower(splitEmail[1])
						stat[resultDomain]++
					}
				}
			}
			out <- stat
		}
	}()

	return out
}

func statMerger(stats ...<-chan DomainStat) <-chan DomainStat {
	out := make(chan DomainStat)

	var wg sync.WaitGroup
	multiplexer := func(p <-chan DomainStat) {
		defer wg.Done()

		for in := range p {
			out <- in
		}
	}

	wg.Add(len(stats))
	for _, in := range stats {
		go multiplexer(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
