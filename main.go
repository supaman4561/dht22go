package main

import (
	"context"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/d2r2/go-dht"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxDBConfig struct {
	token  string `env:"INFLUXDB_TOKEN"`
	url    string `env:"INFLUXDB_URL"`
	org    string `env:"INFLUXDB_ORG"`
	bucket string `env:"INFLUXDB_BUCKET"`
}

type DhtConfig struct {
	pin   int `env:"DHT_PIN" envDefault:"4"`
	retry int `env:"DHT_RETRY" envDefault:"10"`
}

type config struct {
	influxdb InfluxDBConfig
	dht      DhtConfig
}

func main() {

	cfg, err := env.ParseAs[config]()

	client := influxdb2.NewClient(cfg.influxdb.url, cfg.influxdb.token)
	writeAPI := client.WriteAPIBlocking(cfg.influxdb.org, cfg.influxdb.bucket)

	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT22, 4, false, cfg.dht.retry)
	if err != nil {
		log.Fatal(err)
	}

	tags := map[string]string{
		"tagname1": "temperature",
		"tagname2": "humidity",
	}

	fields := map[string]interface{}{
		"temperature": temperature,
		"humidity":    humidity,
	}

	point := write.NewPoint("home_measurement", tags, fields, time.Now())
	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		log.Fatal(err)
	}
}
