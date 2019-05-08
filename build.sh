docker run --rm --privileged \
  -v $PWD:/go/src/github.com/juliosueiras/terraform-lsp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -w /go/src/github.com/juliosueiras/terraform-lsp \
  mailchain/goreleaser-xcgo goreleaser --snapshot --rm-dist

