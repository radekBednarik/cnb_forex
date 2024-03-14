#!/bin/bash

# Prompt user for username and password
read -p "Enter username for PostgreSQL: " username
read -s -p "Enter password for PostgreSQL user $username: " password
echo

# Create PostgreSQL user with administrative privileges
sudo -u postgres psql -c "CREATE ROLE $username WITH SUPERUSER CREATEDB CREATEROLE LOGIN ENCRYPTED PASSWORD '$password';"

echo "User '$username' has been created with administrative privileges."
