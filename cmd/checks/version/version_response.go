package version

// versionResponse responses /dcos-metadata/dcos-version.json
type versionResponse struct {
	Version         string `json:"version"`
	DcosImageCommit string `json:"dcos-image-commit"`
	BootstrapID     string `json:"bootstrap-id"`
}
