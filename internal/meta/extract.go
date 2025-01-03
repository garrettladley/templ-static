package meta

import (
	"bufio"
	"io"
	"strings"
)

func Extract(r io.Reader, path string) (Meta, bool) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// stop parsing when we hit the package declaration
		if strings.HasPrefix(line, "package ") {
			break
		}

		// split the line into fields and check for the directive
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "//templ-static" {
			meta := Meta{FilePath: path}
			// parse the path if provided
			for _, field := range fields[1:] {
				if strings.HasPrefix(field, "path:") {
					meta.Path = strings.TrimSpace(strings.TrimPrefix(field, "path:"))
				}
			}

			return meta, true
		}
	}

	return Meta{}, false
}
