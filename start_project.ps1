
clear
# Navega até a pasta frontend e constrói o projeto Angular
cd frontend/
ng build
cd ..

# Volta para o backend e inicia o servidor Django
python.exe .\backend\manage.py runserver
