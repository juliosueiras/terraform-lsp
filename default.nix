with import <nixpkgs> {};

buildGoModule rec {
  name = "terraform-lsp";
  version = "0.0.3";
  src = ./.;

  modSha256 = "043916n6dpp3qkljqdnaj6bx6zmnmm7csbbjs8gwbx7kg0hw4fss"; 

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
