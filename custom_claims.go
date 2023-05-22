package main

import (
	"context"
	"errors"
)

type CustomClaimExample struct {
	Name         string `json:"name"`
	Username     string `json:"username"`
	ShouldReject bool   `json:"shouldReject,omitempty"`
}

func (c *CustomClaimExample) Validate(ctx context.Context) error {

	if c.ShouldReject {
		return errors.New("should reject was set to true")
	}

	return nil
}
