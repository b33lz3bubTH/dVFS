{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go         # Golang compiler + toolchain
    pkgs.sqlite     # SQLite CLI + development headers
  ];

  shellHook = ''
    echo "🚀 Entered Go + SQLite development shell"
    echo "Go version: $(go version)"
    echo "SQLite version: $(sqlite3 --version)"
  '';
}
