package picture

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// type Req struct {
// 	Name      string   `form:"name"`
// }

func SavePictures(c *gin.Context) {
	// var req := Req
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing form: %s", err.Error()))
		return
	}
	files := form.File["upload[]"]

	for _, file := range files {
		log.Println(file.Filename)

		// Specify the destination directory where the file will be saved
		dst := "./uploads/" + file.Filename

		// Save the uploaded file to the specified destination
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error saving file: %s", err.Error()))
			return
		}
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))

}
