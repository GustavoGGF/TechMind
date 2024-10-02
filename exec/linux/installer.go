package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// Função para transferir um arquivo do servidor remoto para o local
func transferFile(username, password, remoteFile, localFile string) error {
    // Configurar a autenticação SSH
    config := &ssh.ClientConfig{
        User: username,
        Auth: []ssh.AuthMethod{
            ssh.Password(password),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Para simplificar, não verifique a chave do host
    }

    // Estabelecer conexão SSH
    conn, err := ssh.Dial("tcp", "snas01:22", config)
    if err != nil {
        return fmt.Errorf("erro ao conectar: %w", err)
    }
    defer conn.Close()

    // Executar o comando SCP para transferir do remoto para o local
    cmd := exec.Command("scp", "-o", "StrictHostKeyChecking=no", fmt.Sprintf("%s@snas01:%s", username, remoteFile), localFile)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Executar o comando
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("erro ao executar SCP: %w", err)
    }

    return nil
}

func main(){
	reader := bufio.NewReader(os.Stdin)

    // Solicitar o nome de usuário
    fmt.Print("Digite o nome de usuário: ")
    username, _ := reader.ReadString('\n')
	// Remover nova linha do nome de usuário
	username = username[:len(username)-1]

    // Solicitar a senha
    fmt.Print("Digite a senha: ")
    bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
    fmt.Println() // Para nova linha após a senha

    // Converter a senha de bytes para string
    password := string(bytePassword)

    // Definir o arquivo remoto e local
    remoteFile := "/shares/node00/node5/lun0/f-sti01/002-Programas/DEP_Suporte/techmind/techmind" // Substitua pelo caminho do arquivo remoto
    localFile := "/usr/local/bin"

	    // Chamar a função de transferência de arquivo
		if err := transferFile(username, password, remoteFile, localFile); err != nil {
			fmt.Println(err)
			return
		}

    // Torna o programa executável
	erro := os.Chmod(localFile + "/techmind", 0755)
	if erro != nil {
		fmt.Printf("Erro ao tornar o programa executável: %s\n", erro)
		return
	}
	
		fmt.Println("Arquivo transferido com sucesso!")
}