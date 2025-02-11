# Script PowerShell para rodar os comandos necessários

# Remove a pasta bin e seus conteúdos
Remove-Item -Path "C:\Users\adm.gfreitas\Desktop\TechMindInstallerW10\bin" -Recurse -Force

# Executa o comando 'dotnet clean'
dotnet clean

# Executa o comando 'dotnet publish' com configuração Release
dotnet publish -c Release
