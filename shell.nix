with import <nixpkgs> {};

mkShell {
  buildInputs = [
    go_1_12
    dep
  ];

  shellHook = ''
    GOROOT=${pkgs.go_1_12}/share/go
  '';
}
