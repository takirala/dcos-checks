package cmd

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestZookeeperQuorumCheckGetURL validates that parameters Role and ForceTLS
// return the expected URL based on adminrouter configuration.
func TestZookeeperQuorumCheckGetURL(t *testing.T) {
	zk := &ZkQuorumCheck{
		Name: "TEST",
	}

	for _, testCase := range []struct {
		role     string
		forceTLS bool
		expected string
	}{
		{
			role:     "master",
			forceTLS: false,
			expected: "http://127.0.0.1:8181/exhibitor/v1/cluster/state/127.0.0.1",
		},
		{
			role:     "master",
			forceTLS: true,
			expected: "https://127.0.0.1:8181/exhibitor/v1/cluster/state/127.0.0.1",
		},
	} {
		mockCLICfg := &CLIConfigFlags{
			NodeIPStr: "127.0.0.1",
			Role:      testCase.role,
			ForceTLS:  testCase.forceTLS,
		}

		url, err := zk.getURL(nil, mockCLICfg)
		if err != nil {
			t.Fatalf("Error running getURL: %s", err)
		}

		if url.String() != testCase.expected {
			t.Fatalf("Expect %s. Got %s", testCase.expected, url)
		}
	}
}

// TestZookeeperQuorumCheckRun
func TestZookeeperQuorumCheckRun(t *testing.T) {

	// curl 10.0.7.39:8181/exhibitor/v1/cluster/state/10.0.7.39
	//{"response":{"switches":{"restarts":true,"cleanup":true,"backups":true},"state":3,"description":"serving","isLeader":true},"errorMessage":"","success":true}
	for _, testCase := range []struct {
		role     string
		forceTLS bool
		expected string
		response string
	}{
		{
			role:     "master",
			forceTLS: false,
			expected: "http://127.0.0.1/exhibitor/v1/cluster/state/127.0.0.1",
			response: `{"switches":{"restarts":true,"cleanup":true,"backups":true},"state":3,"description":"serving","isLeader":true},"errorMessage":"","success":true}`,
		},
	} {
		mockCLICfg := &CLIConfigFlags{
			NodeIPStr: "127.0.0.1",
			Role:      testCase.role,
			ForceTLS:  testCase.forceTLS,
		}

		masterHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, testCase.response)
		})

		masterServer := httptest.NewServer(masterHandler)
		defer masterServer.Close()

		zk := &ZkQuorumCheck{
			Name: "TEST",
			urlFunc: func(client *http.Client, cfg *CLIConfigFlags) (*url.URL, error) {
				return url.Parse(masterServer.URL)
			},
		}

		_, status, err := zk.Run(context.TODO(), mockCLICfg)
		if err != nil {
			t.Fatalf("Error running the check %s, got status %d", err, status)
		}
	}
}
