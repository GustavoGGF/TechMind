namespace TechMindInstallerW10
{
    using System;
    using System.Drawing;
    using System.Windows.Forms;

    #region Componente LoaderControl
    /// <summary>
    /// Componente customizado que exibe um controle de carregamento (loader) com uma barra de progresso.
    /// Ele permite atualizar o progresso e exibe o progresso visualmente como uma barra de cor verde sobre um fundo cinza.
    /// O progresso é representado por um valor de porcentagem, que é ajustado e renderizado conforme o valor muda.
    /// </summary>
    public class LoaderControl : Control
    {
        // Variável para armazenar a porcentagem de progresso, inicializada com 10%
        private float progressPercentage = 0.1f; 

        /// <summary>
        /// Define o valor de progresso a ser exibido na barra de progresso.
        /// O valor é clamped entre 0 e 100.
        /// </summary>
        /// <param name="percentage">Porcentagem de progresso a ser exibida (0 a 100).</param>
        public void SetProgress(float percentage)
        {
            // Garante que o valor de porcentagem esteja dentro do intervalo permitido (0-100)
            if (percentage < 0f) percentage = 0f;
            if (percentage > 100f) percentage = 100f;

            // Converte a porcentagem para um valor entre 0 e 1
            progressPercentage = percentage / 100f;
            
            // Redesenha o controle para refletir o novo progresso
            Invalidate(); 
        }

        /// <summary>
        /// Redefine o desenho do controle, desenhando a barra de progresso.
        /// O método de pintura é chamado automaticamente quando a tela precisa ser atualizada.
        /// </summary>
        /// <param name="e">Evento de pintura.</param>
        protected override void OnPaint(PaintEventArgs e)
        {
            base.OnPaint(e);

            // Define o tamanho da barra de progresso
            int width = 700;
            int height = 50;

            // Desenha o retângulo de fundo (cinza)
            e.Graphics.FillRectangle(Brushes.Gray, 0, 0, width, height);

            // Desenha o retângulo de progresso (verde) baseado no valor de progressPercentage
            e.Graphics.FillRectangle(Brushes.Green, 0, 0, (int)(width * progressPercentage), height);

            // Desenha a borda preta ao redor da barra de progresso
            e.Graphics.DrawRectangle(Pens.Black, 0, 0, width, height);
        }

        /// <summary>
        /// Garante que o tamanho do controle seja fixo (800x50) ao redimensioná-lo.
        /// </summary>
        /// <param name="e">Evento de redimensionamento.</param>
        protected override void OnResize(EventArgs e)
        {
            base.OnResize(e);
            
            // Mantém o tamanho fixo para o controle
            Size = new Size(800, 50); 
        }
    }
    #endregion

}
