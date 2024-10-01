package e2e

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var serverCmd *exec.Cmd

var _ = BeforeSuite(func() {
	// Start the NATS server using Docker Compose
	fmt.Println("Starting the NATS server")
	err := exec.Command("docker-compose", "up", "-d").Run()
	Expect(err).NotTo(HaveOccurred())

	// Wait for the NATS server to start
	time.Sleep(5 * time.Second)

	// Start the API server
	fmt.Println("Starting the API server")
	err = exec.Command("make", "build").Run()
	Expect(err).NotTo(HaveOccurred())

	// get output from make run
	output, err := exec.Command("make", "run", "deployment=local").Output()
	Expect(err).NotTo(HaveOccurred())
	fmt.Println(string(output))

	// Wait for the server to start
	time.Sleep(5 * time.Second)
	fmt.Println("API server started")
	fmt.Println(serverCmd.Output())
})

var _ = AfterSuite(func() {
	// Stop the API server
	if serverCmd != nil && serverCmd.Process != nil {
		serverCmd.Process.Kill()
	}

	// Stop the NATS server using Docker Compose
	fmt.Println("Stopping the NATS server")
	err := exec.Command("docker-compose", "down").Run()
	Expect(err).NotTo(HaveOccurred())
})

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
