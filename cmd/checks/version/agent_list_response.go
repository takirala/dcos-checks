package version

// agentListResponse response for /slaves
type agentListResponse struct {
	Slaves []struct {
		ID         string `json:"id"`
		Hostname   string `json:"hostname"`
		Port       int    `json:"port"`
		Attributes struct {
			PublicIP string `json:"public_ip"`
		} `json:"attributes"`
	} `json:"slaves"`
	RecoveredSlaves []interface{} `json:"recovered_slaves"`
}
