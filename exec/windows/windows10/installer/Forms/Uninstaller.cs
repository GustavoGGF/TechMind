using System.Diagnostics;
using Microsoft.Win32;

namespace TechMindInstallerW10;

#region 
/// <summary>
/// FUnção parcial de Main, esse arquivo irá manter a Desinstalação do TechMind
/// </summary>
partial class Main
{
    private LoaderControl loader;
    private System.Windows.Forms.Label label3;
    private System.Windows.Forms.Button button4;

    #region Func UninstallationConfirmation
    /// <summary>
    /// Tela de Confirmação para Realizar a Desinstalação do TechMind
    /// </summary>
    private void UninstallationConfirmation()
    {
        try
        {
            // Cria uma nova instância de um contêiner de componentes.
            this.components = new System.ComponentModel.Container();

            // Define o modo de escalonamento automático para o formulário.
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;

            // Obtém o tamanho da tela primária
            var screenSize = Screen.PrimaryScreen.WorkingArea;

            // Calcula 15% da largura e 10% da altura
            int largura = (int)(screenSize.Width * 0.35);
            int altura = (int)(screenSize.Height * 0.25);

            // Define o tamanho do formulário
            this.ClientSize = new Size(largura, altura);

            // Desativa a capacidade de maximizar a janela do formulário.
            this.MaximizeBox = false;

            // Define o estilo da borda do formulário como fixa (não pode ser redimensionada).
            this.FormBorderStyle = FormBorderStyle.FixedSingle;

            // Define o título do formulário.
            this.Text = "Desinstalação TechMind";

            int posXLoader = (int)(this.ClientSize.Width * 0.1); // 10% da largura
            int posYLoader = (int)(this.ClientSize.Height * 0.3); // 30% da altura
            int widthLoader = (int)(this.ClientSize.Width * 0.8); // 80% da largura
            int heightLoader = (int)(this.ClientSize.Height * 0.1); // 10% da altura
            loader = new LoaderControl
            {
                Location = new Point(posXLoader, posYLoader),
                Size = new Size(widthLoader, heightLoader)
            };

            // Cria e configura um label para a opção de concordância.
            this.label3 = new System.Windows.Forms.Label();
            int posXLabel3 = (int)(this.ClientSize.Width * 0.1); // 1% da largura
            int posYLabel3 = (int)(this.ClientSize.Height * 0.2); // 1% da altura
            this.label3.Location = new System.Drawing.Point(posXLabel3, posYLabel3);
            int widthLabel3 = (int)(this.ClientSize.Width * 0.9); // 90% da largura
            int heightLabel3 = (int)(this.ClientSize.Height * 0.2); // 20% da altura
            this.label3.Width = widthLabel3;
            this.label3.Height = heightLabel3;
            this.label3.Font = new Font(label3.Font, FontStyle.Bold);
            this.label3.Text = "Você esta prestes a desinstalar o TechMind, caso esteja ciente disso clique no botão para dar andamento.";

            this.button4 = new System.Windows.Forms.Button();
            int posX = (int)(this.ClientSize.Width * 0.40); // 40% da largura
            int posY = (int)(this.ClientSize.Height * 0.40); // 30% da altura
            this.button4.Location = new Point(posX, posY);
            int widthButton4 = (int)(this.ClientSize.Width * 0.2); // 20% da largura
            int heightButton4 = (int)(this.ClientSize.Height * 0.15); // 15% da altura
            this.button4.Height = heightButton4;
            this.button4.Width = widthButton4;
            this.button4.Text = "Prosseguir";
            this.button4.Click += new System.EventHandler(this.RemoveRegEdit);

            this.Controls.Add(label3);
            this.Controls.Add(button4);
            this.ResumeLayout(false);
            this.PerformLayout();

            this.Icon = Icon.ExtractAssociatedIcon(Application.ExecutablePath);
        }
        catch (Exception ex)
        {
            MessageBox.Show("Erro ao carregar o ícone: " + ex.Message);
        }
    }
    #endregion

    #region Func RemoveRegEdit
    /// <summary>
    /// Remove o Arquivo de Registro referente ao TechMind
    /// </summary>
    /// <param name="sender"></param>
    /// <param name="e"></param>
    private void RemoveRegEdit(object sender, EventArgs e)
    {
        // Atualiza o rótulo indicando que o processo de remoção do registro está em andamento
        this.label3.Text = "Removendo Registro...";

        // Remove o botão 'button4' e adiciona o controle 'loader' para exibir o progresso
        this.Controls.Remove(button4);
        this.Controls.Add(this.loader);
        this.loader.SetProgress(20);

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
            RemoveFirewallRule();
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

    #region Func RemoveFirewallRule
    /// <summary>
    /// Remove a Regra de FireWall do Windows refente ao TechMind
    /// </summary>
    private void RemoveFirewallRule()
    {
        this.label3.Text = "Removendo Regra do FireWall...";
        this.loader.SetProgress(35);

        try
        {
            // Nome da regra que foi adicionada
            string ruleName = "TechMind";

            ProcessStartInfo psi = new()
            {
                FileName = "netsh",
                Arguments = $"advfirewall firewall delete rule name=\"{ruleName}\"",
                UseShellExecute = false,
                CreateNoWindow = true,
                RedirectStandardOutput = true,
                RedirectStandardError = true
            };

            using Process process = Process.Start(psi)!;
            string output = process.StandardOutput.ReadToEnd();
            string error = process.StandardError.ReadToEnd();

            process.WaitForExit();

            if (process.ExitCode == 0)
            {
                this.label3.Text = "✅ Regra de firewall removida com sucesso.";
            }
            else
            {
                MessageBox.Show($"⚠️ Falha ao remover regra. Código: {process.ExitCode}");
                MessageBox.Show($"Erro: {error}");
            }
        }
        catch (Exception ex)
        {
            MessageBox.Show($"❌ Erro inesperado ao remover regra de firewall: {ex.Message}");
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
        this.loader.SetProgress(45);

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
            int widthButton2 = (int)(this.ClientSize.Width * 0.2); // 20% da largura
            int posXButton2 = (int)(this.ClientSize.Width * 0.2); // 20% da largura
            int posXButton2Latter = (int)(this.ClientSize.Width * 0.6); // 60% da largura
            int posYButton2 = (int)(this.ClientSize.Height * 0.5); // 50% da altura
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
                        Location = new System.Drawing.Point(posXButton2, posYButton2),
                        Width = widthButton2,
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
                        Location = new System.Drawing.Point(posXButton2, posYButton2),
                        Width = widthButton2,
                        Text = "Reiniciar Agora"
                    };
                    this.button2.Click += new EventHandler(RestartNow);
                    this.button3 = new System.Windows.Forms.Button
                    {
                        Location = new System.Drawing.Point(posXButton2Latter, posYButton2),
                        Width = widthButton2,
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
                    Location = new System.Drawing.Point(posXButton2, posYButton2),
                    Width = widthButton2,
                    Text = "Reiniciar Agora"
                };
                this.button2.Click += new EventHandler(RestartNow);
                this.button3 = new System.Windows.Forms.Button
                {
                    Location = new System.Drawing.Point(posXButton2Latter, posYButton2),
                    Width = widthButton2,
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
                    Location = new System.Drawing.Point(posXButton2, posYButton2),
                    Width = widthButton2,
                    Text = "Reiniciar Agora"
                };
                this.button2.Click += new EventHandler(RestartNow);
                this.button3 = new System.Windows.Forms.Button
                {
                    Location = new System.Drawing.Point(posXButton2Latter, posYButton2),
                    Width = widthButton2,
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
#endregion