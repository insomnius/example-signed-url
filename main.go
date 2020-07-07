package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/codefluence-x/aurelia"
	"github.com/gin-gonic/gin"
)

func main() {

	engine := gin.Default()
	secretCode := "Programmer is a machine who turns coffe into code"

	engine.GET("/users", func(c *gin.Context) {

		expiresAt := time.Now().Add(15 * time.Second).Unix()
		imageName := "image.jpg"

		signature := aurelia.Hash(secretCode, fmt.Sprintf("%d%s", expiresAt, imageName))

		encoder := json.NewEncoder(c.Writer)
		encoder.SetEscapeHTML(false)
		_ = encoder.Encode(gin.H{
			"data": gin.H{
				"id":        1,
				"name":      "Smitty Werben Men Jensen",
				"image_url": fmt.Sprintf("http://localhost:8080/avatar/image.jpg?signature=%s&expires_at=%d", signature, expiresAt),
			},
		})

		c.Status(200)
		c.Writer.Header().Set("Content-Type", "application/json")
	})

	engine.GET("/avatar/image.jpg", func(c *gin.Context) {
		signature := c.Request.URL.Query().Get("signature")
		expiresAt := c.Request.URL.Query().Get("expires_at")

		if signature == "" || expiresAt == "" {
			c.JSON(400, gin.H{
				"message": "signature and expires_at cannot be empty",
			})
			return
		}

		expiresAtUnix, err := strconv.Atoi(expiresAt)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "invalid format of expires_at",
			})
			return
		}

		fmt.Println("EXPIRES", expiresAtUnix)
		fmt.Println("SIGNATURE", signature)

		if aurelia.Authenticate(secretCode, fmt.Sprintf("%d%s", expiresAtUnix, "image.jpg"), signature) == false {
			c.JSON(403, gin.H{
				"message": "unauthorized",
			})
			return
		}

		if time.Now().After(time.Unix(int64(expiresAtUnix), 0)) {
			c.JSON(404, gin.H{
				"message": "image not found",
			})
			return
		}

		c.File("./image.jpg")
	})

	engine.Run(":8080")
}
