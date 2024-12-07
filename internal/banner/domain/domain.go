package domain

import "github.com/gofiber/fiber/v2"

var InvalidInputError = fiber.NewError(fiber.StatusBadRequest, "Invalid input data. Please check it and try again")
var InvalidPathError = fiber.NewError(fiber.StatusBadRequest, "Invalid URL path. Please check it and try again")
var InvalidParamsError = fiber.NewError(fiber.StatusBadRequest, "Invalid query params. Please check it and try again")
var InternalServerError = fiber.NewError(fiber.StatusInternalServerError, "Internal server error. Please try again later")

type Banner struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type CounterStatistic struct {
	BannerID      int32 `json:"bannerId"`
	TimestampFrom int64 `json:"timestampFrom"`
	TimestampTo   int64 `json:"timestampTo"`
	Count         int64 `json:"count"`
}
