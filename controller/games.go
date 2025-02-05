package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
	"github.com/rayfanaqbil/Zenverse-BP/config"
	inimodel "github.com/rayfanaqbil/zenverse-BE/v2/model"
	cek "github.com/rayfanaqbil/zenverse-BE/v2/module"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Missing ID parameter",
        })
    }

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Invalid ID format",
        })
    }

    game, err := cek.GetGamesByID(objID, config.Ulbimongoconn, "Games")
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": err.Error(),
        })
    }


    return c.Status(http.StatusOK).JSON(game)
}

func EncryptIDHandler(c *fiber.Ctx) error {
    key, err := config.GetEncryptionKey()
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": err.Error(),
        })
    }

    id := c.Query("id")
    if id == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Missing ID parameter",
        })
    }

    encryptedID, err := config.EncryptID(key, id)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": "Failed to encrypt ID",
        })
    }

    return c.JSON(fiber.Map{
        "encrypted_id": encryptedID,
    })
}
//Description
func DecryptIDHandler(c *fiber.Ctx) error {
    key, err := config.GetEncryptionKey()
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": err.Error(),
        })
    }

    encryptedID := c.Query("encrypted_id")
    if encryptedID == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Missing encrypted_id parameter",
        })
    }

    decryptedID, err := config.DecryptID(key, encryptedID)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": "Failed to decrypt ID",
        })
    }

    return c.JSON(fiber.Map{
        "decrypted_id": decryptedID,
    })
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

	if games.Name == "" || games.Desc == "" || len(games.Genre) == 0 ||
    games.Dev_name.Name == "" ||  games.Dev_name.Bio == "" ||
	games.Game_banner == "" || games.Preview == "" ||
     games.Link_games == "" || games.Game_logo == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Fill all the form.",
		})
	}

	if games.Rating == 0 {
		games.Rating = 1.0
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


// GetGameByName godoc
// @Summary Get Game by Name.
// @Description Mengambil data game berdasarkan nama.
// @Tags Games
// @Accept json
// @Produce json
// @Param name query string true "Nama game yang dicari"
// @Success 200 {object} Games
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /games/search [get]
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

func GetGamesByRating(c *fiber.Ctx) error {
    ratingStr := c.Query("rating")
    if ratingStr == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  fiber.StatusBadRequest,
            "message": "Query parameter 'rating' is required",
        })
    }

    rating, err := strconv.ParseFloat(ratingStr, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  fiber.StatusBadRequest,
            "message": "Invalid 'rating' parameter",
        })
    }

    games, err := cek.GetGamesByHighestRating(config.Ulbimongoconn, "Games")
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  fiber.StatusNotFound,
                "message": fmt.Sprintf("No games found with rating %.2f", rating),
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  fiber.StatusInternalServerError,
            "message": fmt.Sprintf("Error retrieving games by rating: %v", err),
        })
    }


    var filteredGames []inimodel.Games
    for _, game := range games {
        if game.Rating >= rating {
            filteredGames = append(filteredGames, game)
        }
    }

    return c.Status(fiber.StatusOK).JSON(filteredGames)
}

func GetAllGamesApps(c *fiber.Ctx) error {
	skipParam := c.Query("skip", "0")
	
	skip, err := strconv.ParseInt(skipParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid skip parameter",
		})
	}

	ps := cek.GetAllDataGamesApps(config.Ulbimongoconn, "Games", skip)
	return c.JSON(ps)
}