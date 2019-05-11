#!/bin/sh
sed -r -e "s;- \\[([^]]*)\\]\\(/docs/providers/([^/]*)/index.html\\);\"\\2\": TerraformProvider{ Name: \"\\1\", Type: \"$2\" },;" $1
