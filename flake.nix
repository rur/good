{
  description = "Good development environment";

  inputs.nixpkgs.url = "nixpkgs/24.11";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (
      system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
          {
            devShell = pkgs.mkShell {
              shellHook = ''
                set -a
                source .env 2>/dev/null || echo "no .env file found"
                set +a
              '';
              buildInputs = with pkgs; [
                go
                gopls
                gops
                gotests
                delve
                go-tools
                errcheck
                reftools
                revive
                gomodifytags
                gotags
                impl
                go-motion
                iferr
              ];
            };
          }
    );
}
