package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

//var bus_dir = flag.String("w1_bus_dir", "/sys/bus/w1/devices", "directory of the 1-wire bus")
var bus_dir = flag.String("w1_bus_dir", "src/github.com/samkalnins/ds18b20-thermometer-prometheus-exporter/fixtures/w1_devices", "directory of the 1-wire bus")

// temperature_c{location="garden",location_type="outside",sensor="28-0417713760ff"} 20

const w1_slave_fname = "w1_slave"

type TempReading struct {
	id     string
	temp_c float64
}

func ReadTemperatureFile(path string) (float64, error) {
	var temp_c float64
	var err error

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return temp_c, err
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

func FindAndReadTemperatures(path string) ([]TempReading, error) {
	out := make([]TempReading, 0)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("Error reading directory %s\n", path)
		return out, err
	}

	for _, file := range files {
		t_file := filepath.Join(path, file.Name(), w1_slave_fname)
		temp_c, err := ReadTemperatureFile(t_file)
		if err == nil {
			t := TempReading{}
			t.id = file.Name()
			t.temp_c = temp_c
			out = append(out, t)
		}
	}
	return out, err
}

func main() {
	flag.Parse()

	readings, err := FindAndReadTemperatures(*bus_dir)
	if err != nil {
		log.Fatal("Error reading temperatures [%s]", err)
	}

	for _, tr := range readings {
		log.Printf("Read sensor %s = %.2f degress C\n", tr.id, tr.temp_c)
	}
}
