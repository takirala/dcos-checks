package version

// masterListResponses response for leader.mesos/master.mesos
type masterListResponses []struct {
	Host string `json:"host"`
	IP   string `json:"ip"`
}
