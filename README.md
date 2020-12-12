# Dallas DS18B20 Thermometer Prometheus Exporter

A simple http server written in Go which reads and exports the temperature of any sensors found on the Dallas 1-wire bus (e.g. a Raspberry Pi) in varz format during each http response.

Use this to export current temperature data from connected sensors into [Prometheus](https://prometheus.io/) for timeseries analysis.

## Usage

```
env GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o ds18b20_exporter
```

## Applying labels to Individual Sensors

You most likely want to include extra label metadata to which Prometheus can use to indicate things like sensor name, location etc.

Find the unique serial numbers of your sensors in `/sys/bus/w1/devices` and include a list of labels like this:

    $ ./ds18b20_exporter --port 8000 \
      --prometheus_labels "28-0416a4a474ff=location=lounge,"28-0417713760ff"=location=garden"

## Varz Output Format

    $ curl http://localhost:8000
    temperature_c{location="lounge",sensor="28-0416a4a474ff"} 18.437000
    temperature_c{location="garden",sensor="28-0417713760ff"} 12.500000

## License

ds18b20-prometheus-exporter is licensed under the MIT License.