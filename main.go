package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

type DomainStats struct {
	Success int
	Total   int
}

const threshold_response_time = 500
const threshold_timeout_duration = 500
//Set above thresholds accordingly in MilliSeconds(ms)

func checkHealth(endpoint Endpoint,stats map[string]*DomainStats) {
	if endpoint.Method == "" {
		endpoint.Method = "GET" // Default to GET if method is omitted
	}

	var client = &http.Client{
		Timeout: threshold_timeout_duration * time.Millisecond, //Capping Request duration at 500ms with Timeout
	}

	var req *http.Request
	var err error

	if endpoint.Body != "" && endpoint.Method != "GET" {
		req, err = http.NewRequest(endpoint.Method, endpoint.URL, bytes.NewBufferString(endpoint.Body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(endpoint.Method, endpoint.URL, nil) // Drafting request with no body for GET methods
	}

	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
	}

	request_start := time.Now()
	resp, err := client.Do(req)
	actual_request_duration := time.Since(request_start)

	domain := extractDomain(endpoint.URL)

	stats[domain].Total++
	if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 && actual_request_duration <= threshold_response_time*time.Millisecond {
		stats[domain].Success++
	}
}

func extractDomain(url string) string {
	urlSplit := strings.Split(url, "//")
	domain := strings.Split(urlSplit[len(urlSplit)-1], "/")[0]
	domain_without_port := strings.Split(domain, ":")[0] // Stripping port if present
	return domain_without_port
}

func monitorEndpoints(endpoints []Endpoint) {
	stats := make(map[string]*DomainStats)
	for _, endpoint := range endpoints {
		domain := extractDomain(endpoint.URL)
		if stats[domain] == nil {
			stats[domain] = &DomainStats{}
		}
	}
		for _, endpoint := range endpoints {
			checkHealth(endpoint,stats)
		}
		logResults(stats)
}

func logResults(stats map[string]*DomainStats) {
	for domain, stat := range stats {
		percentage := int(math.Round(100 * float64(stat.Success) / float64(stat.Total)))
		fmt.Printf("%s has %d%% availability\n", domain, percentage)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <config_file>")
	}

	filePath := os.Args[1]
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var endpoints []Endpoint
	if err := yaml.Unmarshal(data, &endpoints); err != nil {
		log.Fatal("Error parsing YAML:", err)
	}

	for {
		go monitorEndpoints(endpoints)
		time.Sleep(15 * time.Second)
	}
}
