package controllers

import (
	"devcode-pos/helpers"
	"devcode-pos/models"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetPasscode(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var getPasscode models.GetPasscode

	res := models.DB.WithContext(c.Context()).
		Table("cashiers").
		Where(&models.Cashiers{CashierId: uint32(id)}).
		First(&getPasscode)

	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   make(map[string]string),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    getPasscode,
	})
}

func VerifyLogin(c *fiber.Ctx) error {
	var err error
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var cashierModel models.Cashiers
	var input = struct {
		Passcode string `validate:"required" json:"passcode"`
	}{}

	// Parsing Input
	if err = c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Bad Request",
			"error":   make([]int, 0),
		})
	}

	// Run Validation
	if msg, err := helpers.RunValidate(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "test",
			"error":   msg,
		})
	}

	res := models.DB.WithContext(c.Context()).
		Where(&models.Cashiers{CashierId: uint32(id)}).
		First(&cashierModel)

	if res.Error != nil || cashierModel.CashierPasscode != input.Passcode {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Passcode Not Match",
			"error":   make(map[string]string),
		})
	}

	// Generate Jwt Token
	claims := jwt.MapClaims{
		"cashierId": cashierModel.CashierId,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	getToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Bad Request",
			"data":    make(map[string]string),
		})
	}

	cashierModel.CashierToken = getToken
	go models.DB.Save(&cashierModel)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": map[string]string{
			"token": getToken,
		},
	})
}

func VerifyLogout(c *fiber.Ctx) error {
	var err error
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	var cashierModel models.Cashiers
	var input = struct {
		Passcode string `validate:"required" json:"passcode"`
	}{}

	// Parsing Input
	if err = c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Bad Request",
			"error":   make([]int, 0),
		})
	}

	// Run Validation
	if msg, err := helpers.RunValidate(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "test",
			"error":   msg,
		})
	}

	res := models.DB.WithContext(c.Context()).
		Where(&models.Cashiers{CashierId: uint32(id)}).
		First(&cashierModel)

	if res.Error != nil || cashierModel.CashierPasscode != input.Passcode {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Passcode Not Match",
			"error":   make(map[string]string),
		})
	}

	cashierModel.CashierToken = ""
	go models.DB.Save(&cashierModel)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}
