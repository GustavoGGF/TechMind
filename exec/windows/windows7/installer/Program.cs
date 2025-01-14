using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Windows.Forms;
using System.Runtime.InteropServices;

namespace techmind
{
    static class Program
    {
        /// <summary>
        ///  The main entry point for the application.
        /// </summary>
        [DllImport("kernel32.dll", SetLastError = true)]
        [return: MarshalAs(UnmanagedType.Bool)]
        static extern bool AllocConsole();

        [STAThread]
        static void Main(string[] args)
        {
            bool silentMode = args.Length > 0 && args[0].Equals("-silent", StringComparison.OrdinalIgnoreCase);

            if (silentMode)
            {
                AllocConsole();
                Installer.RunSilentInstallation();
            }
            else
            {
            Application.SetHighDpiMode(HighDpiMode.SystemAware);
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.Run(new Form1());
            }
        }
    }
}
