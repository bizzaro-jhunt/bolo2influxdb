package bolo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/starkandwayne/metrics/go-bolo"
)

// FIXME: mock bolo? how to mock zmq failures?
var _ = Describe("Bolo Connections", func() {
	Context("Connect()", func() {
		It("Should return an error when a zmq context creation fails", func() {
			Skip("Figure out how to test this properly")
			pduChan, errChan, err := Connect("")
			Expect(pduChan).Should(BeNil())
			Expect(errChan).Should(BeNil())
			Expect(err).ShouldNot(BeNil())
		})
		It("Should return an error when a zmq socket creation fails", func() {
			Skip("Figure out how to test this properly")
			pduChan, errChan, err := Connect("")
			Expect(pduChan).Should(BeNil())
			Expect(errChan).Should(BeNil())
			Expect(err).ShouldNot(BeNil())
		})
		It("Should return an error when the zmq socket can't connect", func() {
			Skip("Figure out how to test this properly")
			pduChan, errChan, err := Connect("")
			Expect(pduChan).Should(BeNil())
			Expect(errChan).Should(BeNil())
			Expect(err).ShouldNot(BeNil())
		})
		It("Should return an error when the zmq socket can't subscribe", func() {
			Skip("Figure out how to test this properly")
			pduChan, errChan, err := Connect("")
			Expect(pduChan).Should(BeNil())
			Expect(errChan).Should(BeNil())
			Expect(err).ShouldNot(BeNil())
		})
		It("Returns opened channels for pdu + error handling", func() {
			Skip("Figure out how to test this properly")
		})
		It("Sends socket read failures to the error channel", func() {
			Skip("Figure out how to test this properly")
		})
		It("Sends PDU parse failures to the error channel", func() {
			Skip("Figure out how to test this properly")
		})
		It("Sends parsed PDUs to the pdu channel", func() {
			Skip("Figure out how to test this properly")
		})
	})
})
