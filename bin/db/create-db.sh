#!/bin/bash

# Prompt user for database name
read -p "Enter name for the PostgreSQL database: " dbname

# Create PostgreSQL database
sudo -u postgres psql -c "CREATE DATABASE $dbname;"

echo "Database '$dbname' has been created."
