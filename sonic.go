//go:build sonic

package versionchecker

import (
	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

func init() {
	clientPool.New = func() any {
		c := resty.New()
		c.JSONMarshal = sonic.Marshal
		c.JSONUnmarshal = sonic.Unmarshal
		return c
	}
}
