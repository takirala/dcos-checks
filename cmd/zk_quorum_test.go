package cmd

import "testing"

// TestZookeeperQuorumCheckGetHealthURL validates that parameters Role and ForceTLS
// return the expected URL based on adminrouter configuration.
func TestZookeeperQuorumCheckGetHealthURL(t *testing.T) {
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
			expected: "http://127.0.0.1/exhibitor/v1/cluster/state/127.0.0.1",
		},
		{
			role:     "master",
			forceTLS: true,
			expected: "https://127.0.0.1:443/exhibitor/v1/cluster/state/127.0.0.1",
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
