package components

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dcos/dcos-checks/common"
)

// TestComponentCheckGetHealthURL validates that parameters Role and ForceTLS
// return the expected URL based on adminrouter / 3dt configuration.
func TestComponentCheckGetHealthURL(t *testing.T) {
	c := &componentCheck{
		Name: "TEST",
	}

	for _, item := range []struct {
		role     string
		scheme   string
		port     int
		expected string
	}{
		{
			role:     "master",
			scheme:   "http",
			port:     1050,
			expected: "http://127.0.0.1:1050/",
		},
		{
			role:     "master",
			scheme:   "https",
			port:     443,
			expected: "https://127.0.0.1:443/",
		},
		{
			role:     "agent",
			scheme:   "http",
			port:     61001,
			expected: "http://127.0.0.1:61001/",
		},
		{
			role:     "agent",
			scheme:   "https",
			port:     61002,
			expected: "https://127.0.0.1:61002/",
		},
		{
			role:     "agent_public",
			scheme:   "http",
			port:     61001,
			expected: "http://127.0.0.1:61001/",
		},
		{
			role:     "agent_public",
			scheme:   "https",
			port:     61002,
			expected: "https://127.0.0.1:61002/",
		},
	} {
		mockCLICfg := &common.CLIConfigFlags{
			NodeIPStr: "127.0.0.1",
			Role:      item.role,
		}

		url, err := c.getHealthURL(nil, "/", item.scheme, item.port, mockCLICfg)
		if err != nil {
			t.Fatalf("Error running getHealthURL: %s", err)
		}

		if url.String() != item.expected {
			t.Fatalf("Expect %s. Got %s", item.expected, url.String())
		}
	}
}

func TestDiagnosticsResponse(t *testing.T) {
	// A sample from the output of system/health/v1 with checks marked as not healthy
	response := `{"units":[{"id":"dcos-adminrouter.service","health":0,"output":"","description":"exposes a unified control plane proxy for components and services using NGINX","help":"","name":"Admin Router Master"},{"id":"dcos-checks-poststart.service","health":1,"output":"","description":"Run node-poststart checks","help":"","name":"DC/OS Poststart Checks"},{"id":"dcos-checks-poststart.timer","health":1,"output":"","description":"timer for DC/OS Checks service","help":"","name":"DC/OS Checks Timer"},{"id":"dcos-cosmos.service","health":0,"output":"","description":"installs and manages DC/OS packages from DC/OS package repositories, such as the Mesosphere Universe","help":"","name":"DC/OS Package Manager (Cosmos)"},{"id":"dcos-diagnostics.service","health":0,"output":"","description":"aggregates and exposes component health","help":"","name":"DC/OS Diagnostics Master"},{"id":"dcos-diagnostics.socket","health":0,"output":"","description":"socket for DC/OS Diagnostics Agent","help":"","name":"DC/OS Diagnostics Agent Socket"}]}`

	var dr diagnosticsResponse
	if err := json.NewDecoder(strings.NewReader(response)).Decode(&dr); err != nil {
		t.Fatalf("Error decoding")
	}

	complist := []string{"dcos-checks-poststart.service", "dcos-checks-poststart.timer"}
	_, retCode := dr.checkHealth(complist)
	if retCode != 0 {
		t.Fatalf("Component health check failed when it should have passed")
	}

	_, retCode = dr.checkHealth(nil)
	if retCode == 0 {
		t.Fatalf("Component health check passed when it should have failed")
	}
}
