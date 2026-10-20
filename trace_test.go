package main

import (
	"os"
	"reflect"
	"testing"
)

func parseFixture(t *testing.T, name string) []Hop {
	t.Helper()
	output, err := os.ReadFile("fixtures/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return parseTraceroute(string(output))
}

func answered(ttl int, ip string, rttMs float64) Hop {
	return Hop{TTL: ttl, IP: new(ip), RTTMs: new(rttMs)}
}

func TestParseLinuxSkipsHeader(t *testing.T) {
	want := []Hop{
		answered(1, "5.101.109.1", 5.601),
		answered(2, "143.244.192.104", 3.784),
		answered(3, "143.244.224.110", 2.453),
		answered(4, "143.244.224.107", 1.054),
		answered(5, "162.158.84.196", 16.246),
		answered(6, "162.158.84.78", 2.345),
		{TTL: 7},
		answered(8, "1.1.1.1", 2.641),
	}
	if got := parseFixture(t, "linux.txt"); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestParseMacOSWithSilentHops(t *testing.T) {
	want := []Hop{answered(1, "192.168.1.254", 4.059)}
	for ttl := 2; ttl <= 10; ttl++ {
		want = append(want, Hop{TTL: ttl})
	}
	want = append(want, answered(11, "1.1.1.1", 39.38))
	if got := parseFixture(t, "macos.txt"); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
