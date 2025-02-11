using System;
using System.Windows.Forms;
using System.ComponentModel;
using System.Drawing;
using System.IO;
using System.Diagnostics;

namespace techmind
{
    // Definição parcial da classe Form1
    partial class Form1 : Form
    {
        // Container para gerenciar componentes da interface
        private System.ComponentModel.IContainer components = null;
        // Componente PictureBox para exibir a imagem
        private System.Windows.Forms.PictureBox pictureBox1;
        // Componente Label para exibir um texto
        private System.Windows.Forms.Label label1;
        private System.Windows.Forms.Label label2;
        // Componente Button para iniciar a instalação
        private System.Windows.Forms.Button button1;
        // Componente Textbox mostrando local da instalação
        private System.Windows.Forms.TextBox textBox1; 
        // Componente CheckBox concordando com os termos
        private System.Windows.Forms.CheckBox checkBox1;

        // Método responsável por liberar recursos da interface gráfica
        protected override void Dispose(bool disposing)
        {
            // Se for necessário liberar recursos
            if (disposing && (components != null))
            {
                // Libera os componentes gerenciados
                components.Dispose();
            }
            // Chama o método base para liberar outros recursos
            base.Dispose(disposing);
        }

        // Método que inicializa os componentes da interface gráfica
        private void InitializeComponent()
        {
            // Instancia o componente PictureBox
            this.pictureBox1 = new System.Windows.Forms.PictureBox();
            // Instancia o componente Label
            this.label1 = new System.Windows.Forms.Label();
            this.label2 = new System.Windows.Forms.Label();
            // Instancia o componente Button
            this.button1 = new System.Windows.Forms.Button();
            // Instancia o componente TextBox
            this.textBox1 = new System.Windows.Forms.TextBox();
            // Instancia o componente CheckBox
            this.checkBox1 = new System.Windows.Forms.CheckBox();
            // Suspende temporariamente a atualização do layout
            this.SuspendLayout();

            // Define a posição do PictureBox
            this.pictureBox1.Location = new System.Drawing.Point(5, 5);
            // Nomeia o PictureBox
            this.pictureBox1.Name = "pictureBox1";
            // Define o tamanho do PictureBox
            this.pictureBox1.Size = new System.Drawing.Size(20, 20);
            // Configura o modo de exibição da imagem para esticar a imagem no PictureBox
            this.pictureBox1.SizeMode = System.Windows.Forms.PictureBoxSizeMode.StretchImage;

            // Define a posição da label (x=60, y=5)
            this.label1.Location = new System.Drawing.Point(60, 5); 
            // Nomeia a label como "label1" para referenciá-la no código
            this.label1.Name = "label1";
            // Define que o tamanho da label será ajustado automaticamente conforme o texto
            this.label1.AutoSize = true;
            // Define o texto que será exibido na label
            this.label1.Text = "Instalação do TechMind"; 
            // Define a fonte da label como "Segoe UI", tamanho 10, com estilo negrito (bold)
            this.label1.Font = new System.Drawing.Font("Segoe UI", 10, FontStyle.Bold); 
            // Define a cor do texto da label como roxa
            this.label1.ForeColor = System.Drawing.Color.Purple; 
            // Define o fundo da label como transparente
            this.label1.BackColor = System.Drawing.Color.Transparent;

            // Define a posição da label2 (x=10, y=30)
            this.label2.Location = new System.Drawing.Point(10, 30); 
            // Nomeia a label como "label2" para referenciá-la no código
            this.label2.Name = "label2";
            // Define o tamanho da label2 como 280 pixels de largura e 80 pixels de altura
            this.label2.Size = new System.Drawing.Size(280, 80); 
            // Define o texto que será exibido na label2, descrevendo a ferramenta TechMind
            this.label2.Text = "TechMind é a ferramenta de inventário da Lupatech, ao ser instalada ela extrai informações como Hardware, Software conforme as necessidades da empresa.";
            // Define a fonte da label2 como "Segoe UI", tamanho 10, com estilo negrito (bold)
            this.label2.Font = new System.Drawing.Font("Segoe UI", 10, FontStyle.Bold); 
            // Define a cor do texto da label2 como preta
            this.label2.ForeColor = System.Drawing.Color.Black; 
            // Define o fundo da label2 como transparente
            this.label2.BackColor = System.Drawing.Color.Transparent;

            // Define a posição do botão (x=10, y=120)
            this.button1.Location = new System.Drawing.Point(10, 120);
            // Nomeia o botão como "button1" para referenciá-lo no código
            this.button1.Name = "button1";
            // Define o tamanho do botão como 65 pixels de largura e 35 pixels de altura
            this.button1.Size = new System.Drawing.Size(65, 35);
            // Define a ordem de tabulação do botão na interface, permitindo navegação com a tecla Tab
            this.button1.TabIndex = 0;
            // Define o texto que será exibido no botão
            this.button1.Text = "Instalar TechMind";
            // Define que o botão utilizará o estilo visual padrão do Windows
            this.button1.UseVisualStyleBackColor = true;
            // Define que o botão está desativado (não pode ser clicado) inicialmente
            this.button1.Enabled = false;
            // Adiciona um manipulador de eventos para o evento Click do botão, chamando o método button1_Click quando o botão é clicado
            this.button1.Click += new System.EventHandler(this.button1_Click);

            // Tenta carregar a imagem para o PictureBox
            try
            {
                // Caminho da imagem, combinando o caminho de execução com a pasta e o arquivo
                string imagePath = System.IO.Path.Combine(Application.StartupPath, "assets", "logo.png");
                // Carrega a imagem do arquivo especificado
                this.pictureBox1.Image = Image.FromFile(imagePath);
            }
            catch (Exception ex)
            {
                // Exibe uma mensagem de erro se houver falha no carregamento da imagem
                MessageBox.Show($"Erro ao carregar a imagem: {ex.Message}", "Erro", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }

            // Define a posição da textBox1 (x=80, y=125)
            this.textBox1.Location = new System.Drawing.Point(80, 125);
            // Nomeia a textBox como "textBox1" para referenciá-la no código
            this.textBox1.Name = "textBox1";
            // Define o tamanho da textBox como 200 pixels de largura e 10 pixels de altura
            this.textBox1.Size = new System.Drawing.Size(200, 10);
            // Define a ordem de tabulação da textBox na interface, permitindo navegação com a tecla Tab
            this.textBox1.TabIndex = 1;
            // Define o texto que será exibido na textBox, especificando um caminho padrão
            this.textBox1.Text = @"C:\Program Files\techmind\";
            // Define que a textBox é somente leitura, ou seja, o usuário não pode editar seu conteúdo
            this.textBox1.ReadOnly = true;

            // Define a posição do checkBox1 (x=0, y=180)
            this.checkBox1.Location = new System.Drawing.Point(0, 180); 
            // Nomeia o checkBox como "checkBox1" para referenciá-lo no código
            this.checkBox1.Name = "checkBox1";
            // Define o tamanho do checkBox como 280 pixels de largura e 30 pixels de altura
            this.checkBox1.Size = new System.Drawing.Size(280, 30); 
            // Define o texto que será exibido ao lado do checkBox, informando sobre a extração de dados sensíveis
            this.checkBox1.Text = "Estou ciente que serão extraindo dados sensiveis da minha maquina e concordo.";
            // Adiciona um manipulador de eventos para o evento Click do checkBox, chamando o método checkBox1_Click quando o checkBox é clicado
            this.checkBox1.Click += new System.EventHandler(this.checkBox1_Click);

            // Inicializa o container de componentes
            this.components = new System.ComponentModel.Container();
            // Define as dimensões de escala automática da interface, utilizando um fator de 6F (largura) e 13F (altura)
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            // Define o modo de escala automática como fonte, para que a interface ajuste a fonte de acordo com a DPI do sistema
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            // Define o tamanho do cliente da janela como 300 pixels de largura e 220 pixels de altura
            this.ClientSize = new System.Drawing.Size(300, 220);
            // Define o estilo da borda do formulário como fixa (não pode ser redimensionado pelo usuário)
            this.FormBorderStyle = System.Windows.Forms.FormBorderStyle.FixedSingle; 
            // Desativa a capacidade de maximizar a janela, impedindo que o usuário maximize o formulário
            this.MaximizeBox = false; 

            // Adicioando componentes criados
            this.Controls.Add(this.pictureBox1);
            this.Controls.Add(this.label1);
            this.Controls.Add(this.label2);
            this.Controls.Add(this.button1);
            this.Controls.Add(this.textBox1);
            this.Controls.Add(this.checkBox1);
            // Define o texto da janela principal
            this.Text = "Termos de instalação";
        }

        // Método que realiza a instalação, removendo controles da interface do usuário
        private void Installing()
        {
            // Remove os controles do formulário, tornando-os invisíveis
            this.Controls.Remove(this.label1);
            this.Controls.Remove(this.label2);
            this.Controls.Remove(this.button1);
            this.Controls.Remove(this.textBox1);
            this.Controls.Remove(this.checkBox1);

            // Inicializa o componente label1 e label2 como uma nova instância do controle Label
            this.label1 = new System.Windows.Forms.Label();
            this.label2 = new System.Windows.Forms.Label();
            // Suspende temporariamente o layout para evitar atualizações durante a configuração dos controles
            this.SuspendLayout();

            // Define a posição da label1 no formulário (x=60, y=5)
            this.label1.Location = new System.Drawing.Point(60, 5);
            // Nomeia o controle label1 para permitir sua referência no código
            this.label1.Name = "label1";
            // Define a propriedade AutoSize como true para que o tamanho da label se ajuste automaticamente ao seu conteúdo
            this.label1.AutoSize = true;
            // Define o texto que será exibido na label1, informando o título "Instalação TechMind"
            this.label1.Text = "Instalação TechMind";
            // Define a fonte da label1 como "Segoe UI", tamanho 10, com estilo negrito (Bold)
            this.label1.Font = new System.Drawing.Font("Segoe UI", 10, FontStyle.Bold);
            // Define a cor do texto da label1 como roxa (Purple)
            this.label1.ForeColor = System.Drawing.Color.Purple;
            // Define o fundo da label1 como transparente, integrando-se ao fundo do formulário
            this.label1.BackColor = System.Drawing.Color.Transparent;

            // Define a posição da label2 no formulário (x=60, y=70)
            this.label2.Location = new System.Drawing.Point(10, 70);
            // Nomeia o controle label2 para permitir sua referência no código
            this.label2.Name = "label2";
            // Define a propriedade AutoSize como true para que o tamanho da label se ajuste automaticamente ao seu conteúdo
            this.label2.AutoSize = true;
            // Define o texto que será exibido na label2, indicando que a pasta está sendo criada
            this.label2.Text = "Criando pasta...";
            // Define a fonte da label2 como "Segoe UI", tamanho 10, com estilo negrito (Bold)
            this.label2.Font = new System.Drawing.Font("Segoe UI", 10, FontStyle.Bold);
            // Define a cor do texto da label2 como preta (Black)
            this.label2.ForeColor = System.Drawing.Color.Black;
            // Define o fundo da label2 como transparente, integrando-se ao fundo do formulário
            this.label2.BackColor = System.Drawing.Color.Transparent;

            loader.Loading(this);

            // Adiciona o controle label1\label2 ao formulário, para que seja exibido na interface
            this.Controls.Add(this.label1);
            this.Controls.Add(this.label2);
            // Define o texto do título da janela do formulário como "Instalando..."
            this.Text = "Instalando...";
            // Retoma o layout do formulário, aplicando as mudanças configuradas
            this.ResumeLayout(false);

            // Iniciando função para instalar de fato TechMind
            CreateFolder();
        }

        private void Moving()
        {
            // Caminho de origem (servidor)
            string sourcePath = @"\\snas01\node00\node5\lun0\f-sti01\002-Programas\DEP_Suporte\techmind\windows7\techmind.exe";
            
            // Caminho de destino
            string destinationDirectory = @"C:\Program Files\techmind";
            string destinationPath = Path.Combine(destinationDirectory, "tecmind.exe");

            try
            {
                // Verifica se o diretório de destino existe, se não, cria-o
                if (!Directory.Exists(destinationDirectory))
                {
                    Directory.CreateDirectory(destinationDirectory);
                }

                // Copia o arquivo
                File.Copy(sourcePath, destinationPath, true);

                this.label2.Text = "Instalando...";
                this.loader.UpdateProgress(70);

                // Caminho do arquivo .exe que você deseja executar
                string filePath = @"C:\Program Files\techmind\tecmind.exe";

                try
                {
                    // Cria um novo processo para executar o arquivo
                    Process process = new Process();
                    process.StartInfo.FileName = filePath;

                    // Configurações opcionais
                    process.StartInfo.UseShellExecute = false; // Usa o shell padrão do sistema
                    process.StartInfo.CreateNoWindow = true; // Não cria uma janela de console

                    // Inicia o processo
                    process.Start();

                    // Aguarda o processo terminar, se necessário
                    process.WaitForExit();
                    this.loader.UpdateProgress(100);
                }
                catch (Exception ex)
                {
                    MessageBox.Show("Erro ao executar o arquivo: " + ex.Message);
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show("Erro ao copiar o arquivo: " + ex.Message);
            }
    
        }
    }

    public static class InstallerHelperSilent
    {
        public static void CreateFolderSilent()
        {
            Console.WriteLine("Criando Diretório.....");
            // Caminho da nova pasta
            string folderPath = @"C:\Program Files\techmind";

            try
            {
                // Verifica se a pasta já existe
                if (!Directory.Exists(folderPath))
                {
                    // Cria a nova pasta
                    Directory.CreateDirectory(folderPath);
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Erro ao criar a pasta: {ex.Message}", "Erro", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
            
        }

        public static void MoveFilesSilent()
        {
            Console.WriteLine("Movendo Arquivos....");
                        // Caminho de origem (servidor)
            string sourcePath = @"\\snas01\node00\node5\lun0\f-sti01\002-Programas\DEP_Suporte\techmind\windows7\techmind.exe";
            
            // Caminho de destino
            string destinationDirectory = @"C:\Program Files\techmind";
            string destinationPath = Path.Combine(destinationDirectory, "tecmind.exe");

            try
            {
                // Verifica se o diretório de destino existe, se não, cria-o
                if (!Directory.Exists(destinationDirectory))
                {
                    Directory.CreateDirectory(destinationDirectory);
                }

                // Copia o arquivo
                File.Copy(sourcePath, destinationPath, true);

                Console.WriteLine("Instalando...");

                // Caminho do arquivo .exe que você deseja executar
                string filePath = @"C:\Program Files\techmind\tecmind.exe";
                try
                {
                    // Cria um novo processo para executar o arquivo
                    Process process = new Process();
                    process.StartInfo.FileName = filePath;

                    // Configurações opcionais
                    process.StartInfo.UseShellExecute = false; // Usa o shell padrão do sistema
                    process.StartInfo.CreateNoWindow = true; // Não cria uma janela de console

                    // Inicia o processo
                    process.Start();

                    // Aguarda o processo terminar, se necessário
                    process.WaitForExit();
                }
                catch (Exception ex)
                {
                    Console.WriteLine("Erro ao executar o arquivo: " + ex.Message);
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine("Erro ao copiar o arquivo: " + ex.Message);
            }
        }
    }

}
