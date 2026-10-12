{
  description = "FractalForge - Terminal Fractal Renderer";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gopls
        ];

        shellHook = ''
          export GOPATH="$PWD/.go"
          export GOCACHE="$PWD/.cache/go-build"
          export GOTOOLCHAIN=local

          mkdir -p "$GOPATH" "$GOCACHE"

          echo "FractalForge development environment"
        '';
      };
    };
}
