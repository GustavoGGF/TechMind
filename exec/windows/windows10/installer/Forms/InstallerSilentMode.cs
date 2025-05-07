using System.Diagnostics;
using Microsoft.Win32;

namespace TechMindInstallerW10
{
    #region Partial Silent Instalation
    /// <summary>
    /// Essa é uma classe parcial de Program, aqui será rodado a instalação
    /// no modo silencioso
    /// </summary>
    partial class Program
    {
        #region Func RunSilentInstallation
        /// <summary>
        /// Essa função inicia a instalação no modo silencioso criando o diretório para o TechMind
        /// </summary>
        static void RunSilentInstallation()
        {
            // Exibe a mensagem de início da instalação
            Console.WriteLine("Iniciando Instalação...");

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
                    AddFirewallRule();
                }
                else
                {
                    // Se a pasta já existir, apenas chama os métodos para continuar o processo
                    _ = Get_FilesAsyncSilentAsync();
                    CreateREGEditSilent();
                    AddFirewallRule();
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

            string url = "https://techmind.lupatech.com.br/donwload-files/";  // URL do servidor para download.
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

        #region Criando Regra FireWall
        /// <summary>
        /// Essa função cria uma regra no firewall do windows para liberar a porta 9090
        /// </summary>
        static void AddFirewallRule()
        {
            Console.WriteLine("Criando Regra do FireWall...");

            try
            {
                string programPath = @"%ProgramFiles%\techmind\techmind.exe";

                ProcessStartInfo psi = new()
                {
                    FileName = "netsh",
                    Arguments = $"advfirewall firewall add rule name=\"TechMind\" dir=in program=\"{programPath}\" action=allow protocol=TCP localport=9090",
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
                    Console.WriteLine("✅ Regra de firewall adicionada com sucesso.");
                }
                else
                {
                    Console.WriteLine($"⚠️ Falha ao adicionar regra. Código: {process.ExitCode}");
                    Console.WriteLine($"Erro: {error}");
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine($"❌ Erro inesperado ao criar regra de firewall: {ex.Message}");
            }
        }
        #endregion
    }
    #endregion
}