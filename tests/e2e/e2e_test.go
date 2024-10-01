package e2e

import (
	"bytes"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Adapter API", func() {
	It("should handle a request and return a response", func() {
		client := &http.Client{Timeout: 10 * time.Second}
		requestBody := []byte(`{"key": "value"}`)
		req, err := http.NewRequest("POST", "http://localhost:3000/api/adapter", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

	})
})
