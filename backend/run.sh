#!/bin/bash

# Script para correr el servidor de desarrollo

export PATH=$PATH:$HOME/go/bin
export GOTOOLCHAIN=local

# Variables de entorno para PostgreSQL
export DB_HOST=/var/run/postgresql
export DB_USER=devuser
export DB_NAME=buylist_db
export DB_DRIVER=postgres
export PORT=8080
export ENV=development

echo "🚀 Starting BuyList Manager API..."
echo "📍 Database: PostgreSQL (buylist_db)"
echo "🌐 Port: 8080"
echo ""

# Compilar y correr
cd cmd/api && go run main.go
