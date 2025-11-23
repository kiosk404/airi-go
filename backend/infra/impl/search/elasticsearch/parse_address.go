package elasticsearch

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

func parseClusterEndpoints(address string) ([]string, error) {
	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("endpoints environment variable is required")
	}

	endpoints := strings.Split(address, ",")
	var validEndpoints []string
	uniqueEndpoints := make(map[string]bool, len(endpoints))

	for _, endpoint := range endpoints {
		trimmed := strings.TrimSpace(endpoint)
		if trimmed == "" {
			continue
		}
		if !uniqueEndpoints[trimmed] {
			uniqueEndpoints[trimmed] = true
			validEndpoints = append(validEndpoints, trimmed)
		}
	}

	if len(validEndpoints) == 0 {
		return nil, fmt.Errorf("no valid  endpoints found in: %s", address)
	}

	return validEndpoints, nil
}

func getEnvDefaultIntSetting(envVar, defaultValue string) string {
	value := os.Getenv(envVar)
	if value == "" {
		return defaultValue
	}
	if num, err := strconv.Atoi(value); err != nil || num <= 0 {
		logs.Warn("Invalid %s value: %s, using default: %s", envVar, value, defaultValue)
		return defaultValue
	}
	return value
}
