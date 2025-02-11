namespace TechMindInstallerW10;

partial class Form1
{
    /// <summary>
        /// Armazena os componentes do formulário que precisam ser descartados para liberar recursos.
    /// </summary>
    
    private System.ComponentModel.IContainer components = null;
    private System.Windows.Forms.Panel panelEula;
    private System.Windows.Forms.RichTextBox rtextb1; 
    private System.Windows.Forms.Label label1;
    private System.Windows.Forms.Label label2;
    private System.Windows.Forms.Label label3;
    private System.Windows.Forms.CheckBox checkBox1;
    private System.Windows.Forms.Button button1;
    private System.Windows.Forms.Button button2;
    private System.Windows.Forms.Button button3;
    private System.Windows.Forms.Button button4;
    private System.Windows.Forms.TextBox textBox1;
    private LoaderControl loader;

    /// <summary>
        /// Libera os recursos usados pelo formulário.
    /// </summary>
    /// <param name="disposing">
    /// Indica se os recursos gerenciados (true) ou não gerenciados (false) devem ser descartados.
    /// </param>
    

    protected override void Dispose(bool disposing)
    {
        // Verifica se o descarte de recursos gerenciados foi solicitado.
        if (disposing && (components != null))
        {
            // Descarta os componentes gerenciados.
            components.Dispose();
        }

        // Chama o método base para liberar recursos não gerenciados.
        base.Dispose(disposing);
    }
    


    #region Windows Form Designer generated code
    /// <summary>
    /// Método gerado automaticamente para inicializar os componentes do formulário.
    /// Este código é responsável por configurar as propriedades visuais do formulário.
    /// </summary>
    private void InitializeComponent()
    {
        // Cria uma nova instância de um contêiner de componentes.
        this.components = new System.ComponentModel.Container();
        
        // Define o modo de escalonamento automático para o formulário.
        this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
        
        // Define o tamanho do formulário.
        this.ClientSize = new System.Drawing.Size(800, 450);
        
        // Desativa a capacidade de maximizar a janela do formulário.
        this.MaximizeBox = false;
        
        // Define o estilo da borda do formulário como fixa (não pode ser redimensionada).
        this.FormBorderStyle = FormBorderStyle.FixedSingle;
        
        // Define o título do formulário.
        this.Text = "Instalação TechMind";
        
        // Define o ícone do formulário, utilizando o arquivo especificado.
        // this.Icon = new System.Drawing.Icon("./assets/images/logo.ico");
    }
    #endregion

