# Terraform LSP

[![Gitter](https://badges.gitter.im/terraform-lsp/community.svg)](https://gitter.im/terraform-lsp/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
![terraform version](https://img.shields.io/badge/terraform-0.12.13-blue.svg)
[![Nix Build](https://img.shields.io/travis/com/juliosueiras/terraform-lsp.svg?logo=travis&label=Nix%20Build)](https://travis-ci.com/juliosueiras/terraform-lsp)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp?ref=badge_shield)

This is LSP (Language Server Protocol) for Terraform

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

it will need Go 1.12+

### Native

1. Download a [release](https://github.com/juliosueiras/terraform-lsp/releases)
   or clone the repository
2. Run these commands from the `terraform-lsp` directory

```sh
GO111MODULE=on go mod download # Download the modules for the project
make      # Build the project. Alternatively run "go build"
make copy # Install the project
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

## Credits
- LSP structure using [Sourcegraph's go-lsp](https://github.com/sourcegraph/go-lsp)
- JSONRPC2 using [JRPC2](https://bitbucket.org/creachadair/jrpc2)
- [provider communication](./tfstructs/provider.go) is mostly adapted from [tfschema](https://github.com/minamijoyo/tfschema)


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp?ref=badge_large)
