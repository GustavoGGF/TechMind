using Microsoft.Win32;
using System.Diagnostics;

namespace TechMindInstallerW10;

public partial class Form1 : Form
{
    #region Constructor
    /// <summary>
    /// Construtor da classe Form1. Inicializa os componentes da interface do usuário e chama a função para exibir a confirmação da EULA.
    /// </summary>
    public Form1()
    {
        SoftwareExistenceCheck();
    }
    #endregion

    private void SoftwareExistenceCheck()
    {
        string keyPath = @"Software\Microsoft\Windows\CurrentVersion\Run";
        string valueName = "TechMind";

        // Acessando o registro
        using RegistryKey registryKey = Registry.CurrentUser.OpenSubKey(keyPath);
        
        // Verificando se a chave foi aberta com sucesso
        if (registryKey != null)
        {
            // Verificando se o valor existe
            if (registryKey.GetValue(valueName) is string value)
            {
                UninstallationConfirmation();
            }
            else
            {
                // Inicializa os componentes visuais do formulário.
                InitializeComponent();
                // Chama o método para lidar com a confirmação da EULA.
                EULAConfirmation();
            }
        }
        else
        {
            // Inicializa os componentes visuais do formulário.
            InitializeComponent();
            // Chama o método para lidar com a confirmação da EULA.
            EULAConfirmation();
        }
    }

    #region Label1_Click Event Handler
    /// <summary>
    /// Manipulador de evento para o clique no rótulo (label1).
    /// Alterna o estado do checkBox1 e habilita ou desabilita o botão button1 com base no estado atual do checkbox.
    /// </summary>
    /// <param name="sender">O objeto que disparou o evento.</param>
    /// <param name="e">Os argumentos do evento.</param>
    private void Label1_Click(object sender, EventArgs e)
    {
        if (!checkBox1.Checked)
        {
            // Marca o checkbox e habilita o botão se ele estiver desmarcado.
            checkBox1.Checked = true;
            this.button1.Enabled = true;
        }
        else
        {
            // Desmarca o checkbox e desabilita o botão se ele estiver marcado.
            checkBox1.Checked = false;
            this.button1.Enabled = false;
        }
    }
    #endregion


    private void Button1_Click(object sender, EventArgs e)
    {
        MessageBox.Show("EU Disse");
        this.button1.Enabled = checkBox1.Checked;
    }

    #region Event Handlers and Installation Step
    /// <summary>
        /// Evento acionado ao clicar no checkbox.
        /// Habilita ou desabilita o botão "Prosseguir" com base no estado do checkbox.
    /// </summary>
 

    /// <summary>
        /// Evento acionado ao clicar no label "Concordo".
        /// Alterna o estado do checkbox e, consequentemente, do botão "Prosseguir".
    /// </summary>


    /// <summary>
        /// Evento acionado ao clicar no botão "Prosseguir".
        /// Avança para o próximo passo da instalação.
    /// </summary>
    private void Next_Step(object sender, EventArgs e)
    {
        Installing_Step();
    }
    #endregion

    // <summary>
        // Executa o próximo passo da instalação.
        // Inicia a criação de pastas necessárias.
    // </summary>
    private void Installing_Step()
    {
        Create_Folder();
    }

    #region Make_Folder Method
    /// <summary>
    /// Método estático responsável por verificar e criar o diretório especificado, e em seguida invocar a operação assíncrona para obter arquivos.
    /// </summary>
    /// <param name="formInstance">Instância da classe Form1 usada para acessar métodos não estáticos.</param>
    private void Make_Folder(Form1 formInstance)
    {
        // Loader = new LoaderControl();
        string folderPath = @"C:\Program Files\techmind";

        try
        {
            // Verifica se a pasta já existe.
            if (!Directory.Exists(folderPath))
            {
                // Cria a nova pasta no caminho especificado.
                Directory.CreateDirectory(folderPath);
                this.loader.SetProgress(30); 
                _ = formInstance.Get_FilesAsync();  // Chama o método assíncrono para obter arquivos.
            }
            else
            {
                this.loader.SetProgress(30); 
                _ = formInstance.Get_FilesAsync();  // Invoca o método assíncrono mesmo que a pasta já exista.
            }
        }
        catch (Exception ex)
        {
            // Exibe uma mensagem de erro se ocorrer uma exceção durante a criação do diretório.
            MessageBox.Show($"Erro ao criar a pasta: {ex.Message}", "Erro", MessageBoxButtons.OK, MessageBoxIcon.Error);
        }
    }
    #endregion


