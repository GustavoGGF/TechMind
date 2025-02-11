using System.Windows.Forms;

namespace Utils
{
    public class Loader
    {
        private System.Windows.Forms.Panel barraRetangular; // Painel que representará a barra retangular
        private const int MaxWidth = 280; // Largura máxima da barra retangular

        public void Loading(Form form)
        {
            // Cria uma nova instância do painel (barra retangular)
            this.barraRetangular = new System.Windows.Forms.Panel();

            // Define a posição da barra retangular no formulário (x=10, y=90)
            this.barraRetangular.Location = new System.Drawing.Point(10, 90);

            // Define o tamanho inicial da barra retangular (0 pixels de largura e 30 pixels de altura)
            this.barraRetangular.Size = new System.Drawing.Size(0, 30);

            // Define a cor de fundo da barra retangular (pode ser alterada conforme desejado)
            this.barraRetangular.BackColor = System.Drawing.Color.Green; // Exemplo: cor verde

            // Adiciona a barra retangular ao formulário
            form.Controls.Add(this.barraRetangular);
        }

        // Método para atualizar o progresso da barra retangular
        public void UpdateProgress(int percentage)
        {
            // Garante que o valor percentual esteja entre 0 e 100
            if (percentage < 0) percentage = 0;
            if (percentage > 100) percentage = 100;

            // Calcula a nova largura da barra com base na porcentagem
            int newWidth = (MaxWidth * percentage) / 100;

            // Atualiza o tamanho da barra retangular
            this.barraRetangular.Size = new System.Drawing.Size(newWidth, 30);
        }
    }
}
