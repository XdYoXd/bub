package src

import "flag"

type Session struct {
	ConfigFile string
}

type Instructions struct {
	bub string
}

func New() *Session {
	var configFile string

	flag.StringVar(&configFile, "f", "build.bub", "")
	flag.Parse()

	return &Session{
		ConfigFile: configFile,
	}
}
