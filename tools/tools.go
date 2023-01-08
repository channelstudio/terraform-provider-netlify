//go:build tools

package tools

import (
    // document generation
    _ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
    _ "github.com/Azure/go-autorest"
)