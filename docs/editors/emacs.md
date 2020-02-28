### Emacs Support

There are two options

1. Use the latest version of [emacs-lsp/lsp-mode](https://github.com/emacs-lsp/lsp-mode), it had added support for terraform-lsp

2. Work with [emacs-lsp/lsp-mode](https://github.com/emacs-lsp/lsp-mode) while still a little buggy
   ```lisp
   (add-to-list 'lsp-language-id-configuration '(terraform-mode . "terraform"))

   (lsp-register-client
    (make-lsp-client :new-connection (lsp-stdio-connection '("/path/to/terraform-lsp/terraform-lsp" "-enable-log-file"))
                     :major-modes '(terraform-mode)
                     :server-id 'terraform-ls))

   (add-hook 'terraform-mode-hook #'lsp)
   ```
