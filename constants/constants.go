package constants

// NB: it's important to keep StatusOK as the first constant here because
// success=0.  Then warning, failure, and then unknown.
const (
	// StatusOK means that the check succeeded
	StatusOK = iota

	// StatusWarning means that the check generated a warning
	StatusWarning

	// StatusFailure means that the check failed
	StatusFailure

	// StatusUnknown means that the status was not able to be determined
	StatusUnknown
)

const (
	// AdminrouterMasterHTTPSPort is the https port on which adminrouter runs
	AdminrouterMasterHTTPSPort = 443

	// AdminrouterAgentHTTPPort is the port on which adminrouter on the agent
	// listens for HTTP requests
	AdminrouterAgentHTTPPort = 61001

	// AdminrouterAgentHTTPSPort is the port on which adminrouter on the agent
	// listens for HTTPS requests
	AdminrouterAgentHTTPSPort = 61002

	// ExhibitorHTTPPort is the port on which Exhibitor listens on a master
	// node
	ExhibitorHTTPPort = 8181

	// MesosMasterHTTPPort is the port on which the Mesos master listens for
	// HTTP requests.
	MesosMasterHTTPPort = 5050

	// MesosAgentHTTPPort is the port on which the Mesos agent listens for
	// HTTP requests.
	MesosAgentHTTPPort = 5051

	// MesosDNSPort is port on which Mesos DNS listens
	MesosDNSPort = 8123

	// HTTPScheme is the default non-secure http protocol method
	HTTPScheme = "http"

	// HTTPSScheme is the secure http protocol method
	HTTPSScheme = "https"
)
