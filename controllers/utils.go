package controllers

import (
	"k8s.io/apimachinery/pkg/api/errors"
)

func IgnoreIsNotFound(err error) error {
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}
