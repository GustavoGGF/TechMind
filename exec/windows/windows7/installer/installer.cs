using System; // Importa o namespace do sistema, onde encontramos várias funcionalidades básicas
using System.IO;
using System.Collections.Generic; // Importa suporte a coleções genéricas (como listas e dicionários)
using System.ComponentModel; // Para suporte a componentes e controles em Windows Forms
using System.Data; // Para manipulação de dados, geralmente usado com bancos de dados
using System.Drawing; // Para trabalhar com gráficos e imagens (porque às vezes queremos que nossos botões tenham um pouco de cor)
using System.Linq; // Para operações de consulta em coleções (porque escrever 'for' é tão século passado)
using System.Text; // Para manipulação de strings, essencial para quem gosta de palavras
using System.Threading.Tasks; // Para programação assíncrona e multitarefa (porque, quem tem tempo para esperar?)
using System.Windows.Forms; // Para criar aplicações de desktop com interfaces gráficas (a parte visual da nossa magia)

namespace techmind // Nome do nosso namespace, o nosso espaço pessoal para evitar conflitos de nomes
{
    public static class Installer
    {
        public static void RunSilentInstallation()
        {
            try
            {
                InstallerHelperSilent.CreateFolderSilent();
                InstallerHelperSilent.MoveFilesSilent();
                Console.WriteLine("Pressione Enter para finalizar...");
                Console.ReadLine();
            }
            catch (Exception ex)
            {
                Console.WriteLine("Erro detectado: " + ex.Message);
            }
        }
        
    }

    public partial class Form1 : Form // Definindo a classe Form1 que herda de Form, como um filho que herda a beleza dos pais
    {
        private Utils.Loader loader; // Instância do Loader, nossa fiel escudeira para tarefas de instalação

        // Construtor da classe Form1
        public Form1()
        {
            InitializeComponent(); // Inicializa os componentes do formulário (é como arrumar a mesa antes de um jantar)
            loader = new Utils.Loader(); // Cria uma nova instância do Loader (ainda não sabemos o que ele faz, mas deve ser algo legal!)
        }

        // Evento que acontece quando o button1 é clicado
        private void button1_Click(object sender, EventArgs e)
        {
            InitializeInstall(); // Chama o método que inicia a instalação (não, não é uma instalação do Windows!)
        }

        // Evento que acontece quando o checkBox1 é clicado
        private void checkBox1_Click(object sender, EventArgs e)
        {
            this.button1.Enabled = checkBox1.Checked; // Habilita ou desabilita o botão com base no estado da checkbox
            // Se a checkbox estiver marcada, button1 fica habilitado! Se não, ele vai para o timeout... ops, eu quis dizer "desabilitado".
        }

        // Método que inicializa a instalação
        private void InitializeInstall()
        {
            Installing(); // Chama o método de instalação (sinta a tensão no ar!)
        }

        // Método para criar uma pasta (tão emocionante quanto assistir tinta secar)
        private void CreateFolder()
        {
            this.loader.UpdateProgress(10); // Atualiza o progresso do Loader para 10% (porque todo mundo ama ver a barra crescendo!)
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
                MessageBox.Show($"Erro ao criar a pasta: {ex.Message}", "Erro", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }

            MoveFiles();
        }

        private void MoveFiles()
        {
            this.label2.Text = "Movendo arquivos...";
            this.loader.UpdateProgress(30);

            Moving();
        }
    }
}
