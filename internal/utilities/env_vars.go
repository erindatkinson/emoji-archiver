package utilities

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func ConfigOrEnv(prefix, key string) string {
	settings := viper.GetStringMapString(prefix)
	val, ok := settings[key]
	if !ok {
		val, _ = os.LookupEnv(
			fmt.Sprintf("%s_%s", strings.ToUpper(prefix), strings.ToUpper(key)))
	}
	return val
}

func PflagToBool(value pflag.Value) bool {
	if value.String() == "true" {
		return true
	} else if value.String() == "false" {
		return false
	} else {
		panic(fmt.Sprintf("invalid boolean value: %v", value))
	}
}
