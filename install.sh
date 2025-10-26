#!/bin/bash

# EasyGo Panel Installation Script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
INSTALL_DIR="/opt/easygo"
SERVICE_NAME="easygo"
USER="root"

echo -e "${GREEN}EasyGo Panel Installation Script${NC}"
echo "=================================="

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Error: Please run this script as root${NC}"
    exit 1
fi

# Detect OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$NAME
    VERSION=$VERSION_ID
else
    echo -e "${RED}Error: Cannot detect operating system${NC}"
    exit 1
fi

echo "Detected OS: $OS $VERSION"

# Check if EasyGo is already installed
if [ -f "$INSTALL_DIR/easygo" ]; then
    echo -e "${YELLOW}EasyGo Panel is already installed. Updating...${NC}"
    systemctl stop $SERVICE_NAME 2>/dev/null || true
fi

# Create installation directory
echo "Creating installation directory..."
mkdir -p $INSTALL_DIR

# Download or copy binary (assuming binary is in current directory)
echo "Installing EasyGo Panel binary..."
if [ -f "easygo" ]; then
    cp easygo $INSTALL_DIR/
    chmod +x $INSTALL_DIR/easygo
else
    echo -e "${RED}Error: easygo binary not found in current directory${NC}"
    exit 1
fi

# Create systemd service file
echo "Creating systemd service..."
cat > /etc/systemd/system/$SERVICE_NAME.service << EOF
[Unit]
Description=EasyGo Panel - Web Server Management Panel
After=network.target
Wants=network.target

[Service]
Type=simple
User=$USER
Group=$USER
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/easygo web --port 8080 --host 0.0.0.0
ExecReload=/bin/kill -HUP \$MAINPID
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=easygo

# Security settings
NoNewPrivileges=false
PrivateTmp=true
ProtectSystem=strict
ReadWritePaths=$INSTALL_DIR /var/log /var/lib /etc
ProtectHome=true

[Install]
WantedBy=multi-user.target
EOF

# Create CLI symlink
echo "Creating CLI symlink..."
ln -sf $INSTALL_DIR/easygo /usr/local/bin/easygo

# Install dependencies based on OS
echo "Installing dependencies..."
case "$OS" in
    "Ubuntu"|"Debian GNU/Linux")
        apt update
        apt install -y curl wget gnupg2 software-properties-common apt-transport-https ca-certificates
        ;;
    "CentOS Linux"|"Red Hat Enterprise Linux"|"AlmaLinux"|"Rocky Linux")
        if command -v dnf &> /dev/null; then
            dnf install -y curl wget gnupg2
        else
            yum install -y curl wget gnupg2
        fi
        ;;
    *)
        echo -e "${YELLOW}Warning: Unknown OS. Some dependencies might need to be installed manually.${NC}"
        ;;
esac

# Set up PAM for authentication
echo "Configuring PAM authentication..."
if [ ! -f /etc/pam.d/login ]; then
    echo -e "${YELLOW}Warning: PAM login configuration not found. Authentication might not work properly.${NC}"
fi

# Create log directory
mkdir -p /var/log/easygo
chown $USER:$USER /var/log/easygo

# Reload systemd and enable service
echo "Enabling and starting EasyGo Panel service..."
systemctl daemon-reload
systemctl enable $SERVICE_NAME
systemctl start $SERVICE_NAME

# Check service status
sleep 2
if systemctl is-active --quiet $SERVICE_NAME; then
    echo -e "${GREEN}✓ EasyGo Panel service started successfully${NC}"
else
    echo -e "${RED}✗ Failed to start EasyGo Panel service${NC}"
    echo "Check logs with: journalctl -u $SERVICE_NAME -f"
    exit 1
fi

# Configure firewall (if available)
echo "Configuring firewall..."
if command -v ufw &> /dev/null; then
    ufw allow 8080/tcp
    echo "Firewall rule added for port 8080"
elif command -v firewall-cmd &> /dev/null; then
    firewall-cmd --permanent --add-port=8080/tcp
    firewall-cmd --reload
    echo "Firewall rule added for port 8080"
elif command -v iptables &> /dev/null; then
    iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
    echo "Iptables rule added for port 8080"
fi

# Display completion message
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}EasyGo Panel Installation Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "🌐 Web Interface: http://$(hostname -I | awk '{print $1}'):8080"
echo "📱 CLI Usage: easygo help"
echo "🔧 Service Management:"
echo "   - Start:   systemctl start $SERVICE_NAME"
echo "   - Stop:    systemctl stop $SERVICE_NAME"
echo "   - Restart: systemctl restart $SERVICE_NAME"
echo "   - Status:  systemctl status $SERVICE_NAME"
echo "   - Logs:    journalctl -u $SERVICE_NAME -f"
echo ""
echo "📋 Next Steps:"
echo "   1. Access the web interface using your system credentials"
echo "   2. Configure your first web server: easygo apache install"
echo "   3. Add a domain: easygo apache vhost example.com /var/www/example.com"
echo ""
echo -e "${YELLOW}Note: Use your system user credentials to log in to the web panel${NC}"

# Create uninstall script
cat > $INSTALL_DIR/uninstall.sh << 'EOF'
#!/bin/bash
echo "Uninstalling EasyGo Panel..."
systemctl stop easygo 2>/dev/null || true
systemctl disable easygo 2>/dev/null || true
rm -f /etc/systemd/system/easygo.service
rm -f /usr/local/bin/easygo
rm -rf /opt/easygo
systemctl daemon-reload
echo "EasyGo Panel uninstalled successfully"
EOF

chmod +x $INSTALL_DIR/uninstall.sh

echo ""
echo -e "${YELLOW}To uninstall EasyGo Panel, run: $INSTALL_DIR/uninstall.sh${NC}"