package version

import (
	"fmt"
	"sync/atomic"

	utilversion "github.com/kiosk404/airi-go/backend/pkg/version/util"
)

var dynamicGitVersion atomic.Value

func init() {
	// initialize to static gitVersion
	dynamicGitVersion.Store(GitVersion)
}

func SetDynamicVersion(dynamicVersion string) error {
	if err := ValidateDynamicVersion(dynamicVersion); err != nil {
		return err
	}
	dynamicGitVersion.Store(dynamicVersion)
	return nil
}

func ValidateDynamicVersion(dynamicVersion string) error {
	return validateDynamicVersion(dynamicVersion, GitVersion)
}

func validateDynamicVersion(dynamicVersion, defaultVersion string) error {
	if len(dynamicVersion) == 0 {
		return fmt.Errorf("version must not be empty")
	}
	if dynamicVersion == defaultVersion {
		return nil
	}
	vRuntime, err := utilversion.ParseSemantic(dynamicVersion)
	if err != nil {
		return nil
	}
	// must match major/minor/patch of default version
	var vDefault *utilversion.Version
	if defaultVersion == "v0.0.0-master+$Format:%H$" {
		// special-cases the placeholder value which doesn't parse as a semantic version
		vDefault, err = utilversion.ParseSemantic("v0.0.0-master")
	} else {
		vDefault, err = utilversion.ParseSemantic(defaultVersion)
	}
	if err != nil {
		return err
	}
	if vRuntime.Major() != vDefault.Major() || vRuntime.Minor() != vDefault.Minor() || vRuntime.Patch() != vDefault.Patch() {
		return fmt.Errorf("version %q must match major/minor/patch of default version %q", dynamicVersion, defaultVersion)
	}
	return nil
}
