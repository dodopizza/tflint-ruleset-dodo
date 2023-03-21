# TFLint Ruleset Dodo
[![Build Status](https://github.com/dodopizza/tflint-ruleset-dodo/workflows/Test%20and%20publish/badge.svg?branch=main)](https://github.com/dodopizza/tflint-ruleset-dodo/actions)

This is a template repository with custom ruleset.
See also [Writing Plugins](https://github.com/terraform-linters/tflint/blob/master/docs/developer-guide/plugins.md) for more information.

## Requirements

- TFLint v0.30+
- Go v1.17

## Installation

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "dodo" {
  enabled = true

  version = "0.1.0"
  source  = "github.com/dodopizza/tflint-ruleset-dodo"
}
```

## Rules

All rules enabled by default and have `ERROR` severity.

| Name | Description |
| --- | --- |
| dodo_comments | Check that all comments written in consistent way. |
| dodo_file_content | Check that all files looks similarly, mostly focused on vertical alignment |
| dodo_terraform_backend_type | Check that modules specify `azurerm` as backend type |