package influxdb

import (
	"fmt"
	"time"

	"github.com/bolo/go-bolo"
	"github.com/cloudfoundry/sonde-go/events"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/starkandwayne/goutils/log"
)

type Client struct {
	client influx.Client
	cfg    *Config
	bp     *influx.BatchPoints
	Data   chan *influx.Point
	Errs   chan error
}

func Connect(opts Config) (*Client, error) {
	cfg := influx.HTTPConfig{
		Addr:               opts.Addr,
		Username:           opts.User,
		Password:           opts.Password,
		UserAgent:          "firehose2influxdb",
		Timeout:            1 * time.Second,
		InsecureSkipVerify: opts.InsecureSkipVerify,
	}

	client, err := influx.NewHTTPClient(cfg)
	if err != nil {
		return nil, err
	}

	/* FIXME
	// set up channels

	go func() {
		// Read off data channel
		// after 1000 events, submit to client, reset batchpoints
	}()
	*/
	return &Client{client: client, cfg: &opts}, nil
}

func (c *Client) Send(point *influx.Point) error {
	// FIXME: Send point over channel
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{Database: c.cfg.Database})
	if err != nil {
		return err
	}

	bp.AddPoint(point)

	err = c.client.Write(bp)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	// FIXME close channels
	return c.client.Close()
}

func PointsFromFirehoseEnvelope(env *events.Envelope) (*influx.Point, error) {
	var point *influx.Point
	var err error

	tags := env.GetTags()
	if tags == nil {
		tags = map[string]string{}
	}
	tags["origin"] = env.GetOrigin()
	tags["job"] = env.GetJob()
	tags["index"] = env.GetIndex()
	tags["deployment"] = env.GetDeployment()
	tags["ip_addr"] = env.GetIp()

	ts := time.Unix(0, env.GetTimestamp())

	switch *env.EventType {
	case events.Envelope_CounterEvent:
		evt := env.GetCounterEvent()
		fields := map[string]interface{}{
			"Delta": evt.GetDelta(),
			"Total": evt.GetTotal(),
		}
		point, err = influx.NewPoint(evt.GetName(), tags, fields, ts)
		if err != nil {
			return nil, err
		}
	case events.Envelope_ContainerMetric:
		metric := env.GetContainerMetric()
		tags["app_guid"] = metric.GetApplicationId()
		tags["app_instance"] = fmt.Sprintf("%d", metric.GetInstanceIndex())
		fields := map[string]interface{}{
			"CPU":         metric.GetCpuPercentage(),
			"Memory":      metric.GetMemoryBytes(),
			"MemoryQuota": metric.GetMemoryBytesQuota(),
			"Disk":        metric.GetDiskBytes(),
			"DiskQuota":   metric.GetDiskBytesQuota(),
		}
		point, err = influx.NewPoint("ContainerHealth", tags, fields, ts)
		if err != nil {
			return nil, err
		}
	case events.Envelope_ValueMetric:
		metric := env.GetValueMetric()
		fields := map[string]interface{}{
			"Value": metric.GetValue(),
		}
		point, err = influx.NewPoint(metric.GetName(), tags, fields, ts)
		if err != nil {
			return nil, err
		}
	}
	return point, nil
}

func PointFromBoloPDU(pdu bolo.PDU) (*influx.Point, error) {
	var point *influx.Point
	var err error

	tags := map[string]string{}

	switch pdu.Type() {
	case bolo.SAMPLE:
		sample := pdu.(*bolo.SamplePDU)
		fields := map[string]interface{}{
			"Min":      sample.Min,
			"Max":      sample.Max,
			"Sum":      sample.Sum,
			"Samples":  sample.SampleSize,
			"Mean":     sample.Mean,
			"Variance": sample.Variance,
		}
		point, err = influx.NewPoint(sample.Name, tags, fields, sample.Timestamp)
		if err != nil {
			return nil, err
		}

	case bolo.RATE:
		rate := pdu.(*bolo.RatePDU)
		fields := map[string]interface{}{
			"Window": rate.Window,
			"Value":  rate.Value,
		}
		point, err = influx.NewPoint(rate.Name, tags, fields, rate.Timestamp)
		if err != nil {
			return nil, err
		}

	case bolo.COUNTER:
		counter := pdu.(*bolo.CounterPDU)
		fields := map[string]interface{}{
			"Value": counter.Value,
		}
		point, err = influx.NewPoint(counter.Name, tags, fields, counter.Timestamp)
		if err != nil {
			return nil, err
		}
	default:
		log.Debugf("Ignoring PDU type '%s' - not a metric", pdu.Type())
	}
	return point, nil
}
