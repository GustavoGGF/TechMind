using Microsoft.Win32;
using System.Diagnostics;
using System.Text.Json;

namespace TechMindInstallerW10;

public class VersionInfo
{
    public string latest_version { get; set; }
}

public class VersionCurrent
{
    public string current_version { get; set; }
}


partial class Main
{
    private System.Windows.Forms.Panel panelEula;
    private System.Windows.Forms.RichTextBox rtextb1;
    private System.Windows.Forms.Label label1;
    private System.Windows.Forms.CheckBox checkBox1;
    private System.Windows.Forms.Button button1;
    private System.Windows.Forms.Button button2;
    private System.Windows.Forms.Button button3;
    private System.Windows.Forms.Label label2;

    #region Func EULAConfirmation
    /// <summary>
    /// Configura a tela de confirmação do EULA (Contrato de Licença de Usuário Final).
    /// Adiciona os componentes necessários ao painel, incluindo texto, 
    /// checkbox, botão e label.
    /// </summary>
    private void EULAConfirmation()
    {
        //Cria um painel para exibir o conteúdo do EULA.
        this.panelEula = new System.Windows.Forms.Panel();
        this.SuspendLayout();

        // Configura o painel para ocupar todo o espaço disponível no formulário.
        this.panelEula.Dock = DockStyle.Fill;
        this.panelEula.BackColor = Color.FromArgb(254, 250, 224);
        this.panelEula.BorderStyle = BorderStyle.None;

        // Cria e configura uma RichTextBox para exibir o texto do EULA.
        this.rtextb1 = new System.Windows.Forms.RichTextBox();
        int widthRTextb1 = (int)(this.ClientSize.Width * 0.65); // 65% da largura
        int heighRTextb1 = (int)(this.ClientSize.Height * 0.8); // 80% da altura
        this.rtextb1.Size = new Size(widthRTextb1, heighRTextb1);
        int posXTextb1 = (int)(this.ClientSize.Width * 0.2); // 20% da largura
        int posYTextb1 = (int)(this.ClientSize.Height * 0); // 0% da altura
        this.rtextb1.Location = new System.Drawing.Point(posXTextb1, posYTextb1);
        this.rtextb1.ReadOnly = true;
        this.rtextb1.Rtf = @"{\rtf1\ansi\ansicpg1252\uc1 
                \pard\cf0\b CONTRATO DE LICENÇA DE USUÁRIO FINAL (EULA)\b0\par
                \pard\cf0 TechMind - Ferramenta de Inventário\par
                \pard\cf0 Desenvolvido por Gustavo Guilherme de Freitas\par
                \pard\cf0\b 1. INTRODUÇÃO\b0\par
                \pard\cf0 Este Contrato de Licença de Usuário Final (EULA) rege o uso do software TechMind, desenvolvido exclusivamente para a empresa Lupatech. Ao instalar, acessar ou usar o TechMind, você concorda com os termos e condições aqui descritos. Caso não concorde, você não está autorizado a utilizar o software.\par
                \pard\cf0\b 2. TERMOS DE USO\b0\par
                \pard\cf0\b 2.1. Acesso restrito:\b0\par
                \pard\cf0 O TechMind só pode ser executado com permissões administrativas. Usuários com permissões básicas não têm autorização para interagir diretamente com o software.\par
                \pard\cf0\b 2.2. Uso corporativo exclusivo:\b0\par
                \pard\cf0 Este software é licenciado exclusivamente para a empresa Lupatech, e seu uso é estritamente limitado às atividades corporativas internas da mesma.\par
                \pard\cf0\b 2.3. Propriedade intelectual:\b0\par
                \pard\cf0 O TechMind é propriedade exclusiva de Gustavo Guilherme de Freitas, CPF final 853.680, que detém todos os direitos autorais e intelectuais sobre o software.\par
                \pard\cf0\b 3. RESTRIÇÕES\b0\par
                \pard\cf0\b 3.1. Proibições:\b0\par
                \pard\cf0 É estritamente proibido realizar engenharia reversa, distribuir, modificar ou redistribuir o software sem autorização expressa e por escrito do proprietário.\par
                \pard\cf0\b 3.2. Uso não autorizado:\b0\par
                \pard\cf0 Qualquer uso fora do escopo corporativo definido neste contrato, ou por entidades que não pertençam à Lupatech, será considerado uma violação grave deste EULA.\par
                \pard\cf0\b 4. LIMITAÇÕES E RESPONSABILIDADES\b0\par
                \pard\cf0\b 4.1. Limitação de responsabilidade:\b0\par
                \pard\cf0 O desenvolvedor e a equipe de TI da Lupatech não se responsabilizam por quaisquer danos, diretos ou indiretos, causados pelo uso do software ou decorrentes de falhas ou interrupções de sua operação.\par
                \pard\cf0\b 4.2. Revogação de licença:\b0\par
                \pard\cf0 Em caso de violação de dados sensíveis dos usuários por meio do TechMind, a licença perpétua será revogada imediatamente, conforme determinação unilateral do desenvolvedor.\par
                \pard\cf0\b 5. LICENÇA E DISTRIBUIÇÃO\b0\par
                \pard\cf0\b 5.1. Licença:\b0\par
                \pard\cf0 A Lupatech possui uma licença perpétua para a versão atual do TechMind, desde que os termos deste contrato sejam respeitados.\par
                \pard\cf0\b 5.2. Limite de dispositivos:\b0\par
                \pard\cf0 Não há limite para o número de dispositivos conectados ao software, respeitando a capacidade do servidor onde os dados são armazenados.\par
                \pard\cf0\b 6. FUNCIONALIDADE DO SOFTWARE\b0\par
                \pard\cf0\b 6.1. Coleta de dados:\b0\par
                \pard\cf0 O TechMind coleta exclusivamente dados relacionados ao hardware dos dispositivos (como processador, memória RAM, placa-mãe, placa de vídeo, etc.), nomes de softwares instalados, suas versões e informações sobre a licença do Windows.\par
                \pard\cf0\b 6.2. Privacidade dos dados:\b0\par
                \pard\cf0 O software não coleta dados sensíveis dos usuários. O objetivo é unicamente fornecer controle interno para a equipe de TI da Lupatech.\par
                \pard\cf0\b 7. CONSIDERAÇÕES FINAIS\b0\par
                \pard\cf0\b 7.1. Este contrato será regido pelas leis brasileiras.\b0\par
                \pard\cf0\b 7.2. Quaisquer disputas relacionadas ao uso do TechMind deverão ser resolvidas em território nacional, preferencialmente por meio de arbitragem.\b0\par
                \pard\cf0\b 7.3. Ao utilizar este software, você declara ter lido, compreendido e aceitado os termos aqui descritos.\b0\par
                \pard\cf0\b Gustavo Guilherme de Freitas\b0\par
                \pard\cf0 Proprietário e Desenvolvedor do TechMind\par
                \pard\cf0 Versão Atual: 2.0.0\par
            }";

