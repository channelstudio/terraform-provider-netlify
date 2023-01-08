---
layout: "netlify"
page_title: "Netlify: netlify_environment_variable"
sidebar_current: "docs-netlify-resource-environment-variable"
description: |-
  Provides an environment variable resource.
---

# netlify_environment_variable

Allows configuration of either an account-wide or site-specific environment variable
across multiple deploy contexts and scopes.

## Example Usage

```hcl
resource "netlify_environment_variable" "main" {
  key = "my_env_variable"

  repo {
    command       = "middleman build"
    deploy_key_id = "${netlify_deploy_key.key.id}"
    dir           = "/build"
    provider      = "github"
    repo_path     = "username/repo"
    repo_branch   = "master"
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) - The Netlify account ID / slug for the environment variable.
* `key` - (Required) - The key for the environment variable.
* `site_id` - (Optional) - If provided, creates the environment variable on the site level, not the account level.
* `scopes` - (Optional) - Scopes that this environment variable is set to (Netlify Pro plans and above). Use any combination of [`builds`, `functions`, `post_processing`, `runtime`] If unset, defaults to all scopes.
* `value` - (Required) - See [Value](#value)

### Values

`value` supports the following arguments:

* `value` - (Required) - The value of the environment variable in this context
* `context` - (Optional) - The deploy context in which this value is available. `dev` refers to local development when running `netlify dev`. Enum: [`all` `dev` `branch-deploy` `deploy-preview` `production`]
