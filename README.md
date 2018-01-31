# Dallas DS18B20 Thermometer Prometheus Exporter

A simple http server written in Go which reads and exports the temperature of any sensors found on the Dallas 1-wire bus (e.g. a Raspberry Pi) in varz format during each http response.

Use this to export current temperature data from connected sensors into Prometheues for timeseries analysis.


## Compatible Sensors

http://amzn.to/2jGRjKO

