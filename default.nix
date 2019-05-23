with import <nixpkgs> {};

buildGoModule rec {
  name = "terraform-lsp";
  version = "0.0.3";
  src = ./.;

  modSha256 = "0kpw4jvw1cy52v1fsd756mj6p7rq0pw6kf0pqrfnx1gbxj56gmll"; 
  

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
