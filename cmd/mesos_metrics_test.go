package cmd

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestMesosMetricsCheckUrl verifies we get the right url
func TestMesosMetricsCheckUrl(t *testing.T) {
	test := &MesosMetricsCheck{
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
			expected: "http://127.0.0.1:5050/metrics/snapshot",
		},
		{
			role:     "agent_public",
			forceTLS: false,
			expected: "http://127.0.0.1:5051/metrics/snapshot",
		},
		{
			role:     "agent",
			forceTLS: true,
			expected: "https://127.0.0.1:5051/metrics/snapshot",
		},
	} {
		mockCLICfg := &CLIConfigFlags{
			NodeIPStr: "127.0.0.1",
			Role:      testCase.role,
			ForceTLS:  testCase.forceTLS,
		}

		url, err := test.getURL(nil, mockCLICfg)
		if err != nil {
			t.Fatalf("Error running getURL: %s", err)
		}

		if url.String() != testCase.expected {
			t.Fatalf("Expect %s. Got %s", testCase.expected, url)
		}
	}
}

// TestMesosMetricsCheckRun checks run
func TestMesosMetricsCheckRun(t *testing.T) {
	for _, testCase := range []struct {
		role      string
		forceTLS  bool
		status    int
		response  string
		expStatus int
	}{
		{
			role:      "master",
			forceTLS:  false,
			status:    http.StatusOK,
			response:  `{"slave\/tasks_finished":0.0,"slave\/cpus_total":4.0,"slave\/executors_preempted":0.0,"slave\/registered":1.0,"registrar\/log\/recovered": 1.0}`,
			expStatus: statusOK,
		},
		{
			role:      "agent",
			forceTLS:  false,
			status:    http.StatusOK,
			response:  `{"slave\/tasks_finished":0.0,"slave\/cpus_total":4.0,"slave\/executors_preempted":0.0,"slave\/registered":1.0,"registrar\/log\/recovered": 1.0}`,
			expStatus: statusOK,
		},
		{
			role:      "agent",
			forceTLS:  false,
			status:    http.StatusOK,
			response:  `{"slave\/tasks_finished":0.0,"slave\/cpus_total":4.0,"slave\/executors_preempted":0.0,"slave\/registered":0.0,"registrar\/log\/recovered": 1.0}`,
			expStatus: statusFailure,
		},
	} {
		mockCLICfg := &CLIConfigFlags{
			NodeIPStr: "127.0.0.1",
			Role:      testCase.role,
			ForceTLS:  testCase.forceTLS,
		}

		masterHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(testCase.status)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, testCase.response)
		})

		masterServer := httptest.NewServer(masterHandler)
		defer masterServer.Close()

		test := &MesosMetricsCheck{
			Name: "TEST",
			urlFunc: func(client *http.Client, cfg *CLIConfigFlags) (*url.URL, error) {
				return url.Parse(masterServer.URL)
			},
		}
		_, status, err := test.Run(context.TODO(), mockCLICfg)
		if status != testCase.expStatus {
			t.Fatalf("Status %d not as expected status %d %s", status, testCase.expStatus, err)
		}
	}
}
