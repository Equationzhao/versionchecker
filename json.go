//go:build !sonic

package versionchecker

import "github.com/go-resty/resty/v2"

func init() {
	clientPool.New = func() any {
		c := resty.New()

		return c
	}
}
