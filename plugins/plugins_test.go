package main

import (
	"github.com/stangirard/yatas/internal/yatas"
)

var config = yatas.Config{
	Plugins: []yatas.Plugin{
		{
			Name:        "aws",
			Enabled:     true,
			Description: "AWS Plugin",
			Exclude:     []string{},
			Include:     []string{},
		},
	},
}