        // Cria e configura um label para a opção de concordância.
        this.label1 = new System.Windows.Forms.Label();
        int posXLabel1 = (int)(this.ClientSize.Width * 0.22); // 22% da largura
        int posYLabel1 = (int)(this.ClientSize.Height * 0.85); // 85% da altura
        this.label1.Location = new System.Drawing.Point(posXLabel1, posYLabel1);
        this.label1.Text = "Concordo.";
        this.label1.Click += new System.EventHandler(this.Label1_Click);

        // Cria e configura um checkbox para o usuário concordar com o EULA.
        this.checkBox1 = new System.Windows.Forms.CheckBox();
        int posXCheckBox1 = (int)(this.ClientSize.Width * 0.20); // 20% da largura
        int posYCheckBox1 = (int)(this.ClientSize.Height * 0.85); // 85% da altura
        this.checkBox1.Location = new System.Drawing.Point(posXCheckBox1, posYCheckBox1);
        int widthCheckBox1 = (int)(this.ClientSize.Width * 0.029); // 2.9% da largura
        int heighCheckBox1 = (int)(this.ClientSize.Height * 0.029); // 2.9% da altura
        this.checkBox1.Size = new Size(widthCheckBox1, heighCheckBox1);
        this.checkBox1.Click += new System.EventHandler(this.Button1_Click);

