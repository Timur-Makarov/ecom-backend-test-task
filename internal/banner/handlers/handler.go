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

func (h BannerHandler) CreateBanner(c *fiber.Ctx) error {
	type RequestBody struct {
		Name string `json:"name"`
	}

	var reqBody RequestBody

	if err := c.BodyParser(&reqBody); err != nil {
		return domain.InvalidInputError
	}

	err := h.Service.CreateBanner(reqBody.Name)
	if err != nil {
		return domain.InternalServerError
	}

	return nil
}

func (h BannerHandler) UpdateCounterStatistics(c *fiber.Ctx) error {
	bannerIdString := c.Params("bannerID")
	bannerIdInt, err := strconv.ParseInt(bannerIdString, 10, 32)
	if err != nil || bannerIdInt <= 0 {
		return domain.InvalidPathError
	}

	h.Service.UpdateBannerCounterStats(int32(bannerIdInt))
	return nil
}

func (h BannerHandler) GetCounterStatistics(c *fiber.Ctx) error {
	bannerIdString := c.Params("bannerID")
	bannerIdInt, err := strconv.ParseInt(bannerIdString, 10, 32)
	if err != nil || bannerIdInt <= 0 {
		return domain.InvalidPathError
	}

	tsFromString := c.Query("tsFrom")
	tsToString := c.Query("tsTo")

	tsFromInt, err := strconv.ParseInt(tsFromString, 10, 64)
	if err != nil || tsFromInt < 0 {
		return domain.InvalidParamsError
	}

	tsToInt, err := strconv.ParseInt(tsToString, 10, 64)
	if err != nil || tsToInt < 0 {
		return domain.InvalidParamsError
	}

	res, err := h.Service.GetCounterStats(int32(bannerIdInt), tsFromInt, tsToInt)
	if err != nil {
		return err
	}

	err = c.JSON(res)
	if err != nil {
		return domain.InternalServerError
	}

	return nil
}
