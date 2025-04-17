# Fetch Take Home Exercise - Go based Endpoint Availability Monitor

This Go application is an enhanced/fixed version of the base program provided in the [Fetch SRE take-home exercise](https://github.com/fetch-rewards/sre-take-home-exercise-go).

It's accuracy has been enhanced to check the availability of HTTP endpoints defined in a YAML configuration file and logs domain-level availability every 15 seconds. It has been updated to meet all functional requirements, handling YAML structure, domain normalization, and concurrency for scalable check cycles.

---

## Installation & Running Instructions

### 1. Install Go (Prerequisite)

- Download manually from the official Go website: [https://go.dev/dl/](https://go.dev/dl/).
- Run the Installer and follow the instructions.

### 2. Get the Code

- Clone the repository using Git (recommended)

```
git clone https://github.com/praveenchukkapalli/Fetch-sre-exercise-praveen_chukkapalli.git
cd Fetch-sre-exercise-praveen_chukkapalli
```

- Download the ZIP from the [repository page](https://github.com/praveenchukkapalli/Fetch-sre-exercise-praveen_chukkapalli.git) and extract it manually



### 3. Run the Program

```bash
# Download and install any dependencies needed as per go.mod
go mod tidy
```

```bash
# Run the Go program with the provided YAML configuration file
go run main.go sample.yaml 
```

* `main.go` GO program starter that loads the YAML configuration file, runs availability checks and logs the stats.
* `sample.yaml` is a configuration file that defines a list of HTTP endpoints to be monitored, including method, headers, and optional request body.

---

## Issues Identified & Fixes as per requirements provided.

### 1. Perform Healt check for every 15 seconds diregard of HTTP request durations in each cycle (relatively complex).

- **Original**: Monitoring loop ran synchronously and could delay based on endpoint speed.
- **Fix**: Launched `monitorEndpoints()` in a goroutine every 15 seconds using `time.Sleep`. This ensured regular intervals regardless of execution time and allowed concurrent check cycles, If needed.

- This Caused Data racing and corruption when program is stressed with larger no.of endpoints due to shared  `stats` map was shared across all monitoring cycles,s.
    -  **To Fix this,** Defined `stats` inside `monitorEndpoints()` and passed it to other functions. This avoided shared state and enabled each 15s cycle to run without interference.

### 3. Must determine availability as 2xx + ≤500ms

- **Original**: Only the HTTP status code was considered to determine success.
- **Fix**: Added timing logic using `time.Since()` and also customizable 500ms timeout on the HTTP client. This logs Response duration for future scalability and Timeout also restricts the heap buildup with tangled/unresponsive requests.

### 4. Must not send body with GET requests + Must use GET method as default, if missing in YAML file.

- **Original**: The body was sent with every request, and missing method values caused request errors.
- **Fix**: Set the default method to `GET` when missing and ensured the body is only included when the method is not GET and body is non-empty.

### 5. Must ignore port numbers in domain

- **Original**: Domains were tracked with ports included, e.g., `praveenchukkapalli.com:8080`.
- **Fix**: Used `strings.Split` to remove ports from the domain string before grouping. A small string manipulation fix.

---
