package src

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Rule struct {
	Command     string
	Description string /*  optional */
	Output      string /*  optional. stdout | discard  */
	Errors      string /*  optional. stderr | discard  */
	Line        int    /*  bub      */
}

func (i Instructions) ParseInstructions() error {

	for l, v := range strings.Split(i.bub, "\n") {

		sp := strings.Split(v, " ")

		if sp[0] == "BUILD:" {
			/*
				sp[0] == BUILD:
				sp[1:] == RULE NAME
			*/
			if len(sp) < 2 {
				return fmt.Errorf("%d: invalid build: %s", l+1, v)
			}

			if len(sp[1]) == 0 || sp[1] == " " {
				return fmt.Errorf("%d: invalid rule name: %s", l+1, v)
			}

			if len(sp) > 2 {
				for j, x := range sp[1:] {
					rule, err := i.getRule(x, l+1)
					if err != nil {
						return err
					}

					if rule.Description == "" {
						fmt.Printf("[%d/%d] %s\n", j+1, len(sp[1:]), rule.Command)
					} else {
						fmt.Printf("[%d/%d] %s", j+1, len(sp[1:]), rule.Description)
					}

					rule.execute()
				}

				return nil
			}

			rule, err := i.getRule(sp[1], l+1)
			if err != nil {
				return err
			}

			if rule.Description == "" {
				fmt.Printf("[1/1] %s\n", rule.Command)
			} else {
				fmt.Printf("[1/1] %s", rule.Description)
			}

			rule.execute()
			/*  return */
			// fmt.Println(rule)

		}

	}

	return nil

}

func (i Instructions) getRule(rule string, build_line int) (*Rule, error) {
	var inBlock bool
	var found bool
	var r Rule

	for l, v := range strings.Split(i.bub, "\n") {

		sp := strings.Split(v, " ")

		if sp[0] == "RULE" {
			/*
				sp[0] == RULE
				sp[1] == RULE NAME
				sp[2] == {
			*/

			if len(sp) < 3 {
				return nil, fmt.Errorf("%d: invalid rule: '%s'", l+1, v)
			}
			if sp[1] == " " {
				return nil, fmt.Errorf("%d: invalid rule: '%s'", l+1, v)
			}
			if sp[2] != "{" {
				return nil, fmt.Errorf("%d: invalid rule: '%s'", l+1, v)
			}

			if sp[1] != rule {
				continue
			}

			if found {
				return nil, fmt.Errorf("%d: duplicate rule: '%s'", l+1, sp[1])
			}

			r.Line = l + 1
			inBlock = true
			found = true
		}

		if sp[0] == "}" {
			inBlock = false
		}

		if inBlock {
			if sp[0] == "RULE" { /*  ignore rule line */
				continue
			}

			parts := strings.SplitN(v, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				switch key {
				case "COMMAND":
					r.Command = value
				case "DESCRIPTION":
					r.Description = value
				case "ERRORS":
					r.Errors = value
					if r.Errors != "" && r.Errors != "discard" && r.Errors != "stderr" {
						return nil, fmt.Errorf("%d: unknown errors name: %s\nHint: Use 'discard' or 'stderr'", l+1, r.Output)
					}
					if r.Errors == "" { /*  set default value */
						r.Errors = "stderr"
					}

				case "OUTPUT":
					r.Output = value
					if r.Output != "" && r.Output != "discard" && r.Output != "stdout" {
						return nil, fmt.Errorf("%d: unknown output name: %s\nHint: Use 'discard' or 'stdout'", l+1, r.Output)
					}
					if r.Output == "" {
						r.Output = "stdout"
					}

				}
			}

		}

	}

	if !found {
		return nil, fmt.Errorf("%d: unknown build rule: '%s'", build_line, rule)
	}

	if r.Command == "" {
		return nil, fmt.Errorf("%d: expected 'COMMAND =' line", r.Line)
	}

	return &r, nil
}

func (r Rule) execute() {
	temp := os.TempDir()
	shell_name := RandomString(20) + ".sh"

	os.WriteFile(temp+"/"+shell_name, []byte(r.Command), os.ModePerm)

	cmd := exec.Command("bash", temp+"/"+shell_name)

	dir, _ := os.Getwd()
	cmd.Dir = dir
	cmd.Env = os.Environ()

	if r.Output == "stdout" {
		cmd.Stdout = os.Stdout
	}

	if r.Errors == "stderr" {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		fmt.Println(cmd.Err)
		os.Exit(1)
	}

	os.Remove(temp + "/" + shell_name)
}
