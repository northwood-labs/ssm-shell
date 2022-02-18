package main

import "testing"

func TestGetFriendlyNameJS(t *testing.T) {
	expected := "JavaScript as WebAssembly"
	actual := getFriendlyName("js", "wasm")

	if actual != expected {
		t.Errorf("Result was `%s` instead of `%s`.", actual, expected)
	}
}

func TestGetFriendlyNameASi(t *testing.T) {
	expected := "macOS on Apple Silicon"
	actual := getFriendlyName("darwin", "arm64")

	if actual != expected {
		t.Errorf("Result was `%s` instead of `%s`.", actual, expected)
	}
}

func TestGetFriendlyNameLinux64(t *testing.T) {
	expected := "Linux on Intel (64-bit)"
	actual := getFriendlyName("linux", "amd64")

	if actual != expected {
		t.Errorf("Result was `%s` instead of `%s`.", actual, expected)
	}
}

func TestGetFriendlyNameNoOS(t *testing.T) {
	expected := "illumos on Intel (64-bit)"
	actual := getFriendlyName("illumos", "amd64")

	if actual != expected {
		t.Errorf("Result was `%s` instead of `%s`.", actual, expected)
	}
}

func TestGetFriendlyNameNoOSArch(t *testing.T) {
	expected := "z/OS on System/390 (64-bit)"
	actual := getFriendlyName("zos", "s390x")

	if actual != expected {
		t.Errorf("Result was `%s` instead of `%s`.", actual, expected)
	}
}
