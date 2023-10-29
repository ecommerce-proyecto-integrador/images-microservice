package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Definir el directorio donde se almacenarán las imágenes
	imageDir := "./images"
	err := os.MkdirAll(imageDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creando el directorio de imágenes:", err)
		return
	}

	// Configurar las rutas para servir las imágenes
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imageDir))))

	// Iniciar el servidor de imágenes
	port := "8080" // Puedes cambiar el puerto según tus necesidades
	fmt.Printf("Servidor de imágenes en ejecución en el puerto %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
