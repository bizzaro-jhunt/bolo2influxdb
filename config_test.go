package main_test

import (
	. "github.com/bolo/bolo2influxdb"
	"github.com/starkandwayne/metrics/influxdb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configs", func() {
	Context("Load()", func() {
		It("Should return an error when the file cannot be loaded", func() {
			cfg, err := LoadConfig("assets/nonexistent.cfg")
			Expect(cfg).Should(BeNil())
			Expect(err).ShouldNot(BeNil())
		})
		It("Should return an error when the file isn't parseable JSON", func() {
			cfg, err := LoadConfig("assets/invalid.cfg")
			Expect(cfg).Should(BeNil())
			Expect(err).ShouldNot(BeNil())
		})
		It("Returns a parsed JSON config into a Config struct", func() {
			cfg, err := LoadConfig("assets/valid.cfg")
			Expect(err).Should(BeNil())
			Expect(cfg).Should(Equal(&Config{
				Bolo: BoloConfig{
					Addr: "10.10.10.10",
					Port: "2997",
				},
				Influx: influxdb.Config{
					Addr:     "http://10.10.10.11:8086",
					User:     "iuser",
					Password: "ipass",
					Database: "influx",
				},
			}))
		})
	})
})
