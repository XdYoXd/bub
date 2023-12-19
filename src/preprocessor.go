package src

import (
	"fmt"
	"strings"
)

func Preprocessor(instructions []byte) (string, error) {

	variables := make(map[string]string)
	afterPreprocessor := string(instructions)
	var new_lines string

	/* get variables */
	for l, x := range strings.Split(afterPreprocessor, "\n") {
		if strings.HasPrefix(x, "#") {
			continue
		}

		if strings.HasPrefix(x, "const") {
			parts := strings.SplitN(strings.TrimPrefix(x, "const"), "=", 2)
			if len(parts) != 2 {
				return "", fmt.Errorf("%d: invalid const: '%s'", l+1, x)
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			variables[key] = value
		}
	}

	/*  replace consts*/
	for _, x := range strings.Split(afterPreprocessor, "\n") {
		if strings.HasPrefix(x, "#") {
			continue
		}

		if strings.HasPrefix(x, "const") {
			var line = x

			for k, v := range variables {
				line = strings.ReplaceAll(line, "${"+k+"}", v)
			}
			new_lines += line + "\n"

		}

	}

	/*  update variables */

	for l, x := range strings.Split(new_lines, "\n") {
		if strings.HasPrefix(x, "#") {
			continue
		}

		if strings.HasPrefix(x, "const") {
			parts := strings.SplitN(strings.TrimPrefix(x, "const"), "=", 2)
			if len(parts) != 2 {
				return "", fmt.Errorf("%d: invalid const: '%s'", l+1, x)
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			variables[key] = value
		}
	}

	/*  replace other */
	for _, x := range strings.Split(afterPreprocessor, "\n") {
		var line = x

		if strings.HasPrefix(x, "#") {
			continue
		}

		if !strings.HasPrefix(x, "const") {

			for k, v := range variables {
				line = strings.ReplaceAll(line, "${"+k+"}", v)
			}

			new_lines += line + "\n"
		}

	}

	return new_lines, nil
}
