# Sublime Text 3 Support

Add this to the `clients` settings for [tomv564's LSP](https://github.com/tomv564/LSP) also make sure to have `terraform` syntax plugin for sublime text

```json
{
  "clients":
  {
    "terraform":
    {
      "command":
      [
        "terraform-lsp",
        "-enable-log-file",
        "-log-location",
        "/tmp/"
      ],
      "enabled": true,
      "scopes": ["source.terraform"],
      "syntaxes":  ["Packages/Terraform/Terraform.sublime-syntax"],
      "languageId": "terraform"
    }
  }
}
```
