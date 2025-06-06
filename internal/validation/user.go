package validation

import (
	"errors"
	"strconv"
	"strings"

	"pet-project/internal/domain"
	"pet-project/internal/service"
)

func ValidateCreateUser(user domain.User) error {
	if user.FirstName == "" || user.LastName == "" {
		return errors.Join(service.ErrValidation, errors.New("first_name and last_name are required"))
	}
	if user.Age < 18 {
		return errors.Join(service.ErrValidation, errors.New("age must be at least 18"))
	}
	if len(user.Password) < 8 {
		return errors.Join(service.ErrValidation, errors.New("password must be at least characters"))
	}
	return nil
}

func ValidateUserID(path string) (int64, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 3 || parts[1] != "users" {
		return 0, errors.Join(service.ErrValidation, errors.New("invalid URL path"))
	}

	idStr := parts[2]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.Join(service.ErrValidation, errors.New("invalid user ID"))
	}
	return id, nil
}
