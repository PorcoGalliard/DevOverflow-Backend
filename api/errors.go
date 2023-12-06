package api

import "github.com/gofiber/fiber/v2"

type Error struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, message string) Error {
	return Error{
		Code: code,
		Message: message,
	}
}

func (e Error) Error() string {
	return e.Message
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return ctx.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(fiber.StatusInternalServerError, err.Error())
	return ctx.Status(apiError.Code).JSON(apiError)
}

func ErrInvalidID() Error {
	return Error{
		Code: fiber.StatusBadRequest,
		Message: "Invalid ID",
	}
}

func ErrUnauthorized() Error {
	return Error{
		Code: fiber.StatusUnauthorized,
		Message: "Unauthorized User",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: fiber.StatusBadRequest,
		Message: "Invalid JSON Request",
	}
}

func ErrResourceNotFound(res string) Error {
	return Error{
		Code: fiber.StatusNotFound,
		Message: res + " Resource Not Found",
	}
}