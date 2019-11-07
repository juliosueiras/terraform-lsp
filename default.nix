with import <nixpkgs> {};

buildGoModule rec {

  inherit (go_1_12);
  
  name = "terraform-lsp";
  version = "0.0.9";
  src = ./.;

  modSha256 = "1p3g6h0ygh2jfmyqq77rja9ajxqlvv6cj0rdim5wyal1sz5n0sjx"; 

  buildPhase = ''
    runHook preBuild
    runHook renameImports
    go install -ldflags="-s -w -X main.Version=${version} -X main.GitCommit='omitted' -X main.Date='omitted'"
    runHook postBuild
  '';

  goPackagePath = "github.com/juliosueiras/terraform-lsp";
  subPackages = [ "." ];

  meta = with stdenv.lib; {
    description = "Language Server Protocol for Terraform";
    homepage = https://github.com/juliosueiras/terraform-lsp;
    license = licenses.mit;
    maintainers = with maintainers; [ juliosueiras ];
    platforms = platforms.all;
  };
}
