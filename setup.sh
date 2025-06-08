#!/bin/bash

# Update system
echo "Updating system..."
sudo apt update && sudo apt upgrade -y

# Install required packages
echo "Installing required packages..."
sudo apt install -y git golang-go mysql-server

# Create project directory
echo "Creating project directory..."
sudo mkdir -p /var/www/html/simpadu
sudo chown -R $USER:$USER /var/www/html/simpadu

# Clone repository (if using git)
# echo "Cloning repository..."
# git clone <your-repo-url> /var/www/html/simpadu

# Copy project files
echo "Copying project files..."
cp -r ./* /var/www/html/simpadu/

# Setup MySQL
echo "Setting up MySQL..."
sudo mysql -e "CREATE DATABASE IF NOT EXISTS simpadu;"
sudo mysql -e "CREATE USER IF NOT EXISTS 'simpadu_user'@'localhost' IDENTIFIED BY 'password_yang_aman';"
sudo mysql -e "GRANT ALL PRIVILEGES ON simpadu.* TO 'simpadu_user'@'localhost';"
sudo mysql -e "FLUSH PRIVILEGES;"

# Import database
echo "Importing database..."
mysql -u simpadu_user -p'simpadu_user' simpadu < test_data.sql

# Create .env file
echo "Creating .env file..."
cat > /var/www/html/simpadu/.env << EOL
DB_USER=simpadu_user
DB_PASSWORD=password_yang_aman
DB_HOST=localhost
DB_PORT=3306
DB_NAME=simpadu
simpadu_jwt_key=Kj8#mP9\$vL2@nX5&hQ7*wR4!tY6^cF3
GIN_MODE=release
EOL

# Build Go application
echo "Building application..."
cd /var/www/html/simpadu
go mod download
go build -o simpadu

# Create systemd service
echo "Creating systemd service..."
sudo cat > /etc/systemd/system/simpadu.service << EOL
[Unit]
Description=Simpadu Application
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/html/simpadu
ExecStart=/var/www/html/simpadu/simpadu
Restart=always
Environment=DB_USER=simpadu_user
Environment=DB_PASSWORD=password_yang_aman
Environment=DB_HOST=localhost
Environment=DB_PORT=3306
Environment=DB_NAME=simpadu
Environment=simpadu_jwt_key=Kj8#mP9\$vL2@nX5&hQ7*wR4!tY6^cF3
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOL

# Set permissions
echo "Setting permissions..."
sudo chown -R www-data:www-data /var/www/html/simpadu
sudo chmod -R 755 /var/www/html/simpadu

# Create Nginx configuration
echo "Creating Nginx configuration..."
sudo cat > /etc/nginx/sites-available/simpadu << EOL
server {
    listen 80;
    server_name ti054c01.agussbn.my.id;
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl;
    server_name ti054c01.agussbn.my.id;

    ssl_certificate /etc/letsencrypt/live/ti054c01.agussbn.my.id/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/ti054c01.agussbn.my.id/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
    }
}
EOL

# Enable Nginx site
echo "Enabling Nginx site..."
sudo ln -sf /etc/nginx/sites-available/simpadu /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# Start and enable service
echo "Starting service..."
sudo systemctl daemon-reload
sudo systemctl enable simpadu
sudo systemctl start simpadu

echo "Setup completed!"
echo "Please check the service status with: sudo systemctl status simpadu"
echo "Check the logs with: sudo journalctl -u simpadu -f" 
