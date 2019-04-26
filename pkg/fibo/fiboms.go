package fibo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func getEnv(name string, fallback string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		value = fallback
	}
	return value
}

type target struct {
	url  *url.URL
	host string
}

var defaultTarget = newTarget()

func newTarget() *target {
	parsedURL, err := url.Parse(getEnv("FIBO_URL", "http://localhost:8080/fibo/microservice"))
	if err != nil {
		log.Fatal(err)
	}
	return &target{parsedURL, getEnv("FIBO_HOST", "")}
}

func (t *target) invoke(n uint) (uint64, error) {
	q := t.url.Query()
	q.Set("n", strconv.FormatUint(uint64(n), 10))
	t.url.RawQuery = q.Encode()
	client := &http.Client{}
	request, err := http.NewRequest("GET", t.url.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Unexpected statuscode: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(body)), 10, 64)
}

// Microservice callculation of fibonacci number n (NOTE: produces exponential load)
func Microservice(n uint) uint64 {
	var result uint64
	if n < 2 {
		return uint64(n)
	}

	for i := n - 2; i < n; i++ {
		r, err := defaultTarget.invoke(i)
		if err != nil {
			log.Fatal(err)
		}
		result += r
	}
	return result
}
