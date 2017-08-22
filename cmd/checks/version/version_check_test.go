package version

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dcos/dcos-checks/common"
	"github.com/dcos/dcos-checks/constants"
)

func TestVersionCheckUrl(t *testing.T) {
	for _, testCase := range []struct {
		urlopt   common.URLFields
		forceTLS bool
		expected string
	}{
		{
			urlopt: common.URLFields{
				Host: "127.0.0.1",
				Port: 5050,
				Path: "metrics/snapshot"},
			forceTLS: false,
			expected: "http://127.0.0.1:5050/metrics/snapshot",
		},
		{
			urlopt: common.URLFields{
				Host: "127.0.0.1",
				Port: 0,
				Path: "/v1/hosts/master.mesos"},
			forceTLS: false,
			expected: "http://127.0.0.1/v1/hosts/master.mesos",
		},
	} {
		mockCLICfg := &common.CLIConfigFlags{
			NodeIPStr: "127.0.0.1",
			Role:      "master",
			ForceTLS:  testCase.forceTLS,
		}

		url, err := common.GetURL(nil, mockCLICfg, testCase.urlopt)
		if err != nil {
			t.Fatalf("Error running getURL: %s", err)
		}

		if url.String() != testCase.expected {
			t.Fatalf("Expect %s. Got %s", testCase.expected, url)
		}

	}
}

// TestVersionCheckListofmasters gets list of masters from mesos-dns endpoint
func TestVersionCheckListOfMasters(t *testing.T) {
	for _, testCase := range []struct {
		role      string
		forceTLS  bool
		status    int
		response  string
		expStatus int
		expValue  string
	}{
		{
			role:      "master",
			forceTLS:  false,
			status:    http.StatusOK,
			response:  `[{ "host": "leader.mesos.", "ip": "10.0.4.197" }]`,
			expStatus: constants.StatusOK,
			expValue:  "10.0.4.197",
		},
	} {
		mockCLICfg := &common.CLIConfigFlags{
			NodeIPStr: "127.0.0.1",
			Role:      testCase.role,
			ForceTLS:  testCase.forceTLS,
		}

		masterHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(testCase.status)
			w.Header().Set("Content-Type", "application/json")
			if r.URL.Path == "/v1/hosts/master.mesos" {
				io.WriteString(w, testCase.response)
			}
		})

		masterServer := httptest.NewServer(masterHandler)
		defer masterServer.Close()

		test := &versionCheck{
			Name:          "TEST",
			ClusterLeader: "127.0.0.1",
		}

		testurl, err := url.Parse(masterServer.URL)
		if err != nil {
			t.Fatalf("could not parse")
		}
		var masterurlopt common.URLFields
		masterurlopt.Host = testurl.Host
		masterurlopt.Port = 0
		masterurlopt.Path = "/v1/hosts/master.mesos"

		masters, err := test.ListOfMasters(mockCLICfg, masterurlopt)

		if err != nil {
			t.Fatalf("Status %s", err)
		}
		if masters[0] != testCase.expValue {
			t.Fatalf("Getting random nonsense, not the correct value")
		}
	}
}

// TestVersionCheckListofAgents gets list of agents by parsing mesos endpoint /slaves
func TestVersionCheckListOfAgents(t *testing.T) {
	for _, testCase := range []struct {
		role      string
		forceTLS  bool
		status    int
		response  string
		expStatus int
		expValue  string
		version   string
	}{
		{
			role:      "master",
			forceTLS:  false,
			status:    http.StatusOK,
			response:  `{"slaves":[{"id":"529c3971-b5bb-4f9e-b817-bb32def0ede2-S1","hostname":"10.0.6.233","port":5051,"attributes":{"public_ip":"true"},"pid":"slave(1)@10.0.6.233:5051","registered_time":1496728541.24296,"resources":{"disk":35566.0,"mem":14021.0,"gpus":0.0,"cpus":4.0,"ports":"[1-21, 23-5050, 5052-32000]"},"used_resources":{"disk":0.0,"mem":0.0,"gpus":0.0,"cpus":0.0},"offered_resources":{"disk":0.0,"mem":0.0,"gpus":0.0,"cpus":0.0},"reserved_resources":{"slave_public":{"disk":35566.0,"mem":14021.0,"gpus":0.0,"cpus":4.0,"ports":"[1-21, 23-5050, 5052-32000]"}},"unreserved_resources":{"disk":0.0,"mem":0.0,"gpus":0.0,"cpus":0.0},"active":true,"version":"1.3.0","capabilities":["MULTI_ROLE"],"reserved_resources_full":{"slave_public":[{"name":"ports","type":"RANGES","ranges":{"range":[{"begin":1,"end":21},{"begin":23,"end":5050},{"begin":5052,"end":32000}]},"role":"slave_public"},{"name":"disk","type":"SCALAR","scalar":{"value":35566.0},"role":"slave_public"},{"name":"cpus","type":"SCALAR","scalar":{"value":4.0},"role":"slave_public"},{"name":"mem","type":"SCALAR","scalar":{"value":14021.0},"role":"slave_public"}]},"used_resources_full":[],"offered_resources_full":[]}]}`,
			expStatus: constants.StatusOK,
			expValue:  "10.0.6.233",
			version:   `{"version": "1.10-dev", "dcos-image-commit": "ccb53df0da261508249570df577c47bbbcc09f82", "bootstrap-id": "8468e43583e21ccb482ff303ed7496f84bbadb4d"}`,
		},
	} {
		mockCLICfg := &common.CLIConfigFlags{
			NodeIPStr: "127.0.0.1",
			Role:      testCase.role,
			ForceTLS:  testCase.forceTLS,
		}

		agentHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(testCase.status)
			w.Header().Set("Content-Type", "application/json")
			if r.URL.Path == "/slaves" {
				io.WriteString(w, testCase.response)
			}
			if r.URL.Path == "/dcos-metadata/dcos-version.json" {
				io.WriteString(w, testCase.version)
			}
		})

		agentServer := httptest.NewServer(agentHandler)
		defer agentServer.Close()

		test := &versionCheck{
			Name:          "TEST",
			ClusterLeader: "127.0.0.1",
		}

		testurl, err := url.Parse(agentServer.URL)
		if err != nil {
			t.Fatalf("could not parse")
		}
		var agenturlopt common.URLFields
		agenturlopt.Host = testurl.Host
		agenturlopt.Port = 0
		agenturlopt.Path = "/slaves"

		agents, err := test.ListOfAgents(mockCLICfg, agenturlopt)

		if err != nil {
			t.Fatalf("Status %s", err)
		}
		if agents[0] != testCase.expValue {
			t.Fatalf("Getting random nonsense, not the correct value")
		}

		var versionurlopt common.URLFields
		versionurlopt.Host = testurl.Host
		versionurlopt.Port = 0
		versionurlopt.Path = "/dcos-metadata/dcos-version.json"

		version, err := test.GetVersion(mockCLICfg, versionurlopt)

		if err != nil {
			t.Fatalf("Status %s", err)
		}
		if version != "1.10-dev" {
			t.Fatalf("Getting nonsense %s, not the correct value", version)
		}

	}
}
