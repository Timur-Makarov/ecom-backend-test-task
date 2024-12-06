package handlers

import (
	"ecom-backend-test-task/internal/banner/domain"
	"ecom-backend-test-task/internal/banner/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type BannerHandler struct {
	Service services.BannerService
}

func (h BannerHandler) AddBanner(c *fiber.Ctx) error {
	type RequestBody struct {
		Name string `json:"name"`
	}

	var reqBody RequestBody

	if err := c.BodyParser(&reqBody); err != nil {
		return domain.InvalidInputError
	}

	err := h.Service.AddBanner(reqBody.Name)
	if err != nil {
		return domain.InternalServerError
	}

	return nil
}

func (h BannerHandler) UpdateCounterStats(c *fiber.Ctx) error {
	bannerIdString := c.Params("bannerID")
	bannerIdInt, err := strconv.Atoi(bannerIdString)

	if err != nil || bannerIdInt <= 0 {
		return domain.InvalidPathError
	}

	h.Service.UpdateBannerCounterStats(uint64(bannerIdInt))
	return nil
}

func (h BannerHandler) GetCounterStats(c *fiber.Ctx) error {
	bannerIdString := c.Params("bannerID")
	bannerIdInt, err := strconv.Atoi(bannerIdString)

	if err != nil || bannerIdInt <= 0 {
		return domain.InvalidPathError
	}

	tsFromString := c.Query("tsFrom")
	tsToString := c.Query("tsTo")

	tsFromInt, err := strconv.Atoi(tsFromString)
	if err != nil || tsFromInt < 0 {
		return domain.InvalidParamsError
	}

	tsToInt, err := strconv.Atoi(tsToString)
	if err != nil || tsToInt < 0 {
		return domain.InvalidParamsError
	}

	res, err := h.Service.GetCounterStats(uint64(bannerIdInt), uint64(tsFromInt), uint64(tsToInt))
	if err != nil {
		return domain.InternalServerError
	}

	err = c.JSON(res)
	if err != nil {
		return domain.InternalServerError
	}

	return nil
}
