package response

import (
	"errors"

	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	validatorService "github.com/AidlyTeam/Aidly-Backend/pkg/validator_service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BaseResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	DataCount  uint64      `json:"dataCount,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
}

func (r *BaseResponse) Error() string {
	return r.Message
}

func ResponseHandler(c *fiber.Ctx, err error) error {
	base := &BaseResponse{}
	//BaseResponse
	if errors.As(err, &base) {
		return c.Status(err.(*BaseResponse).StatusCode).JSON(err)
	}

	// Validation Errors
	if errors.As(err, &validator.ValidationErrors{}) {
		errs := validatorService.ValidatorErrors(err)
		return c.Status(400).JSON(
			&BaseResponse{
				StatusCode: 400,
				Message:    "validation error",
				Errors:     errs,
			},
		)
	}

	// Fiber Errors
	fiberErr := &fiber.Error{}
	if errors.As(err, &fiberErr) {
		if fiberErr.Code == 404 {
			return c.Status(404).JSON(&BaseResponse{
				StatusCode: 404,
				Message:    "not found",
			})
		} else {
			return c.Status(err.(*fiber.Error).Code).JSON(&BaseResponse{
				StatusCode: err.(*fiber.Error).Code,
				Message:    err.(*fiber.Error).Message,
			})
		}
	}

	// Service Errors
	serviceErr := &serviceErrors.ServiceError{}
	if errors.As(err, &serviceErr) {
		resp := &BaseResponse{
			StatusCode: serviceErr.Code,
			Message:    serviceErr.Message,
		}
		if serviceErr.Error() != "" {
			resp.Errors = serviceErr.Error()
		}
		return c.Status(serviceErr.Code).JSON(resp)
	}

	// Unknown Errors
	return c.Status(500).JSON(&BaseResponse{
		StatusCode: 500,
		Message:    "Internal Server Error (Unknown)",
		Errors:     err.Error(),
	})
}

// Response function for create a new response.
func Response(statusCode int, message string, data interface{}, dataCount ...uint64) error {
	var count uint64
	if len(dataCount) > 0 {
		count = dataCount[0]
	}
	return &BaseResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		DataCount:  count,
	}
}
