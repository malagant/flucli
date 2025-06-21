{
  description = "FluxCLI development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        isDarwin = pkgs.stdenv.isDarwin;
        isLinux = pkgs.stdenv.isLinux;
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go development (available on all platforms)
            go
            gopls
            gotools
            go-tools
            delve
            golangci-lint
            
            # Kubernetes tools (with platform checks)
            kubectl
            kind
            kubernetes-helm  # Use kubernetes-helm instead of helm for better platform support
            fluxcd
            
            # Git and general development tools
            git
            gnumake
            
            # Terminal and testing tools
            tmux
            curl
            jq
            yq-go
            
            # Documentation tools
            mdbook
          ] ++ pkgs.lib.optionals isLinux [
            # Linux-specific packages
            direnv
          ] ++ pkgs.lib.optionals isDarwin [
            # macOS-specific packages if needed
          ];

          shellHook = ''
            echo "🚀 FluxCLI development environment loaded"
            echo "📦 Go version: $(go version)"
            echo "🖥️  Platform: ${system}"
            echo "☸️  kubectl version: $(kubectl version --client --short 2>/dev/null || echo 'kubectl not configured')"
            echo "🌊 Flux version: $(flux version --client 2>/dev/null || echo 'flux not configured')"
            echo ""
            echo "💡 To activate this environment in your current shell, run:"
            echo "   nix develop"
            echo ""
            echo "Available commands:"
            echo "  go mod tidy   - Tidy Go modules"
            echo "  go build      - Build the application"
            echo "  go test       - Run tests"
            echo "  golangci-lint run - Run linter"
            echo ""
            
            # Set up Go environment
            export GOPATH=$HOME/go
            export PATH=$GOPATH/bin:$PATH
            
            # Set up project-specific environment
            export FLUXCLI_DEV=true
            export FLUXCLI_LOG_LEVEL=debug
            
            # Ensure Go is available
            if ! command -v go &> /dev/null; then
              echo "❌ Go not found in PATH"
              exit 1
            fi
          '';
        };

        # Package definition for building FluxCLI
        packages.default = pkgs.buildGoModule {
          pname = "fluxcli";
          version = "0.1.0-dev";
          
          src = ./.;
          
          vendorHash = null; # Will be filled when we have dependencies
          
          buildInputs = with pkgs; [
            git
          ];
          
          meta = with pkgs.lib; {
            description = "Terminal UI for FluxCD Multi-Cluster Management";
            homepage = "https://github.com/malagant/fluxcli";
            license = licenses.mit;
            maintainers = [ ];
            platforms = platforms.unix;
          };
        };
      });
}
