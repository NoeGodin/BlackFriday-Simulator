package Graphics

import (
	"fmt"
	"path/filepath"
	"strings"
)

func extractMapName(mapPath string) string {
	filename := filepath.Base(mapPath)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	dir := filepath.Base(filepath.Dir(mapPath))

	if dir != "." && dir != "" {
		return fmt.Sprintf("%s_%s", dir, name)
	}

	return name
}
