package controllers

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/MichaelYoung87/backend-go-project-remake/database"
	"github.com/MichaelYoung87/backend-go-project-remake/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"path/filepath"
	"strconv"
	"time"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Username: data["username"],
		Password: password,
		ImageURL: data["image_url"],
	}
	database.DB.Create(&user)
	return c.JSON(user)
}
func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	var user models.User
	database.DB.Where("username = ?", data["username"]).First(&user)
	if user.UserID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Username not found",
		})
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}
	currentTime := time.Now()
	fmt.Println("The time is:", currentTime)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(user.UserID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), //1 day - ny lösning
		//ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day - gammal fungerar inte med jwt/v5 .Unix int32 "time" pack = int64
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not login",
		})
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "Success!",
	})
}
func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized jwt",
		})
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized claims",
		})
	}
	var user models.User
	database.DB.Where("user_id = ?", claims.Issuer).First(&user)
	return c.JSON(user)
}
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "Logged out",
	})
}
func CheckIfUsernameExists(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	username := data["username"]
	var user models.User
	database.DB.Where("username = ?", username).First(&user)
	response := fiber.Map{
		"exists": user.UserID != 0,
	}
	return c.JSON(response)
}
func UploadProfilePicture(c *fiber.Ctx) error {
	// Hämtar user_id från JWT token
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized jwt",
		})
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized claims",
		})
	}
	var user models.User
	database.DB.Where("user_id = ?", claims.Issuer).First(&user)
	fmt.Println("User:", user)

	// Hämtar den uppladdade filen från requesten.
	file, err := c.FormFile("profilePicture")
	if err != nil {
		fmt.Println("Error getting file:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	fmt.Println("File:", file)

	// Skapar ny writer för filen som finns i Google Cloud Storage Bucket.
	fileName := uuid.New().String() + filepath.Ext(file.Filename)
	bucketName := "react_project_remake_profile_pictures" //Detta är namnet på min bucket i Google Cloud Storage.
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	fmt.Printf("Client: %+v\n", client)
	bucketHandle := client.Bucket(bucketName)
	fmt.Printf("Bucket handle: %+v\n", bucketHandle)
	objHandle := bucketHandle.Object(fileName)
	writer := objHandle.NewWriter(ctx)
	defer writer.Close()

	// Hämtar in filen och ställer in hur filen ska formateras innan den sparas Google Cloud Storage Bucket.
	writer.ContentType = file.Header.Get("Content-Type")
	writer.ContentDisposition = fmt.Sprintf("inline; filename=\"%s\"", file.Filename)

	// Kopierar innehållet av den uppladdade filen till writern.
	fileData, err := file.Open()
	if err != nil {
		fmt.Println("Error opening file:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer fileData.Close()
	if _, err := io.Copy(writer, fileData); err != nil {
		fmt.Println("Error copying file data:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if err := writer.Close(); err != nil {
		fmt.Println("Error closing writer:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Hämtar den publika URL:en för den uppladdade filen
	objAttrs, err := objHandle.Attrs(ctx)
	if err != nil {
		fmt.Println("Error getting object attributes:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	imageUrl := objAttrs.MediaLink
	fmt.Println("Image URL:", imageUrl)
	// Sätter ACL (Access control lists) för objektet till public read access.
	if err := objHandle.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		fmt.Println("Error setting ACL:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Uppdaterar användaren i databasen med den nya profilbildens URL adress som nu finns i Google Cloud Storage Bucket.
	if err := database.DB.Model(&user).Update("ImageURL", imageUrl).Error; err != nil {
		fmt.Println("Error updating user record:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Profile picture updated successfully",
	})
}