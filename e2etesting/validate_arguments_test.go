package e2etesting

import (
	"fmt"
	"os"
)

func validateEnvironments() error {
	requiredEnvs := []string{
		"AWS_REGION",
		"AWS_ACCESS_KEY_ID",
	}
	for _, envName := range requiredEnvs {
		if os.Getenv(envName) == "" {
			return fmt.Errorf("Missing %s environments", envName)
		}
	}
	return nil
}
