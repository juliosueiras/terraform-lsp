with import <nixpkgs> {};

buildGo112Module rec {

  name = "terraform-lsp";
  version = "0.0.9";
  src = ./.;

  modSha256 = "06427z3lxvqv0vvi2gwni4416aw73qhd0v6nqb86qhkww7zimpip"; 

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
