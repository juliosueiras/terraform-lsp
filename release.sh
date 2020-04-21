docker run --rm --privileged \
  -e GITHUB_TOKEN=$GITHUB_TOKEN \
  -v $PWD:/go/src/github.com/juliosueiras/terraform-lsp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -w /go/src/github.com/juliosueiras/terraform-lsp \
  mailchain/goreleaser-xcgo release --rm-dist
