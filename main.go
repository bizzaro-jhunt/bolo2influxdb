package main

import (
	"fmt"
	"os"

	"github.com/bolo/go-bolo"
	"github.com/starkandwayne/goutils/log"
	"github.com/starkandwayne/metrics/influxdb"
	"github.com/voxelbrain/goptions"
)

var Version = "(development)"

func main() {
	options := struct {
		Debug   bool   `goptions:"-D, --debug, description='Enable debugging'"`
		Version bool   `goptions:"-v, --version, description='Display version information'"`
		Config  string `goptions:"-c, --config, description='Specify the config file for bolo2influxdb'"`
	}{
		Debug:   false,
		Version: false,
		Config:  "bolo2influxdb.conf",
	}
	err := goptions.Parse(&options)
	if err != nil {
		goptions.PrintHelp()
		os.Exit(1)
	}

	if options.Version {
		fmt.Printf("%s - Version %s\n", os.Args[0], Version)
		os.Exit(0)
	}

	logLevel := "info"
	if options.Debug {
		logLevel = "debug"
	}
	log.SetupLogging(log.LogConfig{Type: "console", File: "stderr", Level: logLevel})
	log.Infof("Starting up bolo2influxdb")

	cfg, err := LoadConfig(options.Config)
	if err != nil {
		log.Errorf("Unable to load config file %s: %s", options.Config, err)
		log.Errorf("Bailing out due to errors")
		os.Exit(1)
	}

	log.Debugf("Loaded Config: %#v", *cfg)

	log.Debugf("Connecting to influxdb")
	influxClient, err := influxdb.Connect(cfg.Influx)
	if err != nil {
		log.Errorf("Unable to connect to influx: %s: %s\n", cfg.Influx.Addr, err)
		log.Errorf("Bailing out due to errors")
		os.Exit(1)
	}

	log.Debugf("Connecting to Bolo")
	pduChan, errChan, err := bolo.Connect(fmt.Sprintf("tcp://%s:%s", cfg.Bolo.Addr, cfg.Bolo.Port))
	if err != nil {
		log.Errorf("Error connecting to Bolo: %s", err)
		log.Errorf("Bailing out due to errors")
		os.Exit(1)
	}

	go func() {
		for pdu := range pduChan {
			point, err := influxdb.PointFromBoloPDU(pdu)
			if err != nil {
				log.Errorf("Problem creating influx metric from Bolo data: %s", err)
				continue
			}

			// point is nil + no errors when non-metric event was processed
			if point != nil {
				err = influxClient.Send(point)
				if err != nil {
					log.Errorf("Problem submitting metric: %s", err)
				}
			}
		}
	}()

	for err := range errChan {
		log.Errorf("%s\n", err)
	}
	log.Errorf("Disconnected from Bolo after too many failures. bolo2influxdb is exiting")
}
