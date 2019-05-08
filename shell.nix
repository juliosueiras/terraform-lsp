with import <nixpkgs> {};

mkShell {
  buildInputs = [
    go_1_11
    dep
  ];

  shellHook = ''
    GOROOT=${pkgs.go}/share/go
  '';
}
