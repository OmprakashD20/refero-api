package validator

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/OmprakashD20/refero-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	ValidatedBodyKey  = "validatedBody"
	ValidatedParamKey = "validatedParam"
)

func GetErrorMsg(fe validator.FieldError) string {
	field := utils.FormatFieldName(fe.Field())
	constraint := fe.Param()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "lte":
		return fmt.Sprintf("%s should be less than or equal to %s", field, constraint)
	case "gte":
		return fmt.Sprintf("%s should be greater than or equal to %s", field, constraint)
	case "min":
		return fmt.Sprintf("%s should have at least %s characters", field, constraint)
	case "max":
		return fmt.Sprintf("%s should have at most %s characters", field, constraint)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid ID", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, constraint)
	}

	return fmt.Sprintf("%s has an invalid value", field)
}

func ValidateBody[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body T
		if err := c.ShouldBindJSON(&body); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				fe := ve[0]
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": GetErrorMsg(fe)})
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

		c.Set(ValidatedBodyKey, body)
		c.Next()
	}
}

func ValidateParams[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uriParams T
		if err := c.ShouldBindUri(&uriParams); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				fe := ve[0]
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": GetErrorMsg(fe)})
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

		c.Set(ValidatedParamKey, uriParams)
		c.Next()
	}
}

func GetValidatedData[T any](c *gin.Context, key string) (T, bool) {
	val, exists := c.Get(key)
	if !exists {
		var empty T
		return empty, false
	}

	data, ok := val.(T)
	return data, ok
}
