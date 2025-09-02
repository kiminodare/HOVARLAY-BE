package history

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/internal/middleware"
	dtoHistory "github.com/kiminodare/HOVARLAY-BE/internal/modules/history/dto"
	historyInterface "github.com/kiminodare/HOVARLAY-BE/internal/modules/history/interface"
	"github.com/kiminodare/HOVARLAY-BE/internal/utils"
)

type Handler struct {
	service historyInterface.ServiceInterface
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req dtoHistory.CreateHistoryRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		return middleware.ValidationError(c, utils.FormatValidationErrors(err))
	}

	userIDStr := c.Locals("user_id").(string)

	if userIDStr == "" {
		return middleware.Error(c, "User ID not found", fiber.StatusUnauthorized)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return middleware.Error(c, "Invalid user ID format", fiber.StatusBadRequest)
	}

	history, err := h.service.Create(
		c.Context(),
		userID,
		req.Text,
		req.Voice,
		req.Rate,
		req.Pitch,
		req.Volume,
	)

	if err != nil {
		return middleware.Error(c, "Failed to create history", fiber.StatusInternalServerError)
	}

	return middleware.Success(c, history, "History created successfully", nil)
}

func (h *Handler) GetByUser(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	var query dtoHistory.GetHistoriesQuery
	if err := c.QueryParser(&query); err != nil {
		query.Page = 1
		query.Limit = 10
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}

	if userIDStr == "" {
		return middleware.Error(c, "User ID not found", fiber.StatusUnauthorized)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return middleware.Error(c, "Invalid user ID format", fiber.StatusBadRequest)
	}

	offset := (query.Page - 1) * query.Limit

	histories, err := h.service.GetByUser(c.Context(), userID, offset, query.Limit)
	if err != nil {
		return middleware.Error(c, "Failed to fetch history", fiber.StatusInternalServerError)
	}

	// Optional: hitung total, misal dari repo count
	total, _ := h.service.CountByUser(c.Context(), userID)

	// Buat pagination struct
	var pagination *middleware.Pagination
	if total >= 0 {
		pagination = &middleware.Pagination{
			Page:  query.Page,
			Limit: query.Limit,
			Total: total,
		}
	}

	return middleware.Success(c, histories, "History fetched successfully", pagination)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return middleware.Error(c, utils.ErrInvalidIDFormat, fiber.StatusBadRequest)
	}

	history, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return middleware.Error(c, "Failed to fetch history", fiber.StatusInternalServerError)
	}

	return middleware.Success(c, history, "History fetched successfully", nil)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return middleware.Error(c, "utils.ErrInvalidIDFormat", fiber.StatusBadRequest)
	}

	var req dtoHistory.UpdateHistoryRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		return middleware.ValidationError(c, utils.FormatValidationErrors(err))
	}

	err = h.service.Update(c.Context(), id, req.Text, req.Voice, req.Rate, req.Pitch, req.Volume)
	if err != nil {
		return middleware.Error(c, "Failed to update history", fiber.StatusInternalServerError)
	}

	return middleware.Success(c, nil, "History updated successfully", nil)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return middleware.Error(c, "utils.ErrInvalidIDFormat", fiber.StatusBadRequest)
	}

	err = h.service.Delete(c.Context(), id)
	if err != nil {
		return middleware.Error(c, "Failed to delete history", fiber.StatusInternalServerError)
	}

	return middleware.Success(c, nil, "History deleted successfully", nil)
}
