# Vim Support

Todo: add config snippets

- Should work with all LSP plugin on vim

### coc.nvim

- Install the [coc.nvim plugin](https://github.com/neoclide/coc.nvim)
- Add the following snippet to the `coc-setting.json` file (editable via `:CocConfig` in NeoVim)

```json
{
	"languageserver": {
		"terraform": {
			"command": "terraform-lsp",
			"filetypes": [
				"terraform",
				"tf"
			],
			"initializationOptions": {},
			"settings": {}
		}
	}
}
```
