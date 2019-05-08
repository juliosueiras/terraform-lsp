# Terraform LSP

This is LSP(Language Server Protocol) for Terraform

**NOTE:** This is first stage of the plugin, so is experimental

- [Terraform LSP](#terraform-lsp)
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

## Todo-List (Terraform)

- [X] Provider Configs
- [X] Resources
- [X] Data Sources
- [ ] Backends
- [ ] Provisioner
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

- [ ] Dynamic Block
  - Block

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

- Should work with emacs-lsp(need confirmation)

### Others

- It should work with any lsp client in theory

## Bugs
- Order of completion items
- Issue with block 

## Credits
- LSP structure using [Sourcegraph's go-lsp](https://github.com/sourcegraph/go-lsp)
- JSONRPC2 using [JRPC2](https://bitbucket.org/creachadair/jrpc2)
- (provider communication)[./tfstructs/provider.go] is mostly adapted from [tfschema](https://github.com/minamijoyo/tfschema)
