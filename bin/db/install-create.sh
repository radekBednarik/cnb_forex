#!/bin/bash

# Update package list
sudo apt update

# Install PostgreSQL and its contrib package
sudo apt install postgresql postgresql-contrib -y

# Start PostgreSQL service
sudo service postgresql start

# Enable PostgreSQL service to start on boot
sudo systemctl enable postgresql

echo "PostgreSQL has been installed and started."
