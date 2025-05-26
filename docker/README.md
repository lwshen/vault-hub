# Docker Build Scripts

This directory contains scripts for building Docker images for the vault-hub project.

## Multi-Platform Base Image Build

The `build-base-image.sh` script supports building Docker images for multiple architectures using Docker Buildx.

### Prerequisites

1. **Docker Buildx**: Make sure Docker Buildx is installed and available
   ```bash
   docker buildx version
   ```

2. **QEMU (for cross-platform builds)**: Install QEMU for emulation support
   ```bash
   # On Ubuntu/Debian
   sudo apt-get install qemu-user-static

   # On macOS with Homebrew
   brew install qemu

   # Enable experimental features in Docker (if needed)
   export DOCKER_CLI_EXPERIMENTAL=enabled
   ```

### Usage

#### Basic Multi-Platform Build
Build for both AMD64 and ARM64 architectures:
```bash
./build-base-image.sh
```

#### Custom Platforms
Specify custom target platforms:
```bash
./build-base-image.sh --platforms linux/amd64,linux/arm64,linux/arm/v7
```

#### Single Platform Build
Use regular Docker build for current architecture only:
```bash
./build-base-image.sh --single-platform
```

#### Build and Push to Registry
Build and push to a container registry:
```bash
./build-base-image.sh --registry registry.gitlab.com/your-group/your-project --push
```

#### Custom Tag
Build with a specific tag:
```bash
./build-base-image.sh --tag v1.0.0
```

### Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `--tag TAG` | Set the image tag | `latest` |
| `--registry REG` | Set the registry URL | (none) |
| `--platforms PLAT` | Set target platforms | `linux/amd64,linux/arm64` |
| `--single-platform` | Use regular docker build | (disabled) |
| `--push` | Push image after building | (disabled) |
| `--help` | Show help message | - |

### Examples

1. **Development build** (local, single platform):
   ```bash
   ./build-base-image.sh --single-platform
   ```

2. **Production build** (multi-platform, with push):
   ```bash
   ./build-base-image.sh --tag v1.2.3 --registry registry.gitlab.com/mygroup/vault-hub --push
   ```

3. **ARM-only build** for Raspberry Pi:
   ```bash
   ./build-base-image.sh --platforms linux/arm64 --tag rpi-latest
   ```

4. **Testing specific platforms**:
   ```bash
   ./build-base-image.sh --platforms linux/amd64,linux/arm/v7 --tag test
   ```

### Troubleshooting

#### Multi-Platform Build Cannot Load to Local Docker
**Error**: `docker exporter does not currently support exporting manifest lists`

**Solution**: Multi-platform builds create manifest lists that cannot be loaded to local Docker daemon.

Options:
1. **Push to registry** (recommended for multi-platform):
   ```bash
   ./build-base-image.sh --registry your-registry --push
   ```

2. **Use single platform** for local development:
   ```bash
   ./build-base-image.sh --single-platform
   ```

3. **Build for current platform only**:
   ```bash
   ./build-base-image.sh --platforms linux/amd64
   ```

#### Buildx Not Available
If you get an error about Buildx not being available:
```bash
# Install Docker Buildx plugin
docker buildx install

# Or use single-platform build
./build-base-image.sh --single-platform
```

#### QEMU Issues
If cross-platform builds fail:
```bash
# Install QEMU static binaries
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes

# Verify QEMU is working
docker buildx ls
```

#### Builder Instance Issues
If you encounter builder issues:
```bash
# Remove existing builder
docker buildx rm vault-hub-builder

# The script will create a new one on next run
./build-base-image.sh
```

### CI/CD Integration

For GitLab CI, you can use the built image in your `.gitlab-ci.yml`:

```yaml
image: registry.gitlab.com/your-group/vault-hub/vault-hub-base:latest

stages:
  - build
  - test

build:
  stage: build
  script:
    - go build ./...

test:
  stage: test
  script:
    - go test ./...
```

### Performance Notes

- **Multi-platform builds** take longer due to emulation overhead
- **ARM builds on x86** are significantly slower than native builds
- Consider using **native runners** for ARM builds in production CI/CD
- **Single-platform builds** are faster for development and testing 