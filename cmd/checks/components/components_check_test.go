package components

import (
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
