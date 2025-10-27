#!/bin/bash

# EasyGo Panel Build Script
set -e

PROJECT_NAME="easygo"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitHash=${COMMIT_HASH}"

echo "Building EasyGo Panel ${VERSION}"
echo "Build time: ${BUILD_TIME}"
echo "Commit: ${COMMIT_HASH}"

# Create build directory
mkdir -p build

# Build for current platform
echo "Building for current platform..."
CGO_ENABLED=1 go build -ldflags "${LDFLAGS}" -o build/${PROJECT_NAME} cmd/easygo/main.go

# Build for Linux AMD64 (most common server architecture)
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "${LDFLAGS}" -o build/${PROJECT_NAME}-linux-amd64 cmd/easygo/main.go

# Build for Linux ARM64 (for ARM servers)
echo "Building for Linux ARM64..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o build/${PROJECT_NAME}-linux-arm64 cmd/easygo/main.go

# Make binaries executable
chmod +x build/${PROJECT_NAME}*

# Create installation package
echo "Creating installation package..."
mkdir -p build/package
cp build/${PROJECT_NAME}-linux-amd64 build/package/${PROJECT_NAME}
cp README.md build/package/
cp install.sh build/package/

# Create systemd service file
cat > build/package/easygo.service << EOF
[Unit]
Description=EasyGo Panel
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/easygo
ExecStart=/opt/easygo/easygo web
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# Create installation script
cat > build/package/install.sh << 'EOF'
#!/bin/bash

# EasyGo Panel Installation Script
set -e

if [ "$EUID" -ne 0 ]; then
    echo "Please run as root"
    exit 1
fi

echo "Installing EasyGo Panel..."

# Create installation directory
mkdir -p /opt/easygo
cp easygo /opt/easygo/
chmod +x /opt/easygo/easygo

# Install systemd service
cp easygo.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable easygo

# Create symlink for CLI access
ln -sf /opt/easygo/easygo /usr/local/bin/easygo

echo "EasyGo Panel installed successfully!"
echo "Start the web panel: systemctl start easygo"
echo "Access CLI: easygo help"
echo "Web interface will be available at: http://your-server:8083"
EOF

chmod +x build/package/install.sh

# Create tarball
cd build
tar -czf ${PROJECT_NAME}-${VERSION}-linux-amd64.tar.gz package/
cd ..

echo "Build completed successfully!"
echo "Binaries:"
ls -la build/${PROJECT_NAME}*
echo ""
echo "Installation package: build/${PROJECT_NAME}-${VERSION}-linux-amd64.tar.gz"