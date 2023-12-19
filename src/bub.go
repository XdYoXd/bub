package src

import (
	"fmt"
	"os"
)

func (s *Session) Run() error {

	var i Instructions

	instructions, err := os.ReadFile(s.ConfigFile)
	if err != nil {
		return fmt.Errorf("bub: error: %s", err.Error())
	}

	i.bub, err = Preprocessor(instructions)
	if err != nil {
		return fmt.Errorf("bub: error: %s:%s", s.ConfigFile, err.Error())
	}

	err = i.ParseInstructions()
	if err != nil {
		return fmt.Errorf("bub: error: %s:%s", s.ConfigFile, err.Error())
	}
	return nil

}
