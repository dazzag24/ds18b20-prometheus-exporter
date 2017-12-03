package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//var bus_dir = flag.String("w1_bus_dir", "/sys/devices/w1_bus_master1", "directory of the 1-wire bus")
var bus_dir = flag.String("w1_bus_dir", "src/github.com/samkalnins/ds18b20-thermometer-prometheus-exporter/fixtures/w1_devices", "directory of the 1-wire bus")

var location = flag.String("location", "default_location", "temperature sensor location text label")
var location_type = flag.String("location_type", "inside", "temperature sensor location type (inside or outside)")

func readTemperatureFile(path string) (float64, error) {
	var temp_c float64
	var err error

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")
	if strings.HasSuffix(lines[0], "YES") && strings.Contains(lines[1], "t=") {
		i, err := strconv.ParseFloat(strings.Split(lines[1], "t=")[1], 64)
		temp_c = i / 1000.0
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = errors.New("Unparseable temperature file")
	}
	return temp_c, err
}

func processPath(path string, info os.FileInfo, err error) error {
	if strings.HasSuffix(path, "/w1_slave") {
		p := strings.Split(path, "/")
		id := p[len(p)-2]

		// Check file path for contents
		temp_c, err := readTemperatureFile(path)
		if err == nil {
			log.Printf("Found temp sensor %s (currently %.2f degress C)\n", id, temp_c)
		} else {
			log.Printf("Error processing temp sensor %s [%s]\n", id, err)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	log.Printf("Walking w1 bus dir %s\n", *bus_dir)
	filepath.Walk(*bus_dir, processPath)
}
