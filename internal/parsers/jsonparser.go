package parsers

import (
	"encoding/json"
	"log"
)

func MustParseJson[T any](bts []byte) T {
	var thing T
	if err := json.Unmarshal(bts, &thing); err != nil {
		log.Fatalf("failed to parse, check yourself: %s", err.Error())
	}
	return thing
}