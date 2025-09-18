#!/bin/bash

# Glens Setup Script
# This script sets up the development environment using micromamba

set -e  # Exit on any error

echo "🚀 Glens Setup"
echo "=============================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [[ ! -f "environment.yml" ]]; then
    print_error "environment.yml not found. Please run this script from the glens directory."
    exit 1
fi

# Check for micromamba installation
print_status "Checking for micromamba installation..."
if ! command -v micromamba &> /dev/null; then
    print_warning "micromamba not found. Installing micromamba..."

    # Detect OS
    OS=$(uname -s)
    ARCH=$(uname -m)

    case "$OS" in
        Linux*)
            if [[ "$ARCH" == "x86_64" ]]; then
                MAMBA_URL="https://micro.mamba.pm/api/micromamba/linux-64/latest"
            elif [[ "$ARCH" == "aarch64" ]] || [[ "$ARCH" == "arm64" ]]; then
                MAMBA_URL="https://micro.mamba.pm/api/micromamba/linux-aarch64/latest"
            else
                print_error "Unsupported architecture: $ARCH"
                exit 1
            fi
            ;;
        Darwin*)
            if [[ "$ARCH" == "x86_64" ]]; then
                MAMBA_URL="https://micro.mamba.pm/api/micromamba/osx-64/latest"
            elif [[ "$ARCH" == "arm64" ]]; then
                MAMBA_URL="https://micro.mamba.pm/api/micromamba/osx-arm64/latest"
            else
                print_error "Unsupported architecture: $ARCH"
                exit 1
            fi
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            print_error "Please install micromamba manually from: https://mamba.readthedocs.io/en/latest/installation.html"
            exit 1
            ;;
    esac

    # Create micromamba directory
    MAMBA_ROOT_PREFIX="${HOME}/micromamba"
    mkdir -p "$MAMBA_ROOT_PREFIX/bin"

    # Download and install micromamba
    print_status "Downloading micromamba from $MAMBA_URL..."
    curl -Ls "$MAMBA_URL" | tar -xvj -C "$MAMBA_ROOT_PREFIX" bin/micromamba

    # Add to PATH for this session
    export PATH="$MAMBA_ROOT_PREFIX/bin:$PATH"

    # Initialize micromamba
    print_status "Initializing micromamba..."
    micromamba shell init --shell bash --prefix="$MAMBA_ROOT_PREFIX"

    print_success "micromamba installed successfully!"
    print_warning "Please restart your shell or run: source ~/.bashrc"
    print_warning "Then run this setup script again."

    # Add to shell profile
    SHELL_PROFILE=""
    if [[ -f "$HOME/.bashrc" ]]; then
        SHELL_PROFILE="$HOME/.bashrc"
    elif [[ -f "$HOME/.zshrc" ]]; then
        SHELL_PROFILE="$HOME/.zshrc"
    elif [[ -f "$HOME/.profile" ]]; then
        SHELL_PROFILE="$HOME/.profile"
    fi

    if [[ -n "$SHELL_PROFILE" ]]; then
        echo "export PATH=\"$MAMBA_ROOT_PREFIX/bin:\$PATH\"" >> "$SHELL_PROFILE"
        print_status "Added micromamba to $SHELL_PROFILE"
    fi

    exit 0
else
    print_success "micromamba found: $(which micromamba)"
fi

# Check micromamba version
MAMBA_VERSION=$(micromamba --version 2>/dev/null | head -n1 || echo "unknown")
print_status "micromamba version: $MAMBA_VERSION"

# Create/update environment
print_status "Setting up Glens environment..."
make env

# Verify environment
print_status "Verifying environment setup..."
if micromamba info -e | grep -q "glens"; then
    print_success "Environment 'glens' created successfully!"
else
    print_error "Failed to create environment"
    exit 1
fi

# Check Go version in environment
print_status "Checking Go version in environment..."
GO_VERSION=$(micromamba run -n glens go version 2>/dev/null || echo "Go not found")
print_status "Go version: $GO_VERSION"

# Setup development tools
print_status "Installing development tools..."
make setup

# Run initial tests to verify everything works
print_status "Running initial verification tests..."
if make test; then
    print_success "All tests passed!"
else
    print_warning "Some tests failed, but the environment is set up"
fi

# Build the binary
print_status "Building Glens..."
if make build; then
    print_success "Build completed successfully!"
    BUILD_PATH="$(pwd)/build/glens"
    print_status "Binary location: $BUILD_PATH"
else
    print_error "Build failed"
    exit 1
fi

# Final instructions
echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "Next steps:"
echo "  1. Activate the environment: micromamba activate glens"
echo "  2. Configure your API keys in configs/config.yaml"
echo "  3. Run a test: make run"
echo ""
echo "Available commands:"
echo "  make help     - Show all available commands"
echo "  make shell    - Enter the micromamba environment"
echo "  make run      - Run with example OpenAPI spec"
echo "  make test     - Run the test suite"
echo ""
echo "For GitHub integration, set up these secrets in your repository:"
echo "  - OPENAI_API_KEY"
echo "  - ANTHROPIC_API_KEY"
echo "  - GOOGLE_API_KEY"
echo "  - GOOGLE_PROJECT_ID"
echo ""
print_success "Happy testing! 🚀"
