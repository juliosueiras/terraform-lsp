# Terraform LSP

[![Gitter](https://badges.gitter.im/terraform-lsp/community.svg)](https://gitter.im/terraform-lsp/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
![terraform version](https://img.shields.io/badge/terraform-0.12.2-blue.svg)
[![Nix Build](https://img.shields.io/travis/com/juliosueiras/terraform-lsp.svg?logo=travis&label=Nix%20Build)](https://travis-ci.com/juliosueiras/terraform-lsp)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp?ref=badge_shield)

This is LSP(Language Server Protocol) for Terraform

**NOTE:** This is first stage of the plugin, so is experimental

- [Terraform LSP](#terraform-lsp)
  * [Building](#building)
    + [Native](#native)
    + [Nixpkgs](#nixpkgs)
  * [Currently Supported](#currently-supported)
  * [Todo-List (Main)](#todo-list--main-)
  * [Todo-List (Terraform)](#todo-list--terraform-)
  * [Todo-List (LSP)](#todo-list--lsp-)
  * [Editors Specific](#editors-specific)
    + [VS Code](#vs-code)
      - [Issues](#issues)
    + [Atom](#atom)
      - [Issues](#issues-1)
    + [Intellij](#intellij)
    + [Vim](#vim)
    + [Emacs](#emacs)
    + [Others](#others)
  * [Bugs](#bugs)

[![asciicast](https://asciinema.org/a/245075.svg)](https://asciinema.org/a/245075)

## Building

### Native

- download go modules `GO111MODULE=on go mod download`
- run make or `go build`

### Nixpkgs

- install nixpkgs
- `nix-build`

## Currently Supported

- Variables complex completion(infinite nesting type)
- Provider Config completion
- Resource(with infinite block, looking at you kubernetes provider ;) ) completion
- Data source completion
- Dynamic Error Checking(Terraform and HCL checks)
- Communication using provider binary(so it will support any provider as long as is built with terraform 0.12 sdk)
- Module nesting(inifinte as well) variable completion

## Todo-List (Main)

- A lot of code clean up(right now is mostly getting the feature done)
- CI/CD(Travis)
- Tests
- Getting Started Guide
- Add Gifs & Asciinema for Demo
- Cross Compile for Windows Binary

## Todo-List (Terraform)

- [X] Add Provider Listing(static list)
- [X] Provider Configs
- [X] Resources
  - [X] Provider attribute scope completion(ex. google vs google-beta)
- [X] Data Sources
- [ ] Backends (current focus)
- [ ] Provisioner (current focus)
- [ ] Interpolations
  - Complex nesting interpolations
- [ ] Modules
- [ ] Locals
- [ ] Outputs

- [ ] Variables
  - [ ] Map Interpolation with Object inside
  - [ ] Index Interpolation
  - [ ] List Interpolation with Object inside
  - [ ] Object Interpolation
  - [ ] Functions

- [X] Dynamic Block 
  - [X] For Each Block
    - [ ] Check for complex scenario

## Todo-List (LSP)

- [X] `initialize`
- [X] `textDocument/completion`
- [X] `textDocument/didChange`
- [X] `textDocument/didOpen`
- [X] `textDocument/publishDiagnostics`
- Current Plan: Implement all possible LSP features

## Editors Specific

### VS Code

- There is a seperate [vscode-languageclient-terraform](https://github.com/juliosueiras/vscode-languageclient-terraform)

#### Issues

- Need syntax and highlight etc (possible collab with vscode-terraform?)

### Atom

- There is a seperate [atom-terraform](https://github.com/juliosueiras/atom-terraform)

#### Issues

- Need configuration for linter API

### Intellij

- Work with intellij-lsp plugin(also work with intellij-hcl together)

### Vim

- Should work with all lsp plugin on vim

### Emacs

- Work with [emacs-lsp/lsp-mode](https://github.com/emacs-lsp/lsp-mode) while still a little buggy
   ```lisp
   (add-to-list 'lsp-language-id-configuration '(terraform-mode . "terraform"))

   (lsp-register-client
    (make-lsp-client :new-connection (lsp-stdio-connection '("/path/to/terraform-lsp/terraform-lsp" "-enable-log-file"))
                     :major-modes '(terraform-mode)
                     :server-id 'terraform-ls))

   (add-hook 'terraform-mode-hook #'lsp)
   ```

### Others

- It should work with any lsp client in theory

## Bugs
- Order of completion items
- Issue with block 

## Credits
- LSP structure using [Sourcegraph's go-lsp](https://github.com/sourcegraph/go-lsp)
- JSONRPC2 using [JRPC2](https://bitbucket.org/creachadair/jrpc2)
- [provider communication](./tfstructs/provider.go) is mostly adapted from [tfschema](https://github.com/minamijoyo/tfschema)


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjuliosueiras%2Fterraform-lsp?ref=badge_large)
