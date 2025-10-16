package featureflags

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	OpenPipelinePipelineGroups FeatureFlag = "FEAT_OPENPIPELINE_PIPELINE_GROUP"
)

var defaultValues = map[FeatureFlag]defaultValue{
	OpenPipelinePipelineGroups: false,
}

type (
	// FeatureFlag represents a command line switch to turn certain features
	// ON or OFF. Values are read from environment variables defined by
	// the feature flag. The feature flag can have default value which is used
	// when the resp. environment variable does not exist
	FeatureFlag string

	defaultValue = bool
)

func (ff FeatureFlag) String() string {
	return ff.EnvName()
}

// EnvName gives back the environment variable name for the feature flag
func (ff FeatureFlag) EnvName() string {
	return string(ff)
}

// Enabled look up between known temporary and permanent flags and evaluates it.
// Feature flags are considered to be "enabled" if their resp. environment variable
// is set to 1, t, T, TRUE, true or True.
// Feature flags are considered to be "disabled" if their resp. environment variable
// is set to 0, f, F, FALSE, false or False.
func (ff FeatureFlag) Enabled() bool {
	v, exists := defaultValues[ff]
	if exists {
		return enabled(ff, v)
	}

	panic(fmt.Sprintf("unknown feature flag %s", ff))
}

// enabled evaluates the feature flag.
// Feature flags are considered to be "enabled" if their resp. environment variable
// is set to 1, t, T, TRUE, true or True.
// Feature flags are considered to be "disabled" if their resp. environment variable
// is set to 0, f, F, FALSE, false or False.
func enabled(ff FeatureFlag, d defaultValue) bool {
	if val, ok := os.LookupEnv(ff.EnvName()); ok {
		value, err := strconv.ParseBool(strings.ToLower(val))
		if err != nil {
			return d
		}
		return value
	}
	return d
}
