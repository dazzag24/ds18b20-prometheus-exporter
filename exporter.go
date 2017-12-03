package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/samkalnins/ds18b20-thermometer-prometheus-exporter/temperature"

	"log"
	"net/http"
	"strings"
)

//var bus_dir = flag.String("w1_bus_dir", "/sys/bus/w1/devices", "directory of the 1-wire bus")
var bus_dir = flag.String("w1_bus_dir", "src/github.com/samkalnins/ds18b20-thermometer-prometheus-exporter/fixtures/w1_devices", "directory of the 1-wire bus")
var port = flag.Int("port", 8000, "port to run http server on")

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		readings, err := temperature.FindAndReadTemperatures(*bus_dir)
		if err != nil {
			log.Print("Error reading temperatures [%s]", err)
			// 500
		}

		for _, tr := range readings {
			labels := strings.Join(append(labelMap[tr.Id], fmt.Sprintf("sensor=\"%s\"", tr.Id)), ",")
			log.Printf("Read sensor %s = %.2f degress C {%s}\n", tr.Id, tr.Temp_c, labels)

			// Output varz as both C & F for maximum user happiness
			fmt.Fprintf(w, "temperature_c{%s} %f\n", labels, tr.Temp_c)
			fmt.Fprintf(w, "temperature_f{%s} %f\n", labels, temperature.CentigradeToF(tr.Temp_c))
		}

	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
