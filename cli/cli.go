package cli

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Commands struct {
	BasePath              string
	ConfigFile            string
	Files                 []string
	ExternalLinksToIgnore []string
	InternalLinksToIgnore []string
	FilesToIgnore         []string
	Timeout               int
	RequestRepeats        int8
	AllowRedirect         bool
	AllowCodeBlocks       bool
	IgnoreExternal        bool
	IgnoreInternal        bool
	Verbose               bool
	FlagsSet              map[string]bool
}

func ParseCommands() Commands {
	basePath := flag.String("base-path", "", "The root source directories used to search for files")
	configFile := flag.String("config-file", "milv.config.yaml", "The config file for bot")
	externalLinksToIgnore := flag.String("external-links-to-ignore", "", "The list of external links to ignore")
	internalLinksToIgnore := flag.String("internal-links-to-ignore", "", "The list of internal links to ignore")
	filesToIgnore := flag.String("files-to-ignore", "", "The files to ignore")
	timeout := flag.Int("timeout", 0, "Timeout for http.get reguest")
	requestRepeats := flag.Int("request-repeats", 0, "Times reguest failuring links")
	allowRedirect := flag.Bool("allow-redirect", false, "Allow redirect")
	allowCodeBlocks := flag.Bool("allow-code-blocks", false, "Allow links in code blocks to check")
	ignoreInternal := flag.Bool("ignore-internal", false, "Ignore internal links")
	ignoreExternal := flag.Bool("ignore-external", false, "Ignore external links")
	verbose := flag.Bool("v", false, "Enable verbose logging")

	flag.Parse()
	files := flag.Args()

	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		flagset[f.Name] = true
	})

	if *basePath != "" {
		*configFile = fmt.Sprintf("%s/%s", *basePath, *configFile)
	}

	if len(files) == 0 {
		out := runCmd("find . -name \"*.md\"", true)
		files = strings.Split(string(out), "\n")
		if len(files) > 0 {
			files = files[:len(files)-1]
		}
	}

	return Commands{
		BasePath:              *basePath,
		ConfigFile:            *configFile,
		Files:                 files,
		ExternalLinksToIgnore: strings.Split(*externalLinksToIgnore, ","),
		InternalLinksToIgnore: strings.Split(*internalLinksToIgnore, ","),
		FilesToIgnore:         strings.Split(*filesToIgnore, ","),
		Timeout:               *timeout,
		RequestRepeats:        int8(*requestRepeats),
		AllowRedirect:         *allowRedirect,
		AllowCodeBlocks:       *allowCodeBlocks,
		IgnoreExternal:        *ignoreExternal,
		IgnoreInternal:        *ignoreInternal,
		Verbose:               *verbose,
		FlagsSet:              flagset,
	}
}

func runCmd(cmd string, shell bool) []byte {
	if shell {
		out, err := exec.Command("/bin/bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}
		return out
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}
