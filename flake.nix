{
  description = "Touch Test development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let

      # Wrapper func to make supported archs a little easier
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs {
          inherit system;
        };
      });
    in
    {
      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          hardeningDisable = [ "fortify" ];

          packages = with pkgs; [
            # go (version is specified by overlay)
            go_1_22

            # Go Tools
            golangci-lint
            go-tools
            gotools
            gopls
            delve

            # Dev Tools
            go-task

            # gRPC deps
            buf
            protoc-gen-go
            protoc-gen-connect-go
            grpcurl
          ];
        };
      });
    };
}