    #region Get_FilesAsync Method
    /// <summary>
    /// Método assíncrono responsável por baixar arquivos de um servidor remoto e salvar localmente. 
    /// Atualiza a interface do usuário durante o processo de download.
    /// </summary>
    private async Task Get_FilesAsync()
    {
        // Atualiza o label para informar o status do download.
        this.label2.Text = "Baixando arquivos de SAPPP01...";
        this.loader.SetProgress(60); 
        this.panelEula.Controls.Remove(textBox1);  // Remove o TextBox1 da interface.

        string url = "http://sappp01:3000/donwload-files/";  // URL do servidor para download.
        string localPath = @"C:\Program Files\techmind\techmind.exe";  // Caminho local para salvar o arquivo.

        using HttpClient client = new();  // Inicializa um novo cliente HTTP.
        try
        {
            // Envia uma requisição GET para baixar o arquivo.
            HttpResponseMessage response = await client.GetAsync(url);
            response.EnsureSuccessStatusCode();  // Gera uma exceção se a resposta não for bem-sucedida.

            // Lê os bytes do conteúdo da resposta.
            byte[] fileBytes = await response.Content.ReadAsByteArrayAsync();

            // Salva o conteúdo do arquivo no caminho local especificado.
            await File.WriteAllBytesAsync(localPath, fileBytes);

            // Chama o método para criar entradas no Registro do Windows.
            CreateREGEdit();
        }
        catch (Exception ex)
        {
            // Exibe uma mensagem de erro se o download falhar.
            MessageBox.Show($"Erro ao fazer o download: {ex.Message}");
        }
    }
    #endregion

    #region CreateREGEdit Method
    /// <summary>
    /// Método responsável por criar uma entrada no Registro do Windows para executar o aplicativo no login do usuário.
    /// Atualiza o status da interface para indicar a criação do RegEdit.
    /// </summary>
    private void CreateREGEdit()
    {
        // Atualiza o label para informar que o RegEdit está sendo criado.
        this.label2.Text = "Criando RegEdit...";
        this.loader.SetProgress(95); 

        try
        {
            string registryKeyPath = @"Software\Microsoft\Windows\CurrentVersion\Run";  // Caminho da chave do Registro.
            string appName = "TechMind";  // Nome do valor a ser definido no Registro.
            string appPath = @"C:\Program Files\techmind\techmind.exe";  // Caminho do aplicativo.

            // Abre a chave do Registro no escopo de usuário atual com permissões de gravação.
            using RegistryKey registryKey = Registry.CurrentUser.OpenSubKey(registryKeyPath, true); // Permissão para escrever.
            if (registryKey != null)
            {
                // Define o valor para executar o aplicativo no login do usuário.
                registryKey.SetValue(appName, appPath);
                
                // Chama o método para exibir opções de reinicialização.
                OptionRestart();
            }
            else
            {
                // Exibe uma mensagem se a chave do Registro não for encontrada.
                MessageBox.Show("Erro ao acessar a chave do registro. A chave pode não existir.");
            }
        }
        catch (Exception ex)
        {
            // Exibe uma mensagem de erro se ocorrer uma exceção durante a operação.
            MessageBox.Show($"Erro: {ex.Message}");
        }
    }

    #endregion
    
    #region RestartNow Method
    /// <summary>
    /// Método responsável por reiniciar o sistema imediatamente quando o botão correspondente for clicado.
    /// Executa o comando 'shutdown' com as opções para reiniciar o sistema sem atraso.
    /// </summary>
    /// <param name="sender">O objeto que disparou o evento.</param>
    /// <param name="e">Os argumentos do evento.</param>
    private void RestartNow(object sender, EventArgs e)
    {
        try
        {
            // Inicia o processo de reinicialização do sistema com o comando 'shutdown'.
            System.Diagnostics.Process.Start("shutdown", "/r /t 0");
        }
        catch (Exception ex)
        {
            // Exibe uma mensagem de erro se a reinicialização falhar.
            MessageBox.Show($"Erro: {ex.Message}");
        }
    }
    #endregion

