clear

# Verifica se a pasta 'backend\static' existe e remove se existir
if (Test-Path "backend\static") {
    Remove-Item -Path "backend\static" -Recurse -Force
}

# Navega até a pasta frontend e constrói o projeto Angular
cd frontend/
ng build
cd ..

# Volta para o backend e inicia o servidor Django
cd backend/
python.exe -m daphne -p 8000 techmind.asgi:application
