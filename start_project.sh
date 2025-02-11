#!/bin/bash

# Sai imediatamente se ocorrer um erro
set -e

# Navega até a pasta backend e remove a pasta static
cd backend/
rm -rf static/
cd ..

# Navega até a pasta frontend e constrói o projeto Angular
cd frontend/
ng build
cd ..

# Volta para o backend e inicia o servidor Django
cd backend/
python3 manage.py runserver
