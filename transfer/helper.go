package transfer

import "strings"

func formatPath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	if strings.HasPrefix(".", path) {
		path = path[1:]
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}
