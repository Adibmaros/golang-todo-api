package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func validationErrorResponse(err error) gin.H {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make(map[string]string, len(ve))
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			fields[field] = formatFieldError(field, fe)
		}
		return gin.H{
			"error":  "validasi gagal",
			"fields": fields,
		}
	}
	return gin.H{
		"error": "request body tidak valid",
	}
}

func formatFieldError(field string, fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s harus diisi", field)
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", field, fe.Param())
	default:
		return fmt.Sprintf("%s tidak valid", field)
	}
}
