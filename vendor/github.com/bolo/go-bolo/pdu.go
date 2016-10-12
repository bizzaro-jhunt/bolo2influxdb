package bolo

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	SAMPLE           = "SAMPLE"
	RATE             = "RATE"
	COUNTER          = "COUNTER"
	SET_KEYS         = "SET.KEYS"
	STATE            = "STATE"
	STATE_TRANSITION = "TRANSITION"
	EVENT            = "EVENT"
)

type PDU interface {
	Type() string
	String() string
}

type SamplePDU struct {
	Timestamp  time.Time
	Name       string
	SampleSize int
	Min        float64
	Max        float64
	Sum        float64
	Mean       float64
	Variance   float64
	Tags       map[string]string
}

func (p *SamplePDU) Type() string {
	return SAMPLE
}
func (p *SamplePDU) String() string {
	return fmt.Sprintf("%s: Samples: %d, Min: %f, Max: %f, Sum: %f, Mean: %f, Variance: %f", p.Name, p.SampleSize, p.Min, p.Max, p.Sum, p.Mean, p.Variance)
}

type RatePDU struct {
	Timestamp time.Time
	Window    int
	Value     float64
	Name      string
	Tags      map[string]string
}

func (p *RatePDU) Type() string {
	return RATE
}
func (p *RatePDU) String() string {
	return fmt.Sprintf("%s: %f (%d sec window)", p.Name, p.Value, p.Window)
}

type CounterPDU struct {
	Timestamp time.Time
	Value     float64
	Name      string
	Tags      map[string]string
}

func (p *CounterPDU) Type() string {
	return COUNTER
}
func (p *CounterPDU) String() string {
	return fmt.Sprintf("%s: %f", p.Name, p.Value)
}

type SetKeysPDU struct {
	Keys map[string]string
}

func (p *SetKeysPDU) Type() string {
	return SET_KEYS
}
func (p *SetKeysPDU) String() string {
	var keys []string
	for k, _ := range p.Keys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var output []string
	for _, k := range keys {
		output = append(output, fmt.Sprintf("%s: %s", k, p.Keys[k]))
	}
	return strings.Join(output, "\n")
}

type StatePDU struct {
	Timestamp time.Time
	Name      string
	Stale     int
	State     int
	Summary   string
}

func (p *StatePDU) Type() string {
	return STATE
}
func (p *StatePDU) String() string {
	return fmt.Sprintf("%s: %d %s", p.Name, p.State, p.Summary)
}

type StateTransitionPDU struct {
	Timestamp time.Time
	Name      string
	Stale     int
	State     int
	Summary   string
}

func (p *StateTransitionPDU) Type() string {
	return STATE_TRANSITION
}
func (p *StateTransitionPDU) String() string {
	return fmt.Sprintf("%s: %d %s", p.Name, p.State, p.Summary)
}

type EventPDU struct {
	Timestamp time.Time
	Name      string
	Event     string
}

func (p *EventPDU) Type() string {
	return EVENT
}
func (p *EventPDU) String() string {
	return fmt.Sprintf("%s: %s", p.Name, p.Event)
}

func ParsePDU(message []string) (PDU, error) {
	if len(message) == 0 {
		return nil, fmt.Errorf("Malformed PDU: %v", message)
	}
	// Strip null-byte off strings
	for i, s := range message {
		message[i] = strings.Trim(s, "\x00")
	}

	pduType := message[0]
	switch pduType {
	case SAMPLE:
		return ParseSamplePDU(message)
	case RATE:
		return ParseRatePDU(message)
	case COUNTER:
		return ParseCounterPDU(message)
	case EVENT:
		return ParseEventPDU(message)
	case STATE:
		return ParseStatePDU(message)
	case STATE_TRANSITION:
		return ParseStateTransitionPDU(message)
	case SET_KEYS:
		return ParseSetKeysPDU(message)
	default:
		return nil, fmt.Errorf("Invalid PDU type '%s' detected", message[0])
	}
}

func parseTime(s string) (time.Time, error) {
	t, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return time.Time{}, fmt.Errorf("Invalid timestamp '%s': '%s'", s, err)
	}
	ts := time.Unix(t, 0)
	return ts, nil
}

func ParseSamplePDU(message []string) (*SamplePDU, error) {
	if len(message) < 9 {
		return nil, fmt.Errorf("Malformed %s PDU: %v", message[0], message)
	}
	ts, err := parseTime(message[1])
	if err != nil {
		return nil, err
	}

	name := message[2]

	n, err := strconv.Atoi(message[3])
	if err != nil {
		return nil, fmt.Errorf("Invalid sample size '%s' for %s PDU: %s", message[3], message[0], err)
	}

	min, err := strconv.ParseFloat(fmt.Sprintf("%s", message[4]), 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid MIN '%s' for %s PDU: '%s'", message[4], message[0], err)
	}

	max, err := strconv.ParseFloat(fmt.Sprintf("%s", message[5]), 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid MAX '%s' for %s PDU: '%s'", message[5], message[0], err)
	}

	sum, err := strconv.ParseFloat(fmt.Sprintf("%s", message[6]), 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid SUM '%s' for %s PDU: '%s'", message[6], message[0], err)
	}

	mean, err := strconv.ParseFloat(fmt.Sprintf("%s", message[7]), 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid MEAN '%s' for %s PDU: '%s'", message[7], message[0], err)
	}

	variance, err := strconv.ParseFloat(fmt.Sprintf("%s", message[8]), 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid VARIANCE '%s' for %s PDU: '%s'", message[8], message[0], err)
	}

	return &SamplePDU{
		Timestamp:  ts,
		Name:       name,
		SampleSize: n,
		Min:        min,
		Max:        max,
		Sum:        sum,
		Mean:       mean,
		Variance:   variance,
		Tags:       nil,
	}, nil
}

func ParseRatePDU(message []string) (*RatePDU, error) {
	if len(message) < 5 {
		return nil, fmt.Errorf("Malformed %s PDU: %v", message[0], message)
	}
	ts, err := parseTime(message[1])
	if err != nil {
		return nil, err
	}

	name := message[2]

	window, err := strconv.Atoi(message[3])
	if err != nil {
		return nil, fmt.Errorf("Invalid rate window '%s' for %s PDU: %s", message[3], message[0], err)
	}

	val, err := strconv.ParseFloat(fmt.Sprintf("%s", message[4]), 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid value '%s' for %s PDU: '%s'", message[4], message[0], err)
	}

	return &RatePDU{
		Timestamp: ts,
		Name:      name,
		Window:    window,
		Value:     val,
		Tags:      nil,
	}, nil
}

func ParseCounterPDU(message []string) (*CounterPDU, error) {
	if len(message) < 4 {
		return nil, fmt.Errorf("Malformed %s PDU: %v", message[0], message)
	}
	ts, err := parseTime(message[1])
	if err != nil {
		return nil, err
	}

	name := message[2]

	val, err := strconv.ParseFloat(fmt.Sprintf("%s", message[3]), 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid value '%s' for %s PDU: '%s'", message[3], message[0], err)
	}

	return &CounterPDU{
		Timestamp: ts,
		Name:      name,
		Value:     val,
		Tags:      nil,
	}, nil
}

func ParseEventPDU(message []string) (*EventPDU, error) {
	if len(message) < 4 {
		return nil, fmt.Errorf("Malformed %s PDU: %v", message[0], message)
	}
	ts, err := parseTime(message[1])
	if err != nil {
		return nil, err
	}

	name := message[2]
	event := message[3]

	return &EventPDU{
		Timestamp: ts,
		Name:      name,
		Event:     event,
	}, nil
}

func ParseStatePDU(message []string) (*StatePDU, error) {
	if len(message) < 6 {
		return nil, fmt.Errorf("Malformed %s PDU: %v", message[0], message)
	}
	ts, err := parseTime(message[2])
	if err != nil {
		return nil, err
	}

	name := message[1]

	stale, err := strconv.Atoi(fmt.Sprintf("%s", message[3]))
	if err != nil {
		return nil, fmt.Errorf("Invalid staleness '%s' for %s PDU: '%s'", message[3], message[0], err)
	}

	state, err := strconv.Atoi(fmt.Sprintf("%s", message[4]))
	if err != nil {
		return nil, fmt.Errorf("Invalid state code '%s' for %s PDU: '%s'", message[4], message[0], err)
	}

	summary := message[5]

	return &StatePDU{
		Timestamp: ts,
		Name:      name,
		Stale:     stale,
		State:     state,
		Summary:   summary,
	}, nil
}

func ParseStateTransitionPDU(message []string) (*StateTransitionPDU, error) {
	if len(message) < 6 {
		return nil, fmt.Errorf("Malformed %s PDU: %v", message[0], message)
	}
	ts, err := parseTime(message[2])
	if err != nil {
		return nil, err
	}

	name := message[1]

	stale, err := strconv.Atoi(fmt.Sprintf("%s", message[3]))
	if err != nil {
		return nil, fmt.Errorf("Invalid staleness '%s' for %s PDU: '%s'", message[3], message[0], err)
	}

	state, err := strconv.Atoi(fmt.Sprintf("%s", message[4]))
	if err != nil {
		return nil, fmt.Errorf("Invalid state code '%s' for %s PDU: '%s'", message[4], message[0], err)
	}

	summary := message[5]

	return &StateTransitionPDU{
		Timestamp: ts,
		Name:      name,
		Stale:     stale,
		State:     state,
		Summary:   summary,
	}, nil
}

func ParseSetKeysPDU(message []string) (*SetKeysPDU, error) {
	msgs := message[1:]
	if len(msgs)%2 != 0 {
		return nil, fmt.Errorf("Error parsing SET.KEYS PDU: Not enough Values for all the Keys")
	}
	keys := map[string]string{}
	for i := 0; i < len(msgs); i += 2 {
		keys[msgs[i]] = msgs[i+1]
	}
	return &SetKeysPDU{
		Keys: keys,
	}, nil
}
