package device

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// PairingHandler handles HTTP requests for RFID pairing
type PairingHandler struct {
	service PairingService
}

// NewPairingHandler creates a new pairing handler
func NewPairingHandler(service PairingService) *PairingHandler {
	return &PairingHandler{service: service}
}

// RegisterRoutes registers pairing routes for admin
func (h *PairingHandler) RegisterRoutes(router fiber.Router) {
	pairing := router.Group("/pairing")
	pairing.Post("/start", h.StartPairing)
	pairing.Post("/cancel/:deviceId", h.CancelPairing)
	pairing.Get("/status/:deviceId", h.GetPairingStatus)
}

// RegisterRoutesWithoutGroup registers pairing routes without creating a sub-group
func (h *PairingHandler) RegisterRoutesWithoutGroup(router fiber.Router) {
	router.Post("/start", h.StartPairing)
	router.Post("/cancel/:deviceId", h.CancelPairing)
	router.Get("/status/:deviceId", h.GetPairingStatus)
}

// RegisterPublicRoutes registers public routes for ESP32 devices
func (h *PairingHandler) RegisterPublicRoutes(router fiber.Router) {
	router.Post("/pairing/rfid", h.ProcessRFIDPairing)
}

// StartPairing handles starting a pairing session
// @Summary Start RFID pairing session
// @Description Start a pairing session to link an RFID card to a student
// @Tags Pairing
// @Accept json
// @Produce json
// @Param request body StartPairingRequest true "Pairing request"
// @Success 200 {object} PairingSessionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/pairing/start [post]
func (h *PairingHandler) StartPairing(c *fiber.Ctx) error {
	var req StartPairingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	// Validate required fields
	if req.DeviceID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "ID perangkat wajib diisi",
			},
		})
	}
	if req.StudentID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "ID siswa wajib diisi",
			},
		})
	}

	response, err := h.service.StartPairing(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ProcessRFIDPairing handles RFID tap during pairing mode (from ESP32)
// @Summary Process RFID pairing
// @Description Process an RFID card tap during pairing mode
// @Tags Pairing
// @Accept json
// @Produce json
// @Param request body RFIDPairingRequest true "RFID pairing data"
// @Success 200 {object} RFIDPairingResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/pairing/rfid [post]
func (h *PairingHandler) ProcessRFIDPairing(c *fiber.Ctx) error {
	var req RFIDPairingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	if req.APIKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "API key wajib diisi",
			},
		})
	}
	if req.RFIDCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Kode RFID wajib diisi",
			},
		})
	}

	response, err := h.service.ProcessRFIDPairing(c.Context(), req)
	if err != nil {
		// For pairing errors, still return the response with success=false
		if errors.Is(err, ErrRFIDAlreadyUsed) {
			return c.JSON(fiber.Map{
				"success": false,
				"data":    response,
			})
		}
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": response.Success,
		"data":    response,
	})
}

// CancelPairing handles cancelling a pairing session
// @Summary Cancel pairing session
// @Description Cancel an active pairing session
// @Tags Pairing
// @Produce json
// @Param deviceId path int true "Device ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/pairing/cancel/{deviceId} [post]
func (h *PairingHandler) CancelPairing(c *fiber.Ctx) error {
	deviceID, err := strconv.ParseUint(c.Params("deviceId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID perangkat tidak valid",
			},
		})
	}

	if err := h.service.CancelPairing(c.Context(), uint(deviceID)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Sesi pairing dibatalkan",
	})
}

// GetPairingStatus handles getting pairing session status
// @Summary Get pairing status
// @Description Get the status of a pairing session for a device
// @Tags Pairing
// @Produce json
// @Param deviceId path int true "Device ID"
// @Success 200 {object} PairingSessionResponse
// @Failure 400 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/pairing/status/{deviceId} [get]
func (h *PairingHandler) GetPairingStatus(c *fiber.Ctx) error {
	deviceID, err := strconv.ParseUint(c.Params("deviceId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID perangkat tidak valid",
			},
		})
	}

	response, err := h.service.GetPairingStatus(c.Context(), uint(deviceID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// handleError handles service errors
func (h *PairingHandler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrDeviceNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_DEVICE",
				"message": "Perangkat tidak ditemukan",
			},
		})
	case errors.Is(err, ErrDeviceInactive):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DEVICE_INACTIVE",
				"message": "Perangkat tidak aktif",
			},
		})
	case errors.Is(err, ErrStudentAlreadyPaired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_ALREADY_PAIRED",
				"message": "Siswa sudah memiliki kartu RFID. Hapus kartu lama terlebih dahulu.",
			},
		})
	case errors.Is(err, ErrRFIDAlreadyUsed):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_RFID_USED",
				"message": "Kartu RFID sudah digunakan siswa lain",
			},
		})
	case errors.Is(err, ErrNoPairingSession):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_NO_SESSION",
				"message": "Tidak ada sesi pairing aktif",
			},
		})
	case errors.Is(err, ErrPairingSessionExpired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_SESSION_EXPIRED",
				"message": "Sesi pairing sudah kadaluarsa",
			},
		})
	case errors.Is(err, ErrInvalidAPIKey):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_INVALID_API_KEY",
				"message": "API key tidak valid",
			},
		})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "ERROR",
				"message": err.Error(),
			},
		})
	}
}
