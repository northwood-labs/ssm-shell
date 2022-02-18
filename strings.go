package main

var (
	// OSMap is a mapping of GOOS values to "friendly" values.
	// https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
	// https://github.com/golang/go/blob/master/src/go/build/syslist.go
	OSMap = map[string]string{
		"aix":       "AIX",
		"android":   "Android",
		"darwin":    "macOS",
		"dragonfly": "DragonFly BSD",
		"freebsd":   "FreeBSD",
		"hurd":      "GNU Hurd",
		"illumos":   "illumos",
		"ios":       "iOS",
		"js":        "JavaScript",
		"linux":     "Linux",
		"netbsd":    "NetBSD",
		"openbsd":   "OpenBSD",
		"plan9":     "Plan 9",
		"solaris":   "Solaris",
		"windows":   "Windows",
		"zos":       "z/OS",
	}

	// ArchMap is a mapping of GOARCH values to "friendly" values.
	// https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
	// https://github.com/golang/go/blob/master/src/go/build/syslist.go
	ArchMap = map[string]string{
		"386":         "Intel (32-bit)",
		"amd64":       "Intel (64-bit)",
		"amd64p32":    "Intel (64-bit)",
		"arm":         "ARM (32-bit)",
		"arm64":       "ARM (64-bit)",
		"arm64be":     "ARM (64-bit)",
		"armbe":       "ARM (32-bit)",
		"mips":        "MIPS (32-bit)",
		"mips64":      "MIPS (64-bit)",
		"mips64le":    "MIPS (64-bit)",
		"mips64p32":   "MIPS (64-bit)",
		"mips64p32le": "MIPS (64-bit)",
		"mipsle":      "MIPS (32-bit)",
		"ppc":         "PowerPC (32-bit)",
		"ppc64":       "PowerPC (64-bit)",
		"ppc64le":     "PowerPC (64-bit)",
		"riscv":       "RISC-V (32-bit)",
		"riscv64":     "RISC-V (64-bit)",
		"s390":        "System/390 (32-bit)",
		"s390x":       "System/390 (64-bit)",
		"sparc":       "SPARC (32-bit)",
		"sparc64":     "SPARC (34-bit)",
		"wasm":        "WebAssembly",
	}
)
