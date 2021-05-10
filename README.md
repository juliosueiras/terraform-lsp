# Terraform LSP

[![Gitter](https://badges.gitter.im/terraform-lsp/community.svg)](https://gitter.im/terraform-lsp/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
![terraform version](https://img.shields.io/badge/terraform-0.13.0-blue.svg)
![Release](https://github.com/juliosueiras/terraform-lsp/workflows/Release/badge.svg)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp?ref=badge_shield)

This is LSP (Language Server Protocol) for Terraform

**IMPORTANT:** Currently there is two terraform lsp, one is this one and the other one is [terraform-ls](https://github.com/hashicorp/terraform-ls), which contain details about this repo as well.

**New update to fix for Terraform 0.15 will be up in a day or two**

**Current Focus: Terraform State Reading**

The aim to have a unified lsp for terraform in the future, but for now there is two concurrent development with collabration to each other, this repo is aim for more experimental features, and the terraform-ls is aim for stableness

**NOTE:** This is first stage of the plugin, so is experimental

- [Terraform LSP](#terraform-lsp)
  * [Release](#release)
  * [Building](#building)
    + [Native](#native)
    + [Nixpkgs](#nixpkgs)
  * [Features](#features)
  * [Todo](#todo)
  * [LSP Support Table](#lsp-support-table)
  * [Supported Editors](#supported-editors)
  * [Bugs](#bugs)

[![asciicast](https://asciinema.org/a/245075.svg)](https://asciinema.org/a/245075)

## Release

Release can be found [here](https://github.com/juliosueiras/terraform-lsp/releases)

## Building

### Requirement

it will need Go 1.14+

### Native

1. Download a [release](https://github.com/juliosueiras/terraform-lsp/releases)
   or clone the repository
2. Run these commands from the `terraform-lsp` directory

```sh
GO111MODULE=on go mod download # Download the modules for the project
make      # Build the project. Alternatively run "go build"
make copy # Install the project
```

you may also specify a path to your preferred bin directory with the `DST` parameter

```sh
make copy DST="$your_preferred_bin_path" # Install the project
```

### Nixpkgs

- install nixpkgs
- `nix-build`

## Features

- Variables complex completion (infinite nesting type)
- Provider config completion
- Resource (with infinite block, looking at you Kubernetes provider ;) ) completion
- Data source completion
- Dynamic Error Checking (Terraform and HCL checks)
- Communication using provider binary (so it will support any provider as long as is built with terraform 0.12 SDK)
- Module nesting (infinite as well) variable completion

## Todo

All Todos are listed [here](Todo.md)

## LSP Support Table

| Feature            | Description                                         | Status                                                                                       |
|--------------------|-----------------------------------------------------|----------------------------------------------------------------------------------------------|
| completion         | Autocompletion                                      | Supported for Resources/Data Sources/Variables/Locals, need support for nested interpolation |
| publishDiagnostics | Error checking                                      | Supported, need to check for all possible errors                                             |
| hover              | Hover on function/variables to get result           | Need to Implement                                                                            |
| signatureHelp      | Get docs for resources/data sources name and params | Need to Implement                                                                            |
| declaration        | Go to Declaration                                   | Need to Implement                                                                            |
| references         | Find all references                                 | Need to Implement                                                                            |
| implementation     | Find all implementation                             | Need to Implement (not sure if is applicable)                                                |
| documentHighlight  | Resources/data sources/variables highlight          | Need to Implement                                                                            |
| documentSymbol     | Resources/data sources/variables symbols            | Need to Implement                                                                            |
| codeAction         | Refactoring actions                                 | Need to Implement                                                                            |
| codeLens           | VSCode's code lens                                  | Need to Implement                                                                            |
| formatting         | Formatting                                          | Need to Implement                                                                            |
| rename             | Rename action                                       | Need to Implement                                                                            |
| workspace          | Workspace support                                   | Need to Implement                                                                            |

## Supported Editors

| Editor             | Status    | Docs                             |
|--------------------|-----------|----------------------------------|
| Visual Studio Code | Supported | [Link](docs/editors/vscode.md)   |
| Atom               | Supported | [Link](docs/editors/atom.md)     |
| Vim                | Supported | [Link](docs/editors/vim.md)      |
| Sublime Text 3     | Supported | [Link](docs/editors/sublime3.md) |
| IntelliJ           | Supported | [Link](docs/editors/intellij.md) |
| Emacs              | Supported | [Link](docs/editors/emacs.md)    |

**NOTE:** Please create a issue for a editor that you want to test for

## Bugs
- Order of completion items
- Issue with block
- Issue on Terraform v0.15 due to different behaviour on where it store its provider json

## Credits
- LSP structure using [Sourcegraph's go-lsp](https://github.com/sourcegraph/go-lsp)
- JSONRPC2 using [JRPC2](https://bitbucket.org/creachadair/jrpc2)
- [provider communication](./tfstructs/provider.go) is mostly adapted from [tfschema](https://github.com/minamijoyo/tfschema)
- [awilkins](https://github.com/awilkins) for adding GitHub actions


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp?ref=badge_large)
