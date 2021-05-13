{ pkgs ? import
  (fetchTarball "https://github.com/nixos/nixpkgs/archive/e10c65cdb35.tar.gz")
  { } }:
with pkgs;

buildGoPackage rec {

  name = "terraform-lsp";
  version = "0.0.12";
  src = ./.;

  buildFlagsArray = [
    ''-ldflags="-s -w -X main.Version=${version} -X main.GitCommit='omitted' -X main.Date='omitted'"''
  ];

  goPackagePath = "github.com/juliosueiras/terraform-lsp";

  meta = with stdenv.lib; {
    description = "Language Server Protocol for Terraform";
    homepage = "https://github.com/juliosueiras/terraform-lsp";
    license = licenses.mit;
    maintainers = with maintainers; [ juliosueiras ];
    platforms = platforms.all;
  };
}
