package controller

import (
	"fmt"
	"github.com/rayfanaqbil/Zenverse-BP/model"
	"github.com/gofiber/fiber/v2"
	"github.com/rayfanaqbil/Zenverse-BP/config"
	"github.com/rayfanaqbil/Zenverse-BP/helper"
	"go.mongodb.org/mongo-driver/bson"
)

func PostUploadGithub(c *fiber.Ctx) error {
	fmt.Println("Starting file upload process")


	file, err := c.FormFile("img")
	if err != nil {
		fmt.Println("Error parsing form file:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	folder := c.Query("folder", "") 
	var pathFile string
	if folder != "" {
		pathFile = folder + "/" + file.Filename
	} else {
		pathFile = file.Filename
	}


	gh, err := helper.GetOneDoc[model.Ghcreates](config.Ulbimongoconn, "github", bson.M{})
	if err != nil {
		fmt.Println("Error fetching GitHub credentials:", err)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}


	content, _, err := helper.GithubUpload(
		gh.GitHubAccessToken, gh.GitHubAuthorName, gh.GitHubAuthorEmail, file, 
		"zenverse-assets", "filegambarzenverse", pathFile, false,
	)

	if err != nil {
		fmt.Println("Error uploading file to GitHub:", err)
		return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
			"error": "Gagal upload ke GitHub",
			"detail": err.Error(),
		})
	}

	if content == nil || content.Content == nil {
		fmt.Println("Error: content or content.Content is nil")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error uploading file",
		})
	}

	fmt.Println("File upload process completed successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"name": *content.Content.Name,
		"path": *content.Content.Path,
	})
}
