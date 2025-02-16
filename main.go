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
	Token  string `env:"INFLUXDB_TOKEN"`
	Url    string `env:"INFLUXDB_URL"`
	Org    string `env:"INFLUXDB_ORG"`
	Bucket string `env:"INFLUXDB_BUCKET"`
}

type DhtConfig struct {
	Pin   int `env:"DHT_PIN" envDefault:"4"`
	Retry int `env:"DHT_RETRY" envDefault:"10"`
}

type config struct {
	Influxdb InfluxDBConfig
	Dht      DhtConfig
}

func main() {

	cfg, err := env.ParseAs[config]()
	if err != nil {
		log.Printf("%+v\n", err)
	}

	client := influxdb2.NewClient(cfg.Influxdb.Url, cfg.Influxdb.Token)
	writeAPI := client.WriteAPIBlocking(cfg.Influxdb.Org, cfg.Influxdb.Bucket)

	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT22, cfg.Dht.Pin, false, cfg.Dht.Retry)
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
