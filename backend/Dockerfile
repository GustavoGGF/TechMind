# Use uma imagem base oficial do Python
FROM python:3.10.12

# Define o diretório de trabalho no contêiner
WORKDIR /techmind

# Copia o arquivo de requisitos para o contêiner
COPY requirements.txt ./

# Instala as dependências necessárias
RUN pip install --no-cache-dir -r requirements.txt

# Copia todo o conteúdo do projeto para o diretório de trabalho do contêiner
COPY . .

# Expõe a porta 3000
EXPOSE 3000

# Define o comando para rodar o servidor de desenvolvimento do Django na porta 3000
CMD ["python3", "manage.py", "runserver", "0.0.0.0:3000"]