    #region EULA Confirmation
    /// <summary>
    /// Configura a tela de confirmação do EULA (Contrato de Licença de Usuário Final).
    /// Adiciona os componentes necessários ao painel, incluindo texto, checkbox, botão e label.
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
        this.rtextb1.Size = new Size(400, 350);
        this.rtextb1.Location = new System.Drawing.Point(200, 0);
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
                \pard\cf0 Versão Atual: 1.0.0\par
            }";

        // Cria e configura um label para a opção de concordância.
        this.label1 = new System.Windows.Forms.Label();
        this.label1.Location = new System.Drawing.Point(220, 380);
        this.label1.Text = "Concordo.";
        this.label1.Click += new System.EventHandler(this.Label1_Click);

        // Cria e configura um checkbox para o usuário concordar com o EULA.
        this.checkBox1 = new System.Windows.Forms.CheckBox();
        this.checkBox1.Location = new System.Drawing.Point(200, 380);
        this.checkBox1.Size = new Size(20, 20);
        this.checkBox1.Click += new System.EventHandler(this.Button1_Click);

        // Cria e configura um botão para avançar, que inicialmente está desabilitado.
        this.button1 = new System.Windows.Forms.Button();
        this.button1.Location = new System.Drawing.Point(520, 380);
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

    }
    #endregion

    #region Create_Folder Method
    /// <summary>
    /// Configura os controles da interface do usuário para indicar a criação de uma pasta e chama a função para criá-la.
    /// </summary>
    private void Create_Folder()
    {
        loader = new LoaderControl
        {
            Location = new Point(40, 200),
            Size = new Size(200, 50) 
        };

        // Suspende temporariamente a disposição dos controles para melhor desempenho durante a modificação.
        this.SuspendLayout();
        
        // Remove controles existentes do painel EULA.
        this.panelEula.Controls.Remove(rtextb1);
        this.panelEula.Controls.Remove(checkBox1);
        this.panelEula.Controls.Remove(button1);
        this.panelEula.Controls.Remove(label1);
        
        #region Label Configuration
        /// <summary>
        /// Configuração de um novo rótulo (label) para exibir o status de criação da pasta.
        /// </summary>
        this.label2 = new System.Windows.Forms.Label
        {
            Location = new System.Drawing.Point(200, 10),
            Text = "Criando Pasta...",
            Width = 200,  // Define a largura do rótulo.
            Height = 12  // Define a altura do rótulo.
        };
        #endregion

        #region TextBox Configuration
        /// <summary>
        /// Configuração de uma nova caixa de texto para exibir o caminho do diretório a ser criado.
        /// </summary>
        this.textBox1 = new System.Windows.Forms.TextBox
        {
            Location = new System.Drawing.Point(200, 40),
            Size = new Size(200, 20),
            ReadOnly = true,  // Define a caixa de texto como somente leitura.
            Text = "C:\\Program Files\\TechMind",  // Caminho padrão para a pasta a ser criada.
            AllowDrop = false,  // Desativa o recurso de arrastar e soltar.
            Enabled = false  // Desabilita a edição da caixa de texto.
        };
        #endregion

        // Adiciona os novos controles ao painel.
        this.panelEula.Controls.Add(textBox1);
        this.panelEula.Controls.Add(label2);
        this.panelEula.Controls.Add(this.loader);

        // Retoma a disposição normal dos controles e atualiza o layout.
        this.ResumeLayout(false);
        this.PerformLayout();
        loader.SetProgress(10); 

        // Chama o método responsável por criar a pasta.
        Make_Folder(this);
    }
    #endregion

    #region OptionRestart Method
    /// <summary>
    /// Atualiza a interface do usuário para informar que a instalação foi concluída com sucesso 
    /// e apresenta opções para reiniciar o sistema imediatamente ou depois.
    /// </summary>
    private void OptionRestart()
    {
        // Atualiza o texto do label para informar o status de sucesso.
        this.label2.Text = "TechMind instalado com sucesso!!";

        #region Button Configuration - Restart Now
        /// <summary>
        /// Configura um botão para reiniciar o sistema imediatamente.
        /// </summary>
        this.button2 = new System.Windows.Forms.Button
        {
            Location = new System.Drawing.Point(200, 250),
            Width = 200,
            Text = "Reiniciar Agora"
        };
        // Associa o evento de clique ao método que lida com a reinicialização imediata.
        this.button2.Click += new System.EventHandler(this.RestartNow);
        #endregion

        #region Button Configuration - Restart Later
        /// <summary>
        /// Configura um botão para reiniciar o sistema posteriormente.
        /// </summary>
        this.button3 = new System.Windows.Forms.Button
        {
            Location = new System.Drawing.Point(520, 250),
            Width = 200,
            Text = "Reiniciar Depois"
        };
        // Associa o evento de clique ao método que lida com a reinicialização posterior.
        this.button3.Click += new System.EventHandler(this.RestartLatter);
        #endregion

        // Adiciona os botões configurados ao painel EULA.
        this.panelEula.Controls.Add(button2);
        this.panelEula.Controls.Add(button3);
        this.panelEula.Controls.Remove(this.loader);

        // Retoma a disposição normal dos controles e atualiza o layout.
        this.ResumeLayout(false);
        this.PerformLayout();
    }
    #endregion

    private void UninstallationConfirmation()
    {
        // Cria uma nova instância de um contêiner de componentes.
        this.components = new System.ComponentModel.Container();
        
        // Define o modo de escalonamento automático para o formulário.
        this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
        
        // Define o tamanho do formulário.
        this.ClientSize = new System.Drawing.Size(800, 450);
        
        // Desativa a capacidade de maximizar a janela do formulário.
        this.MaximizeBox = false;
        
        // Define o estilo da borda do formulário como fixa (não pode ser redimensionada).
        this.FormBorderStyle = FormBorderStyle.FixedSingle;
        
        // Define o título do formulário.
        this.Text = "Desinstalação TechMind";

        loader = new LoaderControl
        {
            Location = new Point(40, 200),
            Size = new Size(200, 50) 
        };

        // Cria e configura um label para a opção de concordância.
        this.label3 = new System.Windows.Forms.Label();
        this.label3.Location = new System.Drawing.Point(200, 10);
        this.label3.Width = 400;
        this.label3.Height = 50;
        this.label3.Font = new Font(label3.Font, FontStyle.Bold);
        this.label3.Text = "Você esta prestes a desinstalar o TechMind, caso esteja ciente disso clique em avançar para dar andamento.";

        this.button4 = new System.Windows.Forms.Button();
        this.button4.Location = new System.Drawing.Point(300, 200);
        this.button4.Height = 50;
        this.button4.Width = 80;
        this.button4.Text = "Prosseguir";
        this.button4.Click += new System.EventHandler(this.RemoveRegEdit);

        this.Controls.Add(label3);
        this.Controls.Add(button4);
        this.ResumeLayout(false);
        this.PerformLayout();
    }
}
