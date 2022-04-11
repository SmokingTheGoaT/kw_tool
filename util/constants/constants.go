package constants

const (
	/*----------Helpers------------------*/
	EmptyString                  string = ""
	DefaultSearchCountPerRequest int    = 100
	RedisAddr                    string = "127.0.0.1:6379"
	ServerAddr                   string = "127.0.0.1:8443"

	/*----------Headers------------------*/
	HeaderSessionID string = "X-Session-ID"
	HeaderPlatform  string = "X-Platform"

	/*----------URIs------------------*/
	SuggestURI string = "/suggest"
	HealthzURI string = "/healthz"
)

var (
	Platforms = []string{"google", "amazon"}
)
