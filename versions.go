package main

import (
	"fmt"
)

func getFriendlyName(osStr, archStr string) string {
	var osFriendly string
	var archFriendly string

	osFriendly = osStr
	archFriendly = archStr

	if osStr == "js" && archStr == "wasm" {
		return "JavaScript as WebAssembly"
	}

	if osStr == "darwin" && archStr == "arm64" {
		return "macOS on Apple Silicon"
	}

	if val, ok := OSMap[osStr]; ok {
		osFriendly = val
	}

	if val, ok := ArchMap[archStr]; ok {
		archFriendly = val
	}

	return fmt.Sprintf("%s on %s", osFriendly, archFriendly)
}
