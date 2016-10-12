package bolo_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/starkandwayne/metrics/go-bolo"
)

type TestPDU struct {
	TestCase string
	Payload  []string
	Expect   PDU
}

var _ = Describe("Bolo PDUs", func() {
	Context("SamplePDUs", func() {
		Context("Type()", func() {
			It("Returns the correct type", func() {
				pdu := &SamplePDU{}
				Expect(pdu.Type()).Should(Equal(SAMPLE))
			})
		})
		Context("String()", func() {
			It("Returns the desired formatting string", func() {
				pdu := &SamplePDU{
					Timestamp:  time.Now(),
					Name:       "testSample",
					SampleSize: 1,
					Min:        2.0,
					Max:        3.0,
					Sum:        4.0,
					Mean:       5.0,
					Variance:   6.0,
				}
				Expect(pdu.String()).Should(Equal("testSample: Samples: 1, Min: 2.000000, Max: 3.000000, Sum: 4.000000, Mean: 5.000000, Variance: 6.000000"))
			})
		})
	})
	Context("RatePDUs", func() {
		Context("Type()", func() {
			It("Returns the correct type", func() {
				pdu := &RatePDU{}
				Expect(pdu.Type()).Should(Equal(RATE))
			})
		})
		Context("String()", func() {
			It("Returns the desired formatting string", func() {
				pdu := &RatePDU{
					Name:   "testRate",
					Window: 1,
					Value:  2.0,
				}
				Expect(pdu.String()).Should(Equal("testRate: 2.000000 (1 sec window)"))
			})
		})
	})
	Context("CounterPDUs", func() {
		Context("Type()", func() {
			It("Returns the correct type", func() {
				pdu := &CounterPDU{}
				Expect(pdu.Type()).Should(Equal(COUNTER))
			})
		})
		Context("String()", func() {
			It("Returns the desired formatting string", func() {
				pdu := &CounterPDU{
					Name:  "testCounter",
					Value: 1.0,
				}
				Expect(pdu.String()).Should(Equal("testCounter: 1.000000"))
			})
		})
	})
	Context("SetKeysPDUs", func() {
		Context("Type()", func() {
			It("Returns the correct type", func() {
				pdu := &SetKeysPDU{}
				Expect(pdu.Type()).Should(Equal(SET_KEYS))
			})
		})
		Context("String()", func() {
			It("Returns the desired formatting string", func() {
				pdu := &SetKeysPDU{
					Keys: map[string]string{
						"k1": "v1",
						"k2": "v2",
					},
				}
				Expect(pdu.String()).Should(Equal("k1: v1\nk2: v2"))
			})
		})
	})
	Context("StatePDUs", func() {
		Context("Type()", func() {
			It("Returns the correct type", func() {
				pdu := &StatePDU{}
				Expect(pdu.Type()).Should(Equal(STATE))
			})
		})
		Context("String()", func() {
			It("Returns the desired formatting string", func() {
				pdu := &StatePDU{
					Name:    "testState",
					Stale:   1,
					State:   2,
					Summary: "This is our state",
				}
				Expect(pdu.String()).Should(Equal("testState: 2 This is our state"))
			})
		})
	})
	Context("StateTransitionPDUs", func() {
		Context("Type()", func() {
			It("Returns the correct type", func() {
				pdu := &StateTransitionPDU{}
				Expect(pdu.Type()).Should(Equal(STATE_TRANSITION))
			})
		})
		Context("String()", func() {
			It("Returns the desired formatting string", func() {
				pdu := &StateTransitionPDU{
					Name:    "testStateTransition",
					Stale:   1,
					State:   2,
					Summary: "This is our state transition",
				}
				Expect(pdu.String()).Should(Equal("testStateTransition: 2 This is our state transition"))
			})
		})
	})
	Context("EventPDUs", func() {
		Context("Type()", func() {
			It("Returns the correct type", func() {
				pdu := &EventPDU{}
				Expect(pdu.Type()).Should(Equal(EVENT))
			})
		})
		Context("String()", func() {
			It("Returns the desired formatting string", func() {
				pdu := &EventPDU{
					Name:  "testEvent",
					Event: "The test event occurred",
				}
				Expect(pdu.String()).Should(Equal("testEvent: The test event occurred"))
			})
		})
	})

	Context("ParsePDUs", func() {
		testPDUs := []TestPDU{
			/*
				Bolo Spec:
				EVENT                                    ; broadcast when a [NEW.EVENT] PDU
				<TS>                                     ; is received.
				<NAME>
				<EXTRA>
			*/
			TestPDU{
				TestCase: "returns an error if the time is not parseable",
				Payload:  []string{"EVENT", "invalid timestamp", "testEvent", "This is the event"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an EventPDU with the correct data set",
				Payload:  []string{"EVENT", "0", "testEvent", "This is the event"},
				Expect: &EventPDU{
					Timestamp: time.Unix(0, 0),
					Name:      "testEvent",
					Event:     "This is the event",
				},
			},
			/*
				Bolo Spec:
				TRANSITION                               ; like [STATE], except the state is
				<NAME>                                   ; only broadcast when the status
				<TS>                                     ; changes significantly (i.e. from
				<STALE>                                  ; OK to WARNING, or CRITICAL to OK)
				<CODE>
				<SUMMARY>
			*/
			TestPDU{
				TestCase: "returns an error if the time is not parseable",
				Payload:  []string{"TRANSITION", "testTransition", "invalid timestamp", "1", "2", "State message"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the staleness is not parseable",
				Payload:  []string{"TRANSITION", "testTransition", "0", "invalid stale", "2", "State message"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the state is not parseable",
				Payload:  []string{"TRANSITION", "testTransition", "0", "1", "invalid state", "State message"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns a StateTransitionPDU with the correct data set",
				Payload:  []string{"TRANSITION", "testTransition", "0", "1", "2", "State message"},
				Expect: &StateTransitionPDU{
					Name:      "testTransition",
					Timestamp: time.Unix(0, 0),
					Stale:     1,
					State:     2,
					Summary:   "State message",
				},
			},
			/*
				Bolo Spec:
				STATE                                    ; broadcast in response to *every*
				<NAME>                                   ; [PUT.STATE] PDU, to inform subs
				<TS>                                     ; of an update to state (even if the
				<STALE>                                  ; status didn't change materially).
				<CODE>
				<SUMMARY>
			*/
			TestPDU{
				TestCase: "returns an error if the time is not parseable",
				Payload:  []string{"STATE", "testState", "invalid timestamp", "1", "2", "State message"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the staleness is not parseable",
				Payload:  []string{"STATE", "testState", "0", "invalid stale", "2", "State message"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the state is not parseable",
				Payload:  []string{"STATE", "testState", "0", "1", "invalid state", "State message"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns a StateTransitionPDU with the correct data set",
				Payload:  []string{"STATE", "testState", "0", "1", "2", "State message"},
				Expect: &StatePDU{
					Name:      "testState",
					Timestamp: time.Unix(0, 0),
					Stale:     1,
					State:     2,
					Summary:   "State message",
				},
			},
			/*
				Bolo Spec:
				SET.KEYS                                ; broadcast all known keys, in groups
				<KEY 1>                                 ; of up to 30 key/value pairs.
				<VALUE 1>
				...
				<KEY N>
				<VALUE N>
			*/
			TestPDU{
				TestCase: "returns an error if there are not an even number of messages",
				Payload:  []string{"SET.KEYS", "key1", "val1", "key2"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns a SetKeysPDU with the correct data set",
				Payload:  []string{"SET.KEYS", "key1", "val1", "key2", "val2"},
				Expect: &SetKeysPDU{
					Keys: map[string]string{
						"key1": "val1",
						"key2": "val2",
					},
				},
			},
			/*
				Bolo Spec:
				COUNTER                                  ; broadcast on window rollover.
				<TS>                                     ; subscribers can store the value
				<NAME>                                   ; of the counter (e.g. in RRDs)
				<VALUE>                                  ; for use later.
			*/
			TestPDU{
				TestCase: "returns an error if the time is not parseable",
				Payload:  []string{"COUNTER", "invalid timestamp", "testCounter", "3"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the value is not parseable",
				Payload:  []string{"COUNTER", "0", "testCounter", "invalid value"},
			}, TestPDU{
				TestCase: "returns a CounterPDU with the correct data set",
				Payload:  []string{"COUNTER", "0", "testCounter", "1"},
				Expect: &CounterPDU{
					Name:      "testCounter",
					Timestamp: time.Unix(0, 0),
					Value:     1.0,
				},
			},
			/*
				Bolo Spec:
				RATE                                     ; broadcast on window rollover.
				<TS>                                     ; contains the calculated rate,
				<NAME>                                   ; per window interval.
				<WINDOW>
				<VALUE>
			*/
			TestPDU{
				TestCase: "returns an error if the time is not parseable",
				Payload:  []string{"RATE", "invalid timestamp", "testRate", "60", "1"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the rate window is not parseable",
				Payload:  []string{"RATE", "0", "testRate", "bad window", "1"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the value is not parseable",
				Payload:  []string{"RATE", "0", "testRate", "60", "bad value"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns a RatePDU with the correct data set",
				Payload:  []string{"RATE", "0", "testRate", "60", "1"},
				Expect: &RatePDU{
					Name:      "testRate",
					Timestamp: time.Unix(0, 0),
					Window:    60,
					Value:     1.0,
				},
			},
			/*
				Bolo Spec:
				SAMPLE                                   ; broadcast on window rollover.
				<TS>                                     ; contains all descriptive quantities
				<NAME>                                   ; of the set.  subscribers can use
				<N>                                      ; this data to update stored data
				<MIN>                                    ; (e.g. RRDs) for use later.
				<MAX>
				<SUM>
				<MEAN>
				<VARIANCE>
			*/
			TestPDU{
				TestCase: "returns an error if the time is not parseable",
				Payload:  []string{"SAMPLE", "invalid timestamp", "testSample", "1", "2", "3", "4", "5", "6"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the sample size is not parseable",
				Payload:  []string{"SAMPLE", "0", "testSample", "invalid sample", "2", "3", "4", "5", "6"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the min is not parseable",
				Payload:  []string{"SAMPLE", "0", "testSample", "1", "invalid min", "3", "4", "5", "6"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the max is not parseable",
				Payload:  []string{"SAMPLE", "0", "testSample", "1", "2", "invalid max", "4", "5", "6"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the sum is not parseable",
				Payload:  []string{"SAMPLE", "0", "testSample", "1", "2", "3", "invalid sum", "5", "6"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the mean is not parseable",
				Payload:  []string{"SAMPLE", "0", "testSample", "1", "2", "3", "4", "invalid mean", "6"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns an error if the variance is not parseable",
				Payload:  []string{"SAMPLE", "0", "testSample", "1", "2", "3", "4", "5", "invalid variance"},
				Expect:   nil,
			}, TestPDU{
				TestCase: "returns a SamplePDU with the correct data set",
				Payload:  []string{"SAMPLE", "0", "testSample", "1", "2", "3", "4", "5", "6"},
				Expect: &SamplePDU{
					Name:       "testSample",
					Timestamp:  time.Unix(0, 0),
					SampleSize: 1,
					Min:        2.0,
					Max:        3.0,
					Sum:        4.0,
					Mean:       5.0,
					Variance:   6.0,
				},
			},
			// Test with strings with trailing null characters
			TestPDU{
				TestCase: "works when strings have trailing null characters from zmq",
				Payload:  []string{"SAMPLE\x00", "0\x00", "testSample\x00", "1\x00", "2\x00", "3\x00", "4\x00", "5\x00", "6\x00"},
				Expect: &SamplePDU{
					Name:       "testSample",
					Timestamp:  time.Unix(0, 0),
					SampleSize: 1,
					Min:        2.0,
					Max:        3.0,
					Sum:        4.0,
					Mean:       5.0,
					Variance:   6.0,
				},
			},
			// Test with empty messages
			TestPDU{
				TestCase: "gracefully handles lack of messages in the payload",
				Payload:  []string{},
				Expect:   nil,
			},
			// Test with not enough fields for a message type
			TestPDU{
				TestCase: "gracefully handles not having all the data",
				Payload:  []string{"EVENT"},
				Expect:   nil,
			},
			TestPDU{
				TestCase: "gracefully handles not having all the data",
				Payload:  []string{"STATE"},
				Expect:   nil,
			},
			TestPDU{
				TestCase: "gracefully handles not having all the data",
				Payload:  []string{"TRANSITION"},
				Expect:   nil,
			},
			TestPDU{
				TestCase: "gracefully handles not having all the data",
				Payload:  []string{"COUNTER"},
				Expect:   nil,
			},
			TestPDU{
				TestCase: "gracefully handles not having all the data",
				Payload:  []string{"RATE"},
				Expect:   nil,
			},
			TestPDU{
				TestCase: "gracefully handles not having all the data",
				Payload:  []string{"SAMPLE"},
				Expect:   nil,
			},
			TestPDU{
				TestCase: "gracefully handles not having all the data",
				Payload:  []string{"SET.KEYS"},
				Expect: &SetKeysPDU{
					Keys: map[string]string{},
				},
			},
		}

		for _, t := range testPDUs {
			var test = t
			var kind string
			if len(test.Payload) > 0 {
				kind = test.Payload[0]
			}
			It(fmt.Sprintf("of type %s %s", kind, test.TestCase), func() {
				pdu, err := ParsePDU(test.Payload)
				if test.Expect == nil {
					Expect(pdu).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				} else {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(pdu).Should(Equal(test.Expect))
				}
			})
		}
	})
})
