package temperature

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

//var bus_dir = flag.String("w1_bus_dir", "/sys/bus/w1/devices", "directory of the 1-wire bus")
var bus_dir = flag.String("w1_bus_dir", "src/github.com/samkalnins/ds18b20-thermometer-prometheus-exporter/fixtures/w1_devices", "directory of the 1-wire bus")
var port = flag.Integer("port", 8000, "port to run http server on")

// temperature_c{location="garden",location_type="outside",sensor="28-0417713760ff"} 20

const w1_slave_fname = "w1_slave"

type TempReading struct {
	id     string
	temp_c float64
}

type PrometheusLabel struct {
	temp_id string
	name    string
	value   string
}

type prometheusLabels []PrometheusLabel

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (p *prometheusLabels) String() string {
	return fmt.Sprint(*p)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (p *prometheusLabels) Set(value string) error {
	for _, ls := range strings.Split(value, ",") {
		s := strings.Split(ls, "=")
		if len(s) != 3 {
			errors.New("Bad flag value -- should be temp_id=label=value")
		}
		*p = append(*p, PrometheusLabel{s[0], s[1], s[2]})
	}
	return nil
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

func getLabelsMap(labels prometheusLabels) map[string][]string {
	out := make(map[string][]string)
	for _, label := range labels {
		out[label.temp_id] = append(out[label.temp_id], fmt.Sprintf("%s=\"%s\"", label.name, label.value))
	}
	return out
}

var prometheusLabelsFlag prometheusLabels

func init() {
	flag.Var(&prometheusLabelsFlag, "prometheus_labels", "comma-separated list of labels to apply to sensors e.g. sensor_id_1234=label_a=bar,")
}

func main() {
	flag.Parse()
	labelMap := getLabelsMap(prometheusLabelsFlag)

	readings, err := FindAndReadTemperatures(*bus_dir)
	if err != nil {
		log.Fatal("Error reading temperatures [%s]", err)
	}

	for _, tr := range readings {
		labels := strings.Join(labelMap[tr.id], ",")
		log.Printf("Read sensor %s = %.2f degress C {%s}\n", tr.id, tr.temp_c, labels)
	}
}
