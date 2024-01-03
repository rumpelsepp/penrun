# SPDX-FileCopyrightText: Stefan Tatschner
#
# SPDX-License-Identifier: CC0-1.0

{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs";
  };

  outputs = { self, nixpkgs }:
    with import nixpkgs { system = "x86_64-linux"; };
    let pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in {
      devShell.x86_64-linux = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          bats
          reuse
          shellcheck
          shfmt
          nodePackages_latest.bash-language-server
        ];
        shellHook = ''
          LD_LIBRARY_PATH=${pkgs.lib.makeLibraryPath [stdenv.cc.cc]}
        '';
      };
      packages.x86_64-linux.default =
        pkgs.buildGoModule {
          pname = "penrun";
          src = self;
          version = self.lastModifiedDate;
          buildInputs = [ pkgs.go ];
          vendorHash = "sha256-IfJQHzKlw9Xqb9nEdQ9Z9rVD7mWG2FZvdN5JYhvTURo=";

          meta = with lib; {
            mainProgram = "penrun";
            platforms = platforms.unix;
          };
        };
      formatter.x86_64-linux = pkgs.nixpkgs-fmt;
    };
}