    #region RestartLatter Method
    /// <summary>
    /// Método responsável por agendar a reinicialização do sistema após um atraso de 15 minutos 
    /// quando o botão correspondente for clicado.
    /// Executa o comando 'shutdown' com a opção de reinicialização agendada.
    /// </summary>
    /// <param name="sender">O objeto que disparou o evento.</param>
    /// <param name="e">Os argumentos do evento.</param>
    private void RestartLatter(object sender, EventArgs e)
    {
        try
        {
            // Define o atraso para a reinicialização em 15 minutos (900 segundos).
            int delayInSeconds = 900;  // Tempo de atraso para reinicialização.
            
            // Inicia o processo de reinicialização agendada com o comando 'shutdown'.
            Process.Start("shutdown", $"/r /t {delayInSeconds}");
            this.Close();
        }
        catch (Exception ex)
        {
            // Exibe uma mensagem de erro se o agendamento falhar.
            MessageBox.Show($"Erro: {ex.Message}");
        }
    }
    #endregion

    
    #region Método RemoveRegEdit
    /// <summary>
    /// Este método é responsável por remover uma chave de registro no Windows.
    /// Ele atualiza o texto de um rótulo, altera o botão exibido e inicia um processo de carregamento.
    /// Após, tenta acessar o registro especificado e, se encontrado, exclui o valor associado.
    /// Caso contrário, finaliza o processo sem realizar alterações no registro.
    /// </summary>
    private void RemoveRegEdit(object sender, EventArgs e)
    {
        // Atualiza o rótulo indicando que o processo de remoção do registro está em andamento
        this.label3.Text = "Removendo Registro...";

        // Remove o botão 'button4' e adiciona o controle 'loader' para exibir o progresso
        this.Controls.Remove(button4);
        this.Controls.Add(this.loader);
        this.loader.SetProgress(10); 

        // Caminho da chave de registro a ser modificada
        string keyPath = @"Software\Microsoft\Windows\CurrentVersion\Run";
        // Nome do valor a ser removido
        string valueName = "TechMind";

        // Acessando o registro do usuário atual com permissão de escrita
        using RegistryKey registryKey = Registry.CurrentUser.OpenSubKey(keyPath, writable: true);
        
        // Verifica se a chave existe e o valor também
        if (registryKey != null && registryKey.GetValue(valueName) != null)
        {
            // Se o valor for encontrado, ele é excluído
            registryKey.DeleteValue(valueName);
            // Chama a função para parar qualquer processo relacionado
            StoppingProcess();
        }
        else
        {
            // Se o valor não for encontrado, chama a função para parar o processo sem alterações
            StoppingProcess();
        }
    }
    #endregion


    #region Método StoppingProcess
    /// <summary>
    /// Este método é responsável por finalizar o processo relacionado ao nome especificado (neste caso, "techmind").
    /// Ele atualiza o rótulo e o progresso do carregamento, tenta finalizar o processo se ele estiver em execução
    /// e, em seguida, chama o método para remover os arquivos.
    /// Caso o processo não seja encontrado, o método chama diretamente a função de remoção de arquivos.
    /// </summary>
    private void StoppingProcess()
    {
        // Atualiza o rótulo indicando que o processo de finalização está em andamento
        this.label3.Text = "Finalizando Processos...";
        this.loader.SetProgress(40); 

        // Nome do processo a ser finalizado, sem a extensão .exe
        string processName = "techmind";

        // Obtém todos os processos em execução com o nome especificado
        Process[] processes = Process.GetProcessesByName(processName);

        // Verifica se há algum processo em execução com o nome fornecido
        if (processes.Length > 0)
        {
            // Para cada processo encontrado, tenta finalizar
            foreach (Process process in processes)
            {
                try
                {
                    // Finaliza o processo
                    process.Kill();
                    // Chama o método para remover arquivos após o processo ser finalizado
                    RemovingFiles();
                }
                catch (Exception ex)
                {
                    // Caso haja erro ao finalizar o processo, o erro é registrado no console
                    Console.WriteLine($"Erro ao finalizar o processo: {ex.Message}");
                }
            }
        }
        else
        {
            // Caso nenhum processo seja encontrado, chama o método para remover arquivos
            RemovingFiles();
        }
    }
    #endregion


