package main

import (
	"strings"
)

func parseVars(vars []string) map[string]*string {
	res := make(map[string]*string)
	for _, v := range vars {
		if strings.Contains(v, "=") {
			parts := strings.SplitN(v, "=", 2)
			res[parts[0]] = &parts[1]
		} else {
			res[v] = nil
		}
	}

	return res
}
