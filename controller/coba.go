package controller

import (
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
	cek "github.com/rayfanaqbil/zenverse-BE/module"
)

func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

func GetGetAllGames(c *fiber.Ctx) error {
	ps := cek.GetAllDataGames()
	return c.JSON(ps)
}