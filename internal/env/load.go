package env

import (
	"os"
	"strings"
)

func Load(pathToEnv string) {
	bts, err := os.ReadFile(pathToEnv)
	if err != nil {
		panic(err)
	}
	txt := string(bts)
	for _, line := range strings.Split(txt, "\n") {
		parts := strings.Split(line, "=")
		os.Setenv(parts[0], strings.ReplaceAll(strings.TrimSpace(parts[1]), "'", ""))
	}
}
