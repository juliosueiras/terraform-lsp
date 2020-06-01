with import <nixpkgs> {};

mkShell {
  buildInputs = [
    buildPackages.go_1_14
    dep
  ];

  shellHook = ''
    GOROOT=${buildPackages.go_1_14}/share/go
  '';
}
