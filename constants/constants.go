package constants

// NB: use the zero value (unknown) as the first const in an iota expression
// to convey the right "unknown" meaning until set to some other value.
const (
	// StatusUnknown means that the status was not able to be determined
	StatusUnknown = iota

	// StatusOK means that the check succeeded
	StatusOK

	// StatusWarning means that the check generated a warning
	StatusWarning

	// StatusFailure means that the check failed
	StatusFailure
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
