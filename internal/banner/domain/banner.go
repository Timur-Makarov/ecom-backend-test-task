package domain

import "github.com/gofiber/fiber/v2"

// The idea is that each layer should have it's own domain models
// so that the layers below would have to adjust themselves.
// For now, there is no usage for them.

var InvalidInputError = fiber.NewError(fiber.StatusBadRequest, "Invalid input data. Please check it and try again")
var InvalidPathError = fiber.NewError(fiber.StatusBadRequest, "Invalid URL path. Please check it and try again")
var InvalidParamsError = fiber.NewError(fiber.StatusBadRequest, "Invalid query params. Please check it and try again")
var InternalServerError = fiber.NewError(fiber.StatusInternalServerError, "Internal server error. Please try again later")

type Banner struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type CounterStats struct {
	BannerID      uint64 `json:"bannerId"`
	TimestampFrom uint64 `json:"timestampFrom"`
	TimestampTo   uint64 `json:"timestampTo"`
	Count         uint64 `json:"count"`
}
