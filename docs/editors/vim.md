# Vim Support

## CoC

`terraform-lsp` can be used with [CoC](https://github.com/neoclide/coc.nvim).

#### Requirements

Install the `terraform-lsp` binary and verify it's available in your `$PATH`.

#### Installation

Go to your `coc-settings.json` file (you can use the following command in Vim: `:CocConfig`).

Add an entry for the terraform language server.

```json
{
  "languageserver": {
    "terraform": {
      "command": "terraform-lsp",
      "filetypes": ["terraform"]
    }
  }
}
```
