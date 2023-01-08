//go:build tools

package tools

import (
	// document generation
	_ "github.com/Azure/go-autorest"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
