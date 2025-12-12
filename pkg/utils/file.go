package utils

import (
	"fmt"
	"io"
	"os"
)

func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("impossible d'ouvrir le fichier source: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier destination: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("erreur pendant la copie: %w", err)
	}

	err = dstFile.Sync()
	if err != nil {
		return fmt.Errorf("erreur lors du flush: %w", err)
	}

	return nil
}
