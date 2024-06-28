package controller

import (
	"fmt"
	"github.com/aiteung/musik"
	"errors"
	"github.com/gofiber/fiber/v2"
	cek "github.com/rayfanaqbil/zenverse-BE/module"
	inimodel "github.com/rayfanaqbil/zenverse-BE/model"
	"github.com/rayfanaqbil/Zenverse-BP/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

// GetAllGames godoc
// @Summary Get All Data Games.
// @Description Mengambil semua data games.
// @Tags Games
// @Accept json
// @Produce json
// @Success 200 {object} Games
// @Router /games [get]
func GetAllGames(c *fiber.Ctx) error { 
	ps := cek.GetAllDataGames(config.Ulbimongoconn, "Games")
	return c.JSON(ps)
}

// GetGamesByID godoc
// @Summary Get By ID Data Games.
// @Description Ambil per ID data games.
// @Tags Games
// @Accept json
// @Produce json
// @Param id path string true "Masukan ID"
// @Success 200 {object} Games
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /games/{id} [get]
func GetGamesByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Wrong parameter",
		})
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid id parameter",
		})
	}
	ps, err := cek.GetGamesByID(objID, config.Ulbimongoconn, "Games")
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"status":  http.StatusNotFound,
				"message": fmt.Sprintf("No data found for id %s", id),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error retrieving data for id %s", id),
		})
	}
	return c.JSON(ps)
}

// InsertDataGames godoc
// @Summary Insert data Games.
// @Description Input data games.
// @Tags Games
// @Accept json
// @Produce json
// @Param request body ReqGames true "Payload Body [RAW]" 
// @Success 200 {object} Games
// @Failure 400
// @Failure 500
// @Router /insert [post]
func InsertDataGames(c *fiber.Ctx) error {
	db := config.Ulbimongoconn
	var games inimodel.Games
	if err := c.BodyParser(&games); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	fields := []struct {
		value string
		name  string
	}{
		{games.Name, "Name"},	
		{games.Desc, "Desc"},
		{games.Game_banner, "Game_banner"},
		{games.Preview, "Preview"},
		{games.Link_games, "Link_games"},
		{games.Game_logo, "Game_logo"},
		{games.Dev_name.Name, "Dev_name"},
		{games.Dev_name.Bio, "Bio"},
	}

	// Validasi input
	for _, field := range fields {
		if field.value == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"status":  http.StatusBadRequest,
				"message": "Fill all the form.",
				"field":   field.name,
			})
		}
	}

	if games.Rating == 0 || len(games.Genre) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Fill all the form.",
		})
	}

	insertedID, err := cek.InsertGames(db, "Games",
		games.Name,
		games.Rating,
		games.Desc,
		games.Genre,
		games.Dev_name,
		games.Game_banner,
		games.Preview,
		games.Link_games,
		games.Game_logo)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Data berhasil disimpan.",
		"inserted_id": insertedID,
	})
}

// UpdateDataGames godoc
// @Summary Update data Games.
// @Description Ubah data games.
// @Tags Games
// @Accept json
// @Produce json
// @Param id path string true "Masukan ID"
// @Param request body ReqGames true "Payload Body [RAW]"
// @Success 200 {object} Games
// @Failure 400
// @Failure 500
// @Router /update/{id} [put]
func UpdateDataGames(c *fiber.Ctx) error {
	db := config.Ulbimongoconn

	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	var games inimodel.Games
	if err := c.BodyParser(&games); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	err = cek.UpdateGames(db, "Games",
		objectID,
		games.Name,
		games.Rating,
		games.Desc,
		games.Genre,
		games.Dev_name,
		games.Game_banner,
		games.Preview,
		games.Link_games,
		games.Game_logo)
		if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Data successfully updated",
	})
}

// DeleteGamesByID godoc
// @Summary Delete data Games.
// @Description Hapus data games.
// @Tags Games
// @Accept json
// @Produce json
// @Param id path string true "Masukan ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /delete/{id} [delete]
func DeleteGamesByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Wrong parameter",
		})
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid id parameter",
		})
	}

	err = cek.DeleteGamesByID(objID, config.Ulbimongoconn, "Games")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error deleting data for id %s", id),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Data with id %s deleted successfully", id),
	})
}

func Login(c *fiber.Ctx) error {
	var loginDetails inimodel.Admin
	if err := c.BodyParser(&loginDetails); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
		})
	}

	_, err := cek.Login(config.Ulbimongoconn, "Admin", loginDetails.User_name, loginDetails.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Login successful",
	})
}

// GetDataAdmin godoc
// @Summary Get Data Admin.
// @Description Mengambil semua data admin.
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} Games
// @Router /admin [get]
func GetDataAdmin(c *fiber.Ctx) error { 
	ps := cek.GetDataAdmin(config.Ulbimongoconn, "Admin")
	return c.JSON(ps)
}


func GetGameByName(c *fiber.Ctx) error {
    name := c.Query("name")
    if name == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  fiber.StatusBadRequest,
            "message": "Query parameter 'name' is required",
        })
    }

    db := config.Ulbimongoconn
    games, err := cek.GetGamesByName(db, "Games", name)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  fiber.StatusNotFound,
                "message": fmt.Sprintf("Game with name '%s' not found", name),
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  fiber.StatusInternalServerError,
            "message": fmt.Sprintf("Error searching game by name: %v", err),
        })
    }

    return c.Status(fiber.StatusOK).JSON(games)
}