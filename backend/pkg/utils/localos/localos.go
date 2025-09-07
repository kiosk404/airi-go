package localos

import (
	"fmt"
	"os"
)

func GetLocalOSHost() string {
	protocol := os.Getenv("OSS_PROTOCOL")
	domain := os.Getenv("OSS_DOMAIN")
	port := os.Getenv("OSS_PORT")
	if port == "" {
		return fmt.Sprintf("%s://%s", protocol, domain)
	}
	return fmt.Sprintf("%s:%s", domain, port)
}