    #region Método RemovingFiles
    /// <summary>
    /// Este método é responsável por remover todos os arquivos e subpastas dentro de um diretório especificado.
    /// Após a remoção, ele exibe botões para o usuário reiniciar o sistema imediatamente ou mais tarde.
    /// Caso o diretório não exista, ele ainda exibe os botões de reinício.
    /// </summary>
    private void RemovingFiles()
    {
        // Atualiza o rótulo indicando que os arquivos estão sendo removidos
        this.label3.Text = "Removendo arquivos...";
        this.loader.SetProgress(70); 

        // Caminho do diretório a ser removido
        string folderPath = @"C:\Program Files\techmind";

        try
        {
            // Verifica se o diretório existe
            if (Directory.Exists(folderPath))
            {
                // Deleta todos os arquivos dentro da pasta
                foreach (string file in Directory.GetFiles(folderPath))
                {
                    File.Delete(file);
                    // Após a remoção dos arquivos, atualiza o rótulo e exibe o botão para reiniciar
                    this.label3.Text = "Desinstalação Concluida!";
                    this.button2 = new System.Windows.Forms.Button
                    {
                        Location = new System.Drawing.Point(200, 250),
                        Width = 200,
                        Text = "Reiniciar Agora"
                    };
                    this.button2.Click += new EventHandler(RestartNow);
                    this.Controls.Remove(this.loader);
                    this.Controls.Add(button2);
                }

                // Deleta todas as subpastas dentro da pasta
                foreach (string subDir in Directory.GetDirectories(folderPath))
                {
                    Directory.Delete(subDir, true); // 'true' para excluir recursivamente
                    this.label3.Text = "Desinstalação Concluida!";
                    // Exibe os botões para reiniciar
                    this.button2 = new System.Windows.Forms.Button
                    {
                        Location = new System.Drawing.Point(200, 250),
                        Width = 200,
                        Text = "Reiniciar Agora"
                    };
                    this.button2.Click += new EventHandler(RestartNow);
                    this.button3 = new System.Windows.Forms.Button
                    {
                        Location = new System.Drawing.Point(520, 250),
                        Width = 200,
                        Text = "Reiniciar Depois"
                    };
                    this.button3.Click += new EventHandler(RestartLatter);
                    this.Controls.Add(button2);
                    this.Controls.Add(button3);
                    this.Controls.Remove(this.loader);
                }

                // Deleta a pasta principal
                Directory.Delete(folderPath);
                this.label3.Text = "Desinstalação Concluida!";
                this.button2 = new System.Windows.Forms.Button
                {
                    Location = new System.Drawing.Point(200, 250),
                    Width = 200,
                    Text = "Reiniciar Agora"
                };
                this.button2.Click += new EventHandler(RestartNow);
                this.button3 = new System.Windows.Forms.Button
                {
                    Location = new System.Drawing.Point(520, 250),
                    Width = 200,
                    Text = "Reiniciar Depois"
                };
                this.button3.Click += new EventHandler(RestartLatter);
                this.Controls.Add(button2);
                this.Controls.Add(button3);
                this.Controls.Remove(this.loader);
            }
            else
            {
                // Caso o diretório não exista, ainda exibe os botões de reinício
                this.label3.Text = "Desinstalação Concluida!";
                this.button2 = new System.Windows.Forms.Button
                {
                    Location = new System.Drawing.Point(200, 250),
                    Width = 200,
                    Text = "Reiniciar Agora"
                };
                this.button2.Click += new EventHandler(RestartNow);
                this.button3 = new System.Windows.Forms.Button
                {
                    Location = new System.Drawing.Point(520, 250),
                    Width = 200,
                    Text = "Reiniciar Depois"
                };
                this.button3.Click += new EventHandler(RestartLatter);
                this.Controls.Add(button2);
                this.Controls.Add(button3);
                this.Controls.Remove(this.loader);
            }
        }
        catch (Exception ex)
        {
            // Exibe uma mensagem de erro caso ocorra alguma exceção ao excluir a pasta
            MessageBox.Show($"Erro ao excluir a pasta: {ex.Message}");
        }
    }
    #endregion

}
