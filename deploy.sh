#!/bin/bash
set -e

# Determine the operating system and architecture
os_arch=$(uname -s)_$(uname -m)

# Map the operating system and architecture to Terraform target format
case $os_arch in
  Linux_x86_64)
    target="linux_amd64"
    ;;
  Linux_aarch64)
    target="linux_arm64"
    ;;
  Darwin_x86_64)
    target="darwin_amd64"
    ;;
  Darwin_arm64)
    target="darwin_arm64"
    ;;
  *)
    echo "Unsupported platform: $os_arch"
    exit 1
    ;;
esac

# Set environment variables
export host_name="antunovic.nz"
export namespace="synlestidae"
export type="mattr"
export version="0.0.1"
export target

# Path to the Terraform provider binary
provider_binary="terraform-provider-mattr"

# Remove any existing plugin files
rm -f ~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}/${provider_binary}

# Create the plugin directory and copy the provider binary
mkdir -p ~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}/
cp ${provider_binary} ~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}

# Logging
echo "Architecture determined: $os_arch"
echo "Terraform provider installed at: ~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}/${provider_binary}"