        // Cria e configura um botão para avançar, que inicialmente está desabilitado.
        this.button1 = new System.Windows.Forms.Button();
        int posXButton1 = (int)(this.ClientSize.Width * 0.75); // 75% da largura
        int posYButton1 = (int)(this.ClientSize.Height * 0.85); // 85% da altura
        this.button1.Location = new System.Drawing.Point(posXButton1, posYButton1);
        this.button1.Text = "Prosseguir";
        this.button1.Enabled = false;
        this.button1.Click += new System.EventHandler(this.Next_Step);

        // Adiciona os componentes ao painel do EULA.
        this.panelEula.Controls.Add(rtextb1);
        this.panelEula.Controls.Add(label1);
        this.panelEula.Controls.Add(checkBox1);
        this.panelEula.Controls.Add(button1);

        // Adiciona o painel ao formulário principal.
        this.Controls.Add(panelEula);

        // Finaliza a configuração do layout.
        this.ResumeLayout(false);
        this.PerformLayout();

        this.Icon = Icon.ExtractAssociatedIcon(Application.ExecutablePath);
    }
    #endregion

    #region Func Label1_Click
    /// <summary>
    /// Manipulador de evento para o clique no texto (label1).
    /// Alterna o estado do checkBox1 e habilita ou desabilita o 
    /// botão button1 com base no estado atual do checkbox.
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

    #region Func Button1_Click
    /// <summary>
    /// Manipulador de evento para o clique no botão (button1).
    /// Alterna o estado do checkBox1 e habilita ou desabilita o 
    /// botão button1 com base no estado atual do checkbox.
    /// </summary>
    /// <param name="sender">O objeto que disparou o evento.</param>
    /// <param name="e">Os argumentos do evento.</param>
    private void Button1_Click(object sender, EventArgs e)
    {
        if (button1.Enabled)
        {
            checkBox1.Checked = false;
            this.button1.Enabled = false;
        }
        else
        {
            checkBox1.Checked = true;
            this.button1.Enabled = true;
        }
    }
    #endregion

    #region Func Next_Step
    /// <summary>
    /// Função que incia o processo de instalação do TechMind
    /// Comeã Criando o diretório
    /// </summary>
    /// <param name="sender"></param>
    /// <param name="e"></param> ativado ao clicar no botão de Prosseguir
    private void Next_Step(object sender, EventArgs e)
    {
        Initialize_Component(); //Criando diretório para o programa do TechMind
    }
    #endregion

    #region Func Initialize_Component
    /// <summary>
    /// Inicializa a tela com o loader e informação de status
    /// </summary>
    private void Initialize_Component()
    {
        int posXLoader = (int)(this.ClientSize.Width * 0.058); // 0.58% da largura
        int posYLoader = (int)(this.ClientSize.Height * 0.3); // 3% da altura
        int widthLoader = (int)(this.ClientSize.Width * 0.9); // 90% da largura
        int heighLoader = (int)(this.ClientSize.Height * 0.1); // 10% da altura
        loader = new LoaderControl
        {
            Location = new Point(posXLoader, posYLoader),
            Size = new Size(widthLoader, heighLoader)
        };

        // Suspende temporariamente a disposição dos controles para melhor desempenho durante a modificação.
        this.SuspendLayout();

        // Remove controles existentes do painel EULA.
        this.panelEula.Controls.Remove(rtextb1);
        this.panelEula.Controls.Remove(checkBox1);
        this.panelEula.Controls.Remove(button1);
        this.panelEula.Controls.Remove(label1);

        this.label2 = new System.Windows.Forms.Label();
        int posXLabel2 = (int)(this.ClientSize.Width * 0.3); // 30% da largura
        int posYLabel2 = (int)(this.ClientSize.Height * 0.2); // 20% da altura
        this.label2.Location = new System.Drawing.Point(posXLabel2, posYLabel2);
        int widthLabel2 = (int)(this.ClientSize.Width * 0.3); // 30% da largura
        int heighLabel2 = (int)(this.ClientSize.Height * 0.1); // 10% da altura
        this.label2.Size = new Size(widthLabel2, heighLabel2);

        // Adiciona os novos controles ao painel.
        this.panelEula.Controls.Add(label2);
        this.panelEula.Controls.Add(this.loader);

        // Retoma a disposição normal dos controles e atualiza o layout.
        this.ResumeLayout(false);
        this.PerformLayout();
        loader.SetProgress(10);

        // Chama o método responsável por criar a pasta.
        _ = Create_FolderAsync(this);
    }
    #endregion

    #region Func Create_Folder
    /// <summary>
    /// Cria o Diretorio onde ficará o TechMind
    /// </summary>
    /// <param name="formInstance"></param>
    private async Task Create_FolderAsync(Main formInstance)
    {
        this.label2.Text = "Criando Diretório...";
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
                var version = await GetVersionAsync();
                await Get_FilesAsync(version);  // Chama o método assíncrono para obter arquivos.
            }
            else
            {
                this.loader.SetProgress(30);
                var version = await GetVersionAsync();
                await Get_FilesAsync(version);  // Invoca o método assíncrono mesmo que a pasta já exista.
            }
        }
        catch (Exception ex)
        {
            // Exibe uma mensagem de erro se ocorrer uma exceção durante a criação do diretório.
            MessageBox.Show($"Erro ao criar a pasta: {ex.Message}", "Erro", MessageBoxButtons.OK, MessageBoxIcon.Error);
        }
    }
    #endregion

    #region Func Get_FilesAsync
    /// <summary>
    /// Baixa o TechMind do servidor
    /// </summary>
    /// <returns></returns>
    private async Task Get_FilesAsync(string version)
    {
        // Atualiza o label para informar o status do download.
        this.label2.Text = "Baixando arquivos de SAPPP01...";
        this.loader.SetProgress(60);

        string url = "https://techmind.lupatech.com.br/download-files/techmind/" + version;  // URL do servidor para download.
        string localPath = $@"C:\Program Files\techmind\techmind.exe";  // Caminho local para salvar o arquivo.


        var handler = new HttpClientHandler
        {
            // Ignora validação do certificado
            ServerCertificateCustomValidationCallback = (message, cert, chain, errors) => true
        };

        using HttpClient client = new(handler);  // Inicializa um novo cliente HTTP.
        client.DefaultRequestHeaders.Add("X-Requested-With", "XMLHttpRequest");  // Cabeçalho necessário
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
            CreateREGEdit(version);
        }
        catch (Exception ex)
        {
            // Exibe uma mensagem de erro se o download falhar.
            MessageBox.Show($"Erro ao fazer o download: {ex.Message}");
        }
    }
    #endregion

    #region CreateREGEdit
    /// <summary>
    /// Método responsável por criar uma entrada no Registro do Windows para 
    /// executar o aplicativo no login do usuário.
    /// </summary>
    private void CreateREGEdit(string version)
    {
        // Atualiza o label para informar que o RegEdit está sendo criado.
        this.label2.Text = "Criando RegEdit...";
        this.loader.SetProgress(75);

        try
        {
            string registryKeyPath = @"Software\Microsoft\Windows\CurrentVersion\Run";  // Caminho da chave do Registro.
            string appName = "TechMind";  // Nome do valor a ser definido no Registro.
            string appPath = $@"C:\Program Files\techmind\techmind.exe";  // Caminho do aplicativo.

            // Abre a chave do Registro no escopo de usuário atual com permissões de gravação.
            using RegistryKey registryKey = Registry.CurrentUser.OpenSubKey(registryKeyPath, true); // Permissão para escrever.
            if (registryKey != null)
            {
                // Define o valor para executar o aplicativo no login do usuário.
                registryKey.SetValue(appName, appPath);

                AddFirewallRule();
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

    private async Task<string> GetVersionAsync()
    {
        this.label2.Text = "Verificando Versão do programa...";
        this.loader.SetProgress(35);
        var handler = new HttpClientHandler
        {
            // Ignora validação do certificado
            ServerCertificateCustomValidationCallback = (message, cert, chain, errors) => true
        };
        string url = "https://techmind.lupatech.com.br/get-current-version/windows10";
        using HttpClient client = new(handler);

        try
        {
            HttpResponseMessage response = await client.GetAsync(url);

            response.EnsureSuccessStatusCode(); // Lança exceção se status != 200-299

            string jsonResponse = await response.Content.ReadAsStringAsync();

            // Se quiser desserializar:
            var resultado = JsonSerializer.Deserialize<VersionInfo>(jsonResponse);

            string version = resultado.latest_version;

            SaveVersion(version);
            return version;
        }
        catch (Exception ex)
        {
            MessageBox.Show($"Erro ao fazer a requisição: {ex.Message}");
            this.Close();
            return null;
        }
    }

    private void SaveVersion(string versao)
    {
        this.label2.Text = "Criando arquivo de Configuração...";
        this.loader.SetProgress(45);

        var versaoInfo = new VersionCurrent
        {
            current_version = versao
        };

        string json = JsonSerializer.Serialize(versaoInfo, new JsonSerializerOptions { WriteIndented = true });

        // Caminho: C:\Program Files\techmind\config\version.json
        string pastaBase = @"C:\Program Files\techmind\configs";
        string caminho = Path.Combine(pastaBase, "version.json");

        try
        {
            Directory.CreateDirectory(pastaBase); // Garante que a pasta existe
            File.WriteAllText(caminho, json);
        }
        catch (UnauthorizedAccessException)
        {
            MessageBox.Show("Erro: é necessário executar o programa como administrador para salvar em Program Files.");
        }
        catch (Exception ex)
        {
            MessageBox.Show($"Erro ao salvar versão: {ex.Message}");
        }
    }


    #region Func AddFirewallRule
    /// <summary>
    /// Cria a regra no FireWall do Windows para liberar a porta 9090, que será
    /// usada na aplicação para atualizar o TechMind
    /// </summary>
    private void AddFirewallRule()
    {
        this.label2.Text = "Criando Regra do FireWall...";
        this.loader.SetProgress(90);

        string programPath = @"%ProgramFiles%\techmind\techmind.exe";
        try
        {
            ProcessStartInfo psi = new()
            {
                FileName = "netsh",
                Arguments = $"advfirewall firewall add rule name=\"TechMind\" dir=in action=allow  program=\"{programPath}\" protocol=TCP localport=8080",
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
                this.label2.Text = "✅ Regra de firewall adicionada com sucesso.";
                this.loader.SetProgress(93);
            }
            else
            {
                MessageBox.Show($"⚠️ Falha ao adicionar regra. Código: {process.ExitCode}");
                MessageBox.Show($"Erro: {error}");
            }
        }
        catch (Exception ex)
        {
            MessageBox.Show($"❌ Erro inesperado ao criar regra de firewall: {ex.Message}");
        }
    }
    #endregion

    #region OptionRestart Method
    /// <summary>
    /// Atualiza a interface do usuário para informar que a instalação 
    /// foi concluída com sucesso 
    /// e apresenta opções para reiniciar o sistema imediatamente ou depois.
    /// </summary>
    private void OptionRestart()
    {
        // Atualiza o texto do label para informar o status de sucesso.
        this.label2.Text = "TechMind instalado com sucesso!!";

        int widthButton2 = (int)(this.ClientSize.Width * 0.3); // 30% da largura
        int posXButton2 = (int)(this.ClientSize.Width * 0.2); // 20% da largura
        int posYButton2 = (int)(this.ClientSize.Height * 0.5); // 50% da altura
        this.button2 = new System.Windows.Forms.Button
        {
            Location = new System.Drawing.Point(posXButton2, posYButton2),
            Width = widthButton2,
            Text = "Reiniciar Agora"
        };
        // Associa o evento de clique ao método que lida com a reinicialização imediata.
        int posXButton2Latter = (int)(this.ClientSize.Width * 0.6); // 40% da largura
        this.button2.Click += new System.EventHandler(this.RestartNow);
        this.button3 = new System.Windows.Forms.Button
        {
            Location = new System.Drawing.Point(posXButton2Latter, posYButton2),
            Width = widthButton2,
            Text = "Reiniciar Depois"
        };
        // Associa o evento de clique ao método que lida com a reinicialização posterior.
        this.button3.Click += new System.EventHandler(this.RestartLatter);

        // Adiciona os botões configurados ao painel EULA.
        this.panelEula.Controls.Add(button2);
        this.panelEula.Controls.Add(button3);
        this.panelEula.Controls.Remove(this.loader);

        // Retoma a disposição normal dos controles e atualiza o layout.
        this.ResumeLayout(false);
        this.PerformLayout();
    }
    #endregion

    #region RestartNow Method
    /// <summary>
    /// Método responsável por reiniciar o sistema imediatamente quando o botão 
    /// correspondente for clicado.
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
}
