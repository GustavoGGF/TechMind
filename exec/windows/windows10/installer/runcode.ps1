# Script PowerShell para rodar os comandos necessários

# Executa o comando 'dotnet clean'
dotnet clean

# Executa o comando 'dotnet publish' com configuração Release
dotnet publish -c Release --self-contained -p:PublishSingleFile=true -o .

Clear-Host