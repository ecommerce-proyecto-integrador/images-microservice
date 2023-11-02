package controllers

import (
	"fmt"
	"os"
)

func saveImage(image []byte, fileName string) error {
	err := os.WriteFile(fmt.Sprintf("./images/%s", fileName), image, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Elimina una imagen del disco local.
func deleteImage(fileName string) error {
	err := os.Remove(fmt.Sprintf("./images/%s", fileName))
	if err != nil {
		return err
	}
	return nil
}
