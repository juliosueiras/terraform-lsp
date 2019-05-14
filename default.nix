with import <nixpkgs> {};

buildGoModule rec {
  name = "terraform-lsp";
  version = "0.0.1";
  src = ./.;

  modSha256 = 
  if stdenv.isLinux then
  "04pvx225jffa9kbg6q86mkgnb744axr8nif93zwrmkqz7c6apl2h" 
  else if stdenv.isDarwin then
  "0qswp5lqzpfjvlnb47fkafgwgk8mj569y3ak1ng8xs1ajzyvm1k8"
  else
  "04pvx225jffa9kbg6q86mkgnb744axr8nif93zwrmkqz7c6apl2h"; 
  

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
