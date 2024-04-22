# Use a imagem oficial do Golang como base
FROM golang:1.22.0

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /TechMind

# Copie os arquivos do projeto para o contêiner
COPY . /TechMind

# Baixe e instale as dependências
RUN go mod download

# Compile o código
RUN go build

EXPOSE 3000

# Especifica o comando padrão a ser executado ao iniciar o contêiner
CMD ["./techmind"]
