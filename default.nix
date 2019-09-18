with import <nixpkgs> {};

buildGoModule rec {
  name = "terraform-lsp";
  version = "0.0.6";
  src = ./.;

  modSha256 = "1mwhpabi7n0flw83q6850gg85c0v8b6hvyxjjglzzixjn58nmjxb"; 

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
