using Microsoft.Win32;

namespace TechMindInstallerW10;

#region Função Principal do Código
/// <summary>
/// Main é a função Principal, tudo começa nela
/// Main é Publico, pode ser chamado a partir de qualquer lugar
/// </summary>

public partial class Main : Form
{
    #region Func Main
    /// <summary>
    /// Essa função ao ser iniciaada Chama SoftwareExistenceCheck para fazer
    /// a validação se o TechMind já está istalado
    /// </summary>
    public Main()
    {
        // Verifica se o software já está no registro
        SoftwareExistenceCheck();
    }
    #endregion

    #region Func SoftwareExistenceCheck
    /// <summary>
    /// Essa função faz a verificação se TechMind está instalado ou Não
    /// A verificação é feita atravez do Registro que o Mesmo gera ao ser instalado
    /// </summary>
    private void SoftwareExistenceCheck()
    {
        try
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
        catch (Exception ex)
        {
            MessageBox.Show("Erro na função SoftwareExistenceCheck: " + ex.Message);
        }
    }
    #endregion

}
#endregion