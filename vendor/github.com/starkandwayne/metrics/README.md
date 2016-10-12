# Overview

This is a set of golang libraries and applications revolving around making
metric gathering and submission easier.


*firehose2influxdb* is a go application that connects to Cloud Foundry's Loggregator
Firehose to influxdb

*influxdb* is a wrapper library on top of the [influxdata/influxdb](https://github.com/influxdata/influxdb) library,
which makes things a little easier to work with when handling streams of data like the Firehose.


# Caveats
This is all currently a work in progress. Check out the Issues tab to see what nees to be finishing.

# Using influxdb

`go get github.com/starkandwayne/metrics/influxdb`

# Using firehose2influxdb

```
go get github.com/starkandwayne/metrics/firehose2influxdb
firehose2influxdb -c cfg.json
```

## Config Format
```
{
	"influx": {
		"url":"http://your.influx.ip:port",
		"user":"influx-user",
		"password:"influx-password",
		"database":"influx-database"
	},
	"cf": {
		"api_url": "https://api.system.cloudfoundry.domain",
		"user": "cf-user",
		"password": "cf-password"
	}
}
