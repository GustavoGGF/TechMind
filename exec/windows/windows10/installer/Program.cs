using System.Runtime.InteropServices;
using Microsoft.Win32;

namespace TechMindInstallerW10
{
    static class Program
    {
        // Importa a função ShowWindow da biblioteca user32.dll para esconder a janela de console
        [DllImport("user32.dll")]
        public static extern bool ShowWindow(IntPtr hwnd, int nCmdShow);
        
        [DllImport("kernel32.dll")]
        public static extern IntPtr GetConsoleWindow();
        
        // Constantes para a função ShowWindow
        const int SW_HIDE = 0;

        #region Método Principal
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
                Application.Run(new Form1()); // Inicia o formulário da aplicação
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

        #region Instalação no Modo Silencioso
        /// <summary>
        /// Inicia o processo de instalação no modo silencioso.
        /// Exibe uma mensagem informando que a instalação foi iniciada e chama o método CreateFolder para criar a pasta necessária.
        /// </summary>
        static void RunSilentInstallation()
        {
            // Exibe a mensagem de início da instalação
            Console.WriteLine("Iniciando Instalação...");
            
            // Chama o método CreateFolder para criar a pasta necessária para a instalação
            CreateFolder();
        }
        #endregion

        #region Criação da Pasta para o Programa
        /// <summary>
        /// Cria a pasta necessária para a instalação do programa em "C:\Program Files\techmind".
        /// Se a pasta já existir, continua o processo sem recriar. Em ambos os casos, chama métodos adicionais para
        /// processar arquivos e configurar o registro.
        /// </summary>
        static void CreateFolder()
        {
            // Exibe a mensagem de criação da pasta
            Console.WriteLine("Criando Pasta...");

            string folderPath = @"C:\Program Files\techmind";

            try
            {
                // Verifica se a pasta já existe
                if (!Directory.Exists(folderPath))
                {
                    // Se a pasta não existir, cria a nova pasta no caminho especificado
                    Directory.CreateDirectory(folderPath);
                    
                    // Chama o método assíncrono para obter arquivos e o método para criar o registro
                    _ = Get_FilesAsyncSilentAsync();
                    CreateREGEditSilent();
                }
                else
                {
                    // Se a pasta já existir, apenas chama os métodos para continuar o processo
                    _ = Get_FilesAsyncSilentAsync();
                    CreateREGEditSilent();
                }
            }
            catch (Exception ex)
            {
                // Exibe uma mensagem de erro se ocorrer uma exceção durante a criação do diretório
                Console.WriteLine($"Erro ao criar a pasta: {ex.Message}", "Erro", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }
        #endregion

        #region Baixar Arquivos Necessários
        /// <summary>
        /// Método assíncrono que realiza o download dos arquivos necessários para a instalação do programa.
        /// Ele envia uma requisição HTTP para a URL do servidor, verifica o código de status da resposta, 
        /// e se bem-sucedido, salva o arquivo no diretório local especificado.
        /// </summary>
        static async Task Get_FilesAsyncSilentAsync()
        {
            // Exibe a mensagem de início do download
            Console.WriteLine("Baixando arquivos de SAPPP01...");

            string url = "http://sappp01:3000/donwload-files/";  // URL do servidor para download.
            string localPath = @"C:\Program Files\techmind\techmind.exe";  // Caminho local para salvar o arquivo.

            using HttpClient client = new();
            try
            {
                // Envia a requisição GET para o servidor
                Console.WriteLine("Enviando requisição...");
                HttpResponseMessage response = client.GetAsync(url).GetAwaiter().GetResult(); // Forçando o modo síncrono

                // Exibe o código de status da resposta
                Console.WriteLine($"Código de status da resposta: {response.StatusCode}");
                
                // Verifica se a requisição foi bem-sucedida
                if (!response.IsSuccessStatusCode)
                {
                    // Exibe falha se a resposta não for bem-sucedida
                    Console.WriteLine($"Falha na requisição: {response.ReasonPhrase}");
                }

                // Verifica o status da resposta e lança exceção em caso de erro
                Console.WriteLine("Verificando status da resposta...");
                response.EnsureSuccessStatusCode();  // Lança uma exceção se o código de status for diferente de 2xx

                // Lê o conteúdo da resposta como bytes
                byte[] fileBytes = await response.Content.ReadAsByteArrayAsync();
                Console.WriteLine("Lendo bytes...");

                // Salva os bytes no caminho local especificado
                await File.WriteAllBytesAsync(localPath, fileBytes);               
            }
            catch (HttpRequestException httpEx)
            {
                // Exibe erro de requisição HTTP caso ocorra
                Console.WriteLine($"Erro de requisição HTTP: {httpEx.Message}");
            }
            catch (TaskCanceledException)
            {
                // Informa erro caso a requisição tenha sido cancelada ou expirado
                Console.WriteLine("A requisição expirou. Verifique a conexão com o servidor.");
            }
            catch (Exception ex)
            {
                // Exibe erro genérico para exceções não tratadas
                Console.WriteLine($"Erro geral: {ex.Message}");
            }
        }
        #endregion

        #region Criar Registro no Windows
        /// <summary>
        /// Este método cria uma entrada no registro do Windows para garantir que o aplicativo 
        /// seja executado automaticamente durante o login do usuário. O registro é salvo 
        /// em uma chave específica de inicialização no Windows.
        /// </summary>
        static void CreateREGEditSilent()
        {
            // Exibe a mensagem indicando que o registro está sendo criado
            Console.WriteLine("Criando RegEdit...");

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
                    Console.WriteLine("TechMind instalado com sucesso.");
                    Console.WriteLine("");
                    Console.WriteLine("Necessário reiniciar o computador.");
                    Application.Exit(); 
                    Environment.Exit(0);
                }
                else
                {
                    // Exibe uma mensagem se a chave do Registro não for encontrada.
                    Console.WriteLine("Erro ao acessar a chave do registro. A chave pode não existir.");
                }
            }
            catch (Exception ex)
            {
                // Exibe uma mensagem de erro se ocorrer uma exceção durante a operação.
                Console.WriteLine($"Erro: {ex.Message}");
            }
        }
        #endregion

        static void RunSilentDesinstallation()
        {
            Console.WriteLine("Iniciando Desinstalação...");
            Application.Exit(); 
            Environment.Exit(0);
        }

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
