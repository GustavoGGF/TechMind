namespace TechMindInstallerW10;

#region Partial Main
/// <summary>
/// Uma partição de Main, esse arquivo deve se manter apenas para gerar a tela de instalação
/// </summary>
partial class Main
{
    private System.ComponentModel.IContainer components = null;
    #region Func InitializeComponent
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
    }

    #endregion

}
#endregion