#!/bin/bash

# Build script for vault-hub base Docker image with multi-platform support

set -e

# Configuration
IMAGE_NAME="vault-hub-base"
TAG="latest"
REGISTRY=""  # Set this to your registry URL if needed (e.g., registry.gitlab.com/your-group/your-project)
PLATFORMS="linux/amd64,linux/arm64"  # Default platforms
USE_BUILDX=true

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --tag)
            TAG="$2"
            shift 2
            ;;
        --registry)
            REGISTRY="$2"
            shift 2
            ;;
        --platforms)
            PLATFORMS="$2"
            shift 2
            ;;
        --single-platform)
            USE_BUILDX=false
            shift
            ;;
        --push)
            PUSH=true
            shift
            ;;
        --help)
            echo "Usage: $0 [--tag TAG] [--registry REGISTRY] [--platforms PLATFORMS] [--single-platform] [--push]"
            echo "  --tag TAG           Set the image tag (default: latest)"
            echo "  --registry REG      Set the registry URL"
            echo "  --platforms PLAT    Set target platforms (default: linux/amd64,linux/arm64)"
            echo "  --single-platform   Use regular docker build instead of buildx"
            echo "  --push             Push the image after building"
            echo ""
            echo "⚠️  Important: Multi-platform builds MUST be pushed to a registry."
            echo "   Local loading only works with single-platform builds."
            echo ""
            echo "Examples:"
            echo "  $0 --platforms linux/amd64,linux/arm64 --push"
            echo "  $0 --platforms linux/amd64 --single-platform"
            echo "  $0 --tag v1.0.0 --registry registry.gitlab.com/mygroup/myproject"
            echo "  $0 --single-platform  # For local development"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Construct full image name
if [ -n "$REGISTRY" ]; then
    FULL_IMAGE_NAME="${REGISTRY}/${IMAGE_NAME}:${TAG}"
else
    FULL_IMAGE_NAME="${IMAGE_NAME}:${TAG}"
fi

echo "Building Docker image: $FULL_IMAGE_NAME"

if [ "$USE_BUILDX" = true ]; then
    echo "Using Docker Buildx for multi-platform build"
    echo "Target platforms: $PLATFORMS"
    
    # Check if buildx is available
    if ! docker buildx version >/dev/null 2>&1; then
        echo "❌ Docker Buildx is not available. Please install Docker Buildx or use --single-platform flag."
        exit 1
    fi
    
    # Create and use a new builder instance if it doesn't exist
    BUILDER_NAME="vault-hub-builder"
    if ! docker buildx inspect "$BUILDER_NAME" >/dev/null 2>&1; then
        echo "Creating new buildx builder: $BUILDER_NAME"
        docker buildx create --name "$BUILDER_NAME" --driver docker-container --bootstrap
    fi
    
    echo "Using buildx builder: $BUILDER_NAME"
    docker buildx use "$BUILDER_NAME"
    
    # Build command with buildx
    BUILD_CMD="docker buildx build"
    BUILD_CMD="$BUILD_CMD --platform $PLATFORMS"
    BUILD_CMD="$BUILD_CMD -f Dockerfile-base"
    BUILD_CMD="$BUILD_CMD -t $FULL_IMAGE_NAME"
    
    # Check if building for multiple platforms
    PLATFORM_COUNT=$(echo "$PLATFORMS" | tr ',' '\n' | wc -l)
    
    # Add push flag if requested
    if [ "$PUSH" = true ]; then
        BUILD_CMD="$BUILD_CMD --push"
        echo "Will push to registry after building"
    elif [ "$PLATFORM_COUNT" -gt 1 ]; then
        echo "⚠️  Multi-platform builds cannot be loaded to local Docker daemon."
        echo "   Multi-platform images must be pushed to a registry."
        echo ""
        echo "Options:"
        echo "  1. Add --push flag to push to registry: $FULL_IMAGE_NAME"
        echo "  2. Use --single-platform for local development"
        echo "  3. Use --platforms with a single platform (e.g., --platforms linux/amd64)"
        echo ""
        read -p "Do you want to push to registry instead? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            BUILD_CMD="$BUILD_CMD --push"
            echo "Will push to registry after building"
        else
            echo "❌ Build cancelled. Please use one of the suggested options."
            exit 1
        fi
    else
        BUILD_CMD="$BUILD_CMD --load"
        echo "Will load image to local Docker daemon"
    fi
    
    BUILD_CMD="$BUILD_CMD ."
    
    echo "Executing: $BUILD_CMD"
    eval "$BUILD_CMD"
    
else
    echo "Using regular Docker build (single platform)"
    
    # Build the image with regular docker build
    docker build -f Dockerfile-base -t "$FULL_IMAGE_NAME" .
    
    # Push if requested
    if [ "$PUSH" = true ]; then
        echo "Pushing image to registry..."
        docker push "$FULL_IMAGE_NAME"
        echo "✅ Successfully pushed: $FULL_IMAGE_NAME"
    fi
fi

echo "✅ Successfully built: $FULL_IMAGE_NAME"

echo "🎉 Done!"
echo ""
if [ "$USE_BUILDX" = true ]; then
    echo "Multi-platform image built for: $PLATFORMS"
else
    echo "Single-platform image built for current architecture"
fi
echo ""
echo "To use this image in your GitLab CI, update your .gitlab-ci.yml:"
echo "  image: $FULL_IMAGE_NAME" 