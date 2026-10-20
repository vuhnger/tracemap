package main

import (
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/vuhnger/tracemap/constants"
)

type Hop struct {
	TTL   int      `json:"ttl"`
	IP    *string  `json:"ip"`
	RTTMs *float64 `json:"rttMs"`
}

type Route struct {
	Target    string    `json:"target"`
	StartedAt time.Time `json:"startedAt"`
	Hops      []Hop     `json:"hops"`
}

var hopLine = regexp.MustCompile(`(?m)^ *(\d+) +(?:(\S+) +([\d.]+) ms|\*)`)

func parseTraceroute(output string) []Hop {
	hops := []Hop{}
	for _, match := range hopLine.FindAllStringSubmatch(output, -1) {
		ttl, _ := strconv.Atoi(match[1])
		hop := Hop{TTL: ttl}
		if match[2] != "" {
			rtt, _ := strconv.ParseFloat(match[3], 64)
			hop.IP, hop.RTTMs = new(match[2]), new(rtt)
		}
		hops = append(hops, hop)
	}
	return hops
}

func trace(target string) (Route, error) {
	startedAt := time.Now().UTC()
	output, err := exec.Command(
		"traceroute", "-I", "-n", "-q", "1",
		"-w", strconv.Itoa(constants.ProbeWaitSeconds),
		"-m", strconv.Itoa(constants.MaxHops),
		target,
	).Output()
	if err != nil {
		return Route{}, err
	}
	return Route{Target: target, StartedAt: startedAt, Hops: parseTraceroute(string(output))}, nil
}
