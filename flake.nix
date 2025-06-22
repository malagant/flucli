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
            kubecolor  # Colorized kubectl output (user has k=kubecolor alias)
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
            
            # Shell enhancement tools (respecting user's existing setup)
            starship          # Modern shell prompt (user already uses this)
            fzf               # Fuzzy finder (user already uses this)
            eza               # Modern ls replacement (user has ls=eza alias)
            direnv            # Environment management (user already uses this)
            atuin             # Shell history management (user already uses this)
            
            # Zsh plugins and enhancements
            zsh-autosuggestions    # Fish-like autosuggestions
            zsh-syntax-highlighting # Syntax highlighting
            zsh-completions        # Additional completions
            
            # Development utilities
            bat               # Better cat with syntax highlighting
            ripgrep           # Fast grep replacement
            fd                # Better find
            tree              # Directory tree display
            
          ] ++ pkgs.lib.optionals isLinux [
            # Linux-specific packages
          ] ++ pkgs.lib.optionals isDarwin [
            # macOS-specific packages if needed
          ];

          shellHook = ''
            # Only show the welcome message if not already shown in this session
            # This prevents repeated messages when sourcing .zshrc
            if [ "$FLUXCLI_WELCOME_SHOWN" != "1" ]; then
              echo "🚀 FluxCLI development environment loaded"
              echo "📦 Go version: $(go version)"
              echo "🖥️  Platform: ${system}"
              echo "☸️  kubectl version: $(kubectl version --client --short 2>/dev/null || echo 'kubectl not configured')"
              echo "🌊 Flux version: $(flux version --client 2>/dev/null || echo 'flux not configured')"
              echo "⭐ Starship prompt: $(starship --version 2>/dev/null || echo 'not available')"
              echo ""
              echo "Available commands:"
              echo "  ./dev.sh build    - Build FluxCLI"
              echo "  ./dev.sh test     - Run tests"
              echo "  ./dev.sh lint     - Run linter"
              echo "  make help         - Show all available make targets"
              echo ""
              echo "FluxCLI-specific aliases:"
              echo "  fk                - FluxCLI with kubectl context"
              echo "  fks               - FluxCLI show all resources"
              echo "  fkr               - FluxCLI reconcile"
              echo ""
              export FLUXCLI_WELCOME_SHOWN=1
            fi
            
            # Preserve existing GOPATH if set by user, otherwise set default
            if [ -z "$GOPATH" ]; then
              export GOPATH=$HOME/go
            fi
            
            # Preserve existing PATH and only prepend if not already present
            if [[ ":$PATH:" != *":$GOPATH/bin:"* ]]; then
              export PATH="$GOPATH/bin:$PATH"
            fi
            
            # Set up project-specific environment variables
            export FLUXCLI_DEV=true
            # Use existing FLUXCLI_LOG_LEVEL if set, otherwise default to debug
            export FLUXCLI_LOG_LEVEL=''${FLUXCLI_LOG_LEVEL:-debug}
            
            # Enable zsh enhancements if in zsh and not already loaded by user
            if [ -n "$ZSH_VERSION" ]; then
              # Only load our enhancements if not already available
              
              # Autosuggestions (respect user's existing config)
              if [ -f "${pkgs.zsh-autosuggestions}/share/zsh-autosuggestions/zsh-autosuggestions.zsh" ] && ! typeset -f _zsh_autosuggest_start > /dev/null; then
                source "${pkgs.zsh-autosuggestions}/share/zsh-autosuggestions/zsh-autosuggestions.zsh"
                export ZSH_AUTOSUGGEST_STRATEGY=(history completion)
                export ZSH_AUTOSUGGEST_BUFFER_MAX_SIZE=20
              fi
              
              # Syntax highlighting (load last)
              if [ -f "${pkgs.zsh-syntax-highlighting}/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh" ] && ! typeset -f _zsh_highlight > /dev/null; then
                source "${pkgs.zsh-syntax-highlighting}/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
              fi
              
              # Additional completions
              if [ -d "${pkgs.zsh-completions}/share/zsh/site-functions" ]; then
                fpath+="${pkgs.zsh-completions}/share/zsh/site-functions"
                # Only run compinit if it hasn't been run recently
                autoload -U compinit
                # Simple compinit without complex glob patterns that need escaping
                compinit -C
              fi
            fi
            
            # FluxCLI-specific aliases and functions
            # Only add if not already defined by user
            if ! alias fk > /dev/null 2>&1; then
              alias fk='./fluxcli'
            fi
            
            if ! alias fks > /dev/null 2>&1; then
              alias fks='./fluxcli --all-namespaces'
            fi
            
            if ! alias fkr > /dev/null 2>&1; then
              alias fkr='flux reconcile source git'
            fi
            
            # FluxCLI development helpers
            if ! alias fluxdev > /dev/null 2>&1; then
              alias fluxdev='./dev.sh build && ./fluxcli'
            fi
            
            if ! alias fluxtest > /dev/null 2>&1; then
              alias fluxtest='./dev.sh test && ./dev.sh lint'
            fi
            
            # Enhanced kubectl aliases for FluxCD (only if not already defined)
            if ! alias kgf > /dev/null 2>&1; then
              alias kgf='kubectl get --all-namespaces gitrepositories,helmrepositories,kustomizations,helmreleases'
            fi
            
            if ! alias kdf > /dev/null 2>&1; then
              alias kdf='kubectl describe'
            fi
            
            # Function to quickly test FluxCLI against different clusters
            if ! typeset -f fluxcli-test > /dev/null; then
              fluxcli-test() {
                local context=''${1:-$(kubectl config current-context)}
                echo "Testing FluxCLI against context: $context"
                ./dev.sh build && kubectl config use-context "$context" && ./fluxcli
              }
            fi
            
            # Function to quickly switch between flux contexts
            if ! typeset -f flux-ctx > /dev/null; then
              flux-ctx() {
                kubectl config use-context "$1"
                echo "Switched to context: $1"
                kubectl cluster-info
              }
            fi
            
            # Source user's .zshrc if it exists and we're in zsh
            # This ensures user's aliases, functions, and customizations are preserved
            if [ -n "$ZSH_VERSION" ] && [ -f "$HOME/.zshrc" ] && [ -z "$FLUXCLI_ZSHRC_SOURCED" ]; then
              echo "📝 Sourcing your existing .zshrc configuration..."
              export FLUXCLI_ZSHRC_SOURCED=1
              # Source the .zshrc but ignore any errors to prevent breaking the nix shell
              source "$HOME/.zshrc" 2>/dev/null || true
            fi
            
            # Verify Go is available (but don't exit, just warn)
            if ! command -v go &> /dev/null; then
              echo "⚠️  Warning: Go not found in PATH. Some development features may not work."
            fi
            
            # Initialize starship if available and not already initialized
            if command -v starship &> /dev/null && [ -z "$STARSHIP_SESSION_KEY" ] && [ "$TERM" != "dumb" ]; then
              eval "$(starship init zsh)"
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
