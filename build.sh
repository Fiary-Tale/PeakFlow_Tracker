#!/bin/bash

# Set the output directory
OUTPUT_DIR="./make"
mkdir -p $OUTPUT_DIR

# Define the targets
TARGETS=(
    "darwin arm64"
    "darwin amd64"
    "linux amd64"
    "linux arm64"
)

# Build the project for each target
for TARGET in "${TARGETS[@]}"; do
    OS=$(echo $TARGET | cut -d ' ' -f 1)
    ARCH=$(echo $TARGET | cut -d ' ' -f 2)

    OUTPUT_NAME="PeakFlow_Tracker-$OS-$ARCH"

    echo "Building for $OS $ARCH..."

    # Set the environment variables
    GOOS=$OS GOARCH=$ARCH go build -o "$OUTPUT_DIR/$OUTPUT_NAME"

    if [ $? -ne 0 ]; then
        echo "Failed to build for $OS $ARCH"
        exit 1
    fi

    echo "Successfully built $OUTPUT_NAME"
done

echo "All builds completed successfully."
