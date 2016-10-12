package main

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/starkandwayne/goutils/log"
	"github.com/starkandwayne/metrics/influxdb"
	"github.com/voxelbrain/goptions"
)

var Version = "(development)"

func main() {
	options := struct {
		Debug   bool   `goptions:"-D, --debug, description='Enable debugging'"`
		Version bool   `goptions:"-v, --version, description='Display version information'"`
		Config  string `goptions:"-c, --config, description='Specify the config file for firehose2influxdb'"`
	}{
		Debug:   false,
		Version: false,
		Config:  "firehose2influxdb.conf",
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
	log.Infof("Starting up firehose2influxdb")

	cfg, err := LoadConfig(options.Config)
	if err != nil {
		log.Errorf("Unable to load config file %s: %s", options.Config, err)
		log.Errorf("Bailing out due to errors")
		os.Exit(1)
	}

	log.Debugf("Loaded Config: %#v", *cfg.CF)

	log.Debugf("Connecting to influxdb")
	influxClient, err := influxdb.Connect(cfg.Influx)
	if err != nil {
		log.Errorf("Unable to connect to influx: %s: %s\n", cfg.Influx.Addr, err)
		log.Errorf("Bailing out due to errors")
		os.Exit(1)
	}

	log.Debugf("Connecting to CloudFoundry")
	cf, err := cfclient.NewClient(cfg.CF)
	if err != nil {
		log.Errorf("Unable to connect to %s: %s\n", cfg.CF.ApiAddress, err)
		log.Errorf("Bailing out due to errors")
		os.Exit(1)
	}
	log.Debugf("Connected to CF. Requesting token from UAA for the firehose")
	token, err := cf.GetToken()
	if err != nil {
		log.Errorf("ERROR retrieving token from %s: %s\n", cf.Endpoint.TokenEndpoint, err)
		log.Errorf("Bailing out due to errors")
		os.Exit(1)
	}
	noaa := consumer.New(cf.Endpoint.DopplerEndpoint, &tls.Config{InsecureSkipVerify: cfg.CF.SkipSslValidation}, nil)

	log.Debugf("Hooking up a firehose stream")
	evtChan, errChan := noaa.Firehose("firehose2influxdb", token)
	go func() {
		for evt := range evtChan {
			// For influxdb, submit to influx
			point, err := influxdb.PointsFromFirehoseEnvelope(evt)
			if err != nil {
				log.Errorf("Problem creating influx metric from firehose data: %s", err)
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

	log.Errorf("Disconnected from Firehose after too many failures. firehose2influxdb is exiting")
}
