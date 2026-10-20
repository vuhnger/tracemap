package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/vuhnger/tracemap/constants"
)

func readLocations() (Locations, error) {
	data, err := os.ReadFile(constants.LocationsPath)
	if errors.Is(err, fs.ErrNotExist) {
		return Locations{}, nil
	}
	if err != nil {
		return nil, err
	}
	var locations Locations
	return locations, json.Unmarshal(data, &locations)
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func printRoute(route Route, locations Locations) {
	table := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(table, "TTL\tIP\tRTT MS\tCITY")
	for _, hop := range route.Hops {
		ip, rtt, city := "", "", ""
		if hop.IP != nil {
			ip, rtt = *hop.IP, fmt.Sprint(*hop.RTTMs)
			if location := locations[ip]; location != nil {
				city = location.City
			}
		}
		fmt.Fprintf(table, "%d\t%s\t%s\t%s\n", hop.TTL, ip, rtt, city)
	}
	table.Flush()
}

func run(target string) error {
	route, err := trace(target)
	if err != nil {
		return err
	}
	known, err := readLocations()
	if err != nil {
		return err
	}
	ips := []string{}
	for _, hop := range route.Hops {
		if hop.IP != nil {
			ips = append(ips, *hop.IP)
		}
	}
	locations, err := locate(constants.IPMapLocateURL, ips, known)
	if err != nil {
		return err
	}
	if err := writeJSON(constants.LocationsPath, locations); err != nil {
		return err
	}
	routePath := fmt.Sprintf("%s/%s-%d.json", constants.RoutesDir, target, time.Now().UnixMilli())
	if err := writeJSON(routePath, route); err != nil {
		return err
	}
	printRoute(route, locations)
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: tracemap <host>")
	}
	if err := run(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
