package config

const (
	// Localhost port number to connect for server
	PortNum = "localhost:5221"
	// Localhost port number to connect for Seller API
	APIPortNum = "localhost:5000"

	// Expiry time in seconds = 300 seconds = 5 minutes
	SessionExpireTime int = 300

	// Time to clean MapSessions 120 seconds = 2 minutes
	CleanSessionTime int = 120

	// Directory that stores self generated cerificate.
	CertPath = "./cert/"

	// Other directory Paths
	FilePath = "./"
	LogPath  = "./log/"
)
