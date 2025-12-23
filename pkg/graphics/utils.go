package Graphics

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
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

var (
	defaultColorScale = ebiten.ColorScale{}
	WinItemColorScale = func() ebiten.ColorScale {
		var cs ebiten.ColorScale
		cs.Scale(0.8, 2.5, 0.8, 1.0)
		return cs
	}()
	LooseItemColorScale = func() ebiten.ColorScale {
		var cs ebiten.ColorScale
		cs.Scale(2.5, 0.3, 0.3, 1.0)
		return cs
	}()

	guardsFOVColorScale = func() ebiten.ColorScale {
		var cs ebiten.ColorScale
		cs.Scale(0.2, 0.2, 0.2, 0.3)
		return cs
	}()
)
