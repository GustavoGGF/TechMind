using System.Runtime.InteropServices;

namespace TechMindInstallerW10
{
    public partial class Program
    {
        // Importa a função ShowWindow da biblioteca user32.dll para esconder a janela de console
        [DllImport("user32.dll")]
        public static extern bool ShowWindow(IntPtr hwnd, int nCmdShow);

        [DllImport("kernel32.dll")]
        public static extern IntPtr GetConsoleWindow();

        // Constantes para a função ShowWindow
        const int SW_HIDE = 0;
        
        #region Iniciando Código
        /// <summary>
        /// Função principal que verifica os argumentos passados para determinar o modo de execução do programa.
        /// Se o argumento --silent for fornecido, o programa executa a instalação ou desinstalação silenciosa, caso contrário, executa o modo normal com a interface gráfica.
        /// </summary>
        /// <param name="args">Argumentos passados na linha de comando.</param>
        [STAThread]
        static void Main(string[] args)
        {
            // Verifica se o argumento --silent foi passado para rodar o modo silencioso
            if (args.Length > 0 && args[0].Equals("--silent", StringComparison.CurrentCultureIgnoreCase))
            {
                // Se o --silent foi passado, verifica o segundo argumento para determinar a ação
                if (args.Length > 1)
                {
                    // Se o segundo argumento for -install, chama a função de instalação silenciosa
                    if (args[1].Equals("-install", StringComparison.CurrentCultureIgnoreCase))
                    {
                        RunSilentInstallation();
                    }
                    // Se o segundo argumento for -remove, chama a função de desinstalação silenciosa
                    else if (args[1].Equals("-remove", StringComparison.CurrentCultureIgnoreCase))
                    {
                        RunSilentDesinstallation();
                    }
                    // Caso contrário, executa o modo silencioso padrão
                    else
                    {
                        RunSilentMode();
                    }
                }
                else
                {
                    // Se apenas --silent for passado, executa o modo silencioso padrão
                    RunSilentMode();
                }
            }
            else
            {
                // Caso o --silent não seja passado, executa o programa com a interface gráfica normal
                HideConsoleWindow(); // Esconde a janela do console no modo normal
                Application.EnableVisualStyles();
                Application.SetCompatibleTextRenderingDefault(false);
                Application.Run(new Main()); // Inicia o formulário da aplicação
            }
        }
        #endregion

        #region Modo Silencioso - Argumentos Faltando
        /// <summary>
        /// Exibe uma mensagem de erro informando que os argumentos necessários estão faltando.
        /// Informa como os comandos de instalação e desinstalação devem ser executados no modo silencioso.
        /// Após exibir as instruções, o aplicativo é encerrado.
        /// </summary>
        static void RunSilentMode()
        {
            // Exibe a mensagem de instrução para a instalação no modo silencioso
            Console.WriteLine("Para realizar a instalação do programa, por favor, execute o seguinte comando no prompt de comando: TechMindInstallerW10.exe --silent -install.");
            Console.WriteLine();

            // Exibe a mensagem de instrução para a desinstalação no modo silencioso
            Console.WriteLine("Para proceder com a desinstalação do programa, execute o comando a seguir no prompt de comando: TechMindInstallerW10.exe --silent -remove.");

            // Encerra a aplicação após exibir as mensagens
            Application.Exit(); 
            Environment.Exit(0);
        }
        #endregion
        
        #region Esconder Janela do Console
        /// <summary>
        /// Esta função utiliza a API do Windows para esconder a janela de console 
        /// quando o programa está sendo executado no modo normal, com a interface 
        /// gráfica sendo usada para interação.
        /// </summary>
        static void HideConsoleWindow()
        {
            IntPtr consoleWindow = GetConsoleWindow(); // Obtém o identificador da janela do console
            if (consoleWindow != IntPtr.Zero)
            {
                ShowWindow(consoleWindow, SW_HIDE); // Esconde a janela do console
            }
        }
        #endregion
    }
}
