// Package dto forms the bridge between the api data and the internal models
package dto

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	steamMatchTokenRegex = regexp.MustCompile(`^CSGO-[A-Za-z0-9]{5}(?:-[A-Za-z0-9]{5}){4}$`)
	steamAuthTokenRegex  = regexp.MustCompile(`^[A-Za-z0-9]{4}-[A-Za-z0-9]{5}-[A-Za-z0-9]{4}$`)
)

var Validate *validator.Validate

func InitValidator() error {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	return registerSteamValidators(Validate)
}

func registerSteamValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("steammatchtoken", func(fl validator.FieldLevel) bool {
		return steamMatchTokenRegex.MatchString(fl.Field().String())
	}); err != nil {
		return fmt.Errorf("register steam match token %w", err)
	}

	if err := v.RegisterValidation("steamauthtoken", func(fl validator.FieldLevel) bool {
		return steamAuthTokenRegex.MatchString(fl.Field().String())
	}); err != nil {
		return fmt.Errorf("register steam authentication token %w", err)
	}

	return nil
}
