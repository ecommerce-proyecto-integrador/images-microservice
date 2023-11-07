package controllers

import (
	db "github.com/ecommerce-proyecto-integrador/images-microservice/config"
	"github.com/ecommerce-proyecto-integrador/images-microservice/models"

	/*amqp "github.com/rabbitmq/amqp091-go"*/
	"net/http"
	"os"
)

func UploadImage(image models.Image, fileName string, imageData []byte) (models.Image, error) {
	err := db.DB.Create(&image).Error

	if err != nil {
		return image, err
	}
	// Procesa la solicitud para subir la imagen utilizando el nombre del archivo y los datos de la imagen.
	// Guarda el archivo en el sistema de archivos.
	imagePath := "../images/" + fileName
	dest, err := os.Create(imagePath)
	if err != nil {
		return image, err
	}
	defer dest.Close()

	_, err = dest.Write(imageData)
	if err != nil {
		return image, err
	}

	return image, err
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	/*vars := mux.Vars(r)
	imageID := vars["id"]

	imagePath := "../images/" + imageID
	fmt.Println(imagePath)*/

	http.StripPrefix("/images/", http.FileServer(http.Dir("images"))).ServeHTTP(w, r)

}

func DeleteImage(imageName string) error {
	var image models.Image
	err := db.DB.Select("id, name").Where("name = ?", imageName).First(&image).Error

	if err != nil {
		return err
	}

	// Ruta completa al archivo de imagen que deseas eliminar.
	imagePath := "../images/" + imageName

	// Intenta eliminar el archivo.
	if err := os.Remove(imagePath); err != nil {
		return err
	}

	err = db.DB.Delete(&image).Error

	return err
}

/*
func (c *ImageController) GetProductWithImage(w http.ResponseWriter, r *http.Request) {
	// Implementa la lógica para obtener un producto con su imagen correspondiente
}*/
