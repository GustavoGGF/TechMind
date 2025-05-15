import { Component, OnInit } from "@angular/core";
import { UtilitiesModule } from "../../utilities/utilities.module";
import { CommonModule } from "@angular/common";
import { HttpClient, HttpClientModule } from "@angular/common/http";
import "../../../assets/bootstrap-5.3.3-dist/js/bootstrap.js";
import "../../../assets/bootstrap-5.3.3-dist/js/bootstrap.bundle.min.js";
import { catchError, throwError } from "rxjs";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";

declare var bootstrap: any;
@Component({
  selector: "app-panel-adm",
  standalone: true,
  imports: [UtilitiesModule, CommonModule, HttpClientModule],
  templateUrl: "./panel-adm.component.html",
  styleUrl: "./panel-adm.component.css",
})
export class PanelAdmComponent implements OnInit {
  constructor(private http: HttpClient) {}

  popoverTriggerList: NodeListOf<Element> = [] as any;
  popoverList: any[] = [];

  machines: any;
  listMachines: any;
  menuCustomSelect: any;
  name: any;

  buttomSize = "/static/assets/images/minimize.png";
  computers_class: string = "";
  device_class: string = "";
  errorType: string = "";
  home_class: string = "";
  machineHeader: string = "Máquina";
  machineName: string = "";
  messageError: string = "";
  panel_class: string = "active";
  processAnimation: string = "";
  processExec: string = "";
  processHeader: string = "Processo";
  statusDot: string = "";
  statusDotContent: string = "";
  statusDotPing: string = "";
  statusDotTitle: string = "";
  statusPorcentage: string = "0%";
  statusHeader: string = "Status";

  canView: boolean = false;
  canViewLoadingSearch: boolean = true;
  canViewProcessTab: boolean = false;
  expenseProcessTab: boolean = true;
  dotActive: boolean = false;
  menuSingular: boolean = false;
  menuVisible: boolean = false;
  showMessage: boolean = false;

  status: number = 0;
  tabsMachines: number = 0;

  menuPosition = { x: 0, y: 0 };

  /**
   * ngOnInit é um lifecycle hook do Angular executado automaticamente
   * após a criação do componente e inicialização de suas propriedades.
   *
   * Esta função realiza três operações principais:
   * 1. Recupera o nome do usuário armazenado no navegador.
   * 2. Valida a existência desses dados e trata erros de ausência.
   * 3. Inicializa a chamada à API para obter dados e registra um ouvinte
   *    global de cliques para controlar elementos interativos como popovers e menus.
   */
  ngOnInit() {
    // Recupera o valor da chave "name" do armazenamento local do navegador
    try {
      this.name = localStorage.getItem("name");

      // Verifica se o valor de 'name' é nulo ou uma string vazia
      if (this.name.length == 0 || this.name == null) {
        // Caso os dados estejam ausentes, define tipo e mensagem de erro
        this.errorType = "Falta de Dados";
        this.messageError =
          "Ouve um erro ao acessar dados do LDAP, contatar a TI";
        // Sinaliza para a interface que a mensagem de erro deve ser exibida
        this.showMessage = true;
      } else {
        // Dados válidos: permite a visualização da interface protegida
        this.canView = true;
        // Inicia chamada à API para buscar informações das máquinas
        this.getMachines();
      }

      // Registra um ouvinte para o evento de clique global no documento
      document.addEventListener("click", (event: MouseEvent) => {
        // Se um popover estiver ativo, executa lógica de encerramento
        if (this.dotActive) {
          this.closePoPOver(event);
        }
        // Se um menu contextual estiver posicionado (visível), oculta o menu
        if (this.menuPosition) {
          this.hideMenu();
        }
      });
    } catch (err: string | any) {
      this.errorType = "Erro Inesperado";
      this.messageError = err;
      this.showMessage = true;
      return console.error(err);
    }
  }

  /**
   * Função responsável por ocultar o componente de mensagem exibido na interface.
   *
   * Define a propriedade 'showMessage' como falsa, removendo a visibilidade do alerta ou notificação.
   * Utilizada, por exemplo, após o usuário fechar manualmente a mensagem de erro ou aviso.
   */
  hideMessage() {
    this.showMessage = false;
  }

  /**
   * Função responsável por ativar um popover específico na interface.
   *
   * Antes de ativar o novo popover, a função garante que todos os outros popovers
   * com a classe "statusDot" sejam fechados, evitando múltiplos popovers abertos simultaneamente.
   *
   * Após garantir que todos os anteriores foram ocultados, ativa o popover do elemento alvo
   * identificado pelo ID fornecido.
   * Também sinaliza que um popover está ativo através da flag 'dotActive'.
   *
   * @param element - ID do elemento HTML no qual o popover deve ser ativado.
   */
  activePopOver(element: any) {
    try {
      // Obtém o elemento DOM pelo ID fornecido
      const getElement = document.getElementById(element);
      console.log(getElement);

      // Obtém ou cria uma instância do Popover Bootstrap associada ao elemento
      const popover = bootstrap.Popover.getOrCreateInstance(getElement);
      console.log(popover);

      // Seleciona todos os elementos com a classe "statusDot"
      const dots = document.querySelectorAll(".statusDot");
      console.log(dots);

      // Itera por todos os elementos com statusDot e oculta qualquer popover ativo
      dots.forEach((element) => {
        const popover = bootstrap.Popover.getInstance(element);
        if (popover) {
          popover.hide();
        }
      });

      // Marca que há um popover ativo
      this.dotActive = true;
      // Exibe o popover no elemento alvo
      popover.show();
    } catch (err: string | any) {
      this.errorType = "Erro Inesperado";
      this.messageError = err;
      this.showMessage = true;
      return console.error(err);
    }
  }

  /**
   * Função responsável por fechar todos os popovers ativos quando o usuário clica fora deles.
   *
   * Esta função é invocada a partir de um eventListener global de clique registrado no documento.
   *
   * Ela verifica se o clique foi realizado fora de um popover ou do botão de ativação (identificado por "status_dot").
   * Se o clique ocorrer em uma área externa, todos os popovers com a classe "statusDot" são ocultados.
   *
   * Em caso de erro inesperado durante a execução (ex: manipulação de DOM), a função captura a exceção,
   * exibe uma mensagem de erro ao usuário e registra o erro no console.
   *
   * @param event - Evento de clique do mouse capturado globalmente no documento.
   */
  closePoPOver(event: MouseEvent) {
    try {
      // Obtém o elemento que foi clicado
      const target = event.target as HTMLElement;

      // Verifica se o clique ocorreu fora do botão (id com "status_dot") e fora das áreas internas do popover
      if (
        !target.classList.contains("dot_validate") &&
        !(
          target.classList.contains("popover-header") ||
          target.classList.contains("popover-body")
        )
      ) {
        // Seleciona todos os elementos com classe "statusDot"
        const dots = document.querySelectorAll(".statusDot");

        // Para cada dot, se houver um popover associado, ele será ocultado
        dots.forEach((element) => {
          const popover = bootstrap.Popover.getInstance(element);
          if (popover) {
            // Atualiza o estado do componente informando que nenhum popover está ativo
            this.dotActive = false;
            popover.hide();
          }
        });
      }
    } catch (err: string | any) {
      // Em caso de erro, define o tipo e a mensagem de erro e registra no console
      this.errorType = "Erro Inesperado";
      this.messageError = err;
      this.showMessage = true;
      return console.error(err);
    }
  }

  /**
   * Função responsável por buscar a lista de máquinas através de uma requisição HTTP GET.
   *
   * A requisição é feita para o endpoint "/home/panel-adm/get-machines".
   * Caso ocorra algum erro na chamada, ele é capturado e tratado com uma mensagem de erro apropriada.
   *
   * Ao receber uma resposta válida, os dados são processados por uma função externa chamada `groupSplitter`,
   * que organiza os resultados em grupos de 100 elementos para facilitar a exibição ou manipulação paginada.
   *
   * Após o processamento, a flag `canViewLoadingSearch` é desativada para indicar o término do carregamento.
   */
  getMachines() {
    this.http
      // Executa requisição GET para obter os dados das máquinas
      .get("/home/panel-adm/get-machines", {})
      .pipe(
        // Intercepta e trata erros da requisição HTTP
        catchError((error) => {
          this.errorType = "Erro de Conexão";
          this.messageError = "Erro ao fazer fetch para get-machines: " + error;
          console.error(error);
          this.showMessage = true;
          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        // Se os dados forem recebidos com sucesso
        if (data) {
          // Chama a função externa para agrupar os dados das máquinas em blocos de 100
          this.machines = this.groupSplitter(data.machines[0], 100);
          this.listMachines = this.machines[0];
          this.tabsMachines = Object.keys(this.machines).length;

          // Atualiza o estado indicando que o carregamento foi concluído
          return (this.canViewLoadingSearch = false);
        }
        return;
      });
  }

  /**
   * Função genérica que divide um array em subarrays de tamanho definido, facilitando o processamento em grupos.
   *
   * Utilizada para segmentar grandes conjuntos de dados em "pedaços" menores, por exemplo, para paginação ou
   * processamento em blocos.
   *
   * @template T - Tipo genérico dos elementos do array.
   * @param array - Array original que será dividido em grupos.
   * @param size - Número máximo de elementos em cada grupo (tamanho dos subarrays).
   * @returns Array de arrays, onde cada subarray possui até 'size' elementos.
   */
  groupSplitter<T>(array: T[], size: number): T[][] {
    const groups: T[][] = [];
    // Percorre o array original avançando 'size' elementos a cada iteração
    for (let i = 0; i < array.length; i += size) {
      // Cria um subarray do índice atual até o índice atual + size (limite)
      groups.push(array.slice(i, i + size));
    }
    // Retorna o array contendo os grupos segmentados
    return groups;
  }

  /**
   * Função que ajusta a exibição do status de conexão de uma máquina com base na data da última atividade.
   *
   * Recebe uma string representando a data/hora da última conexão da máquina e calcula a diferença em relação ao momento atual.
   * Dependendo se a máquina esteve ativa nas últimas 48 horas, configura os estilos e títulos para o indicador visual (dot) correspondente.
   *
   * - Se a última conexão ocorreu dentro das últimas 48 horas, indica que a máquina está "Online",
   *   atualiza classes CSS para dot ativo e retorna o horário formatado (HH:mm).
   *
   * - Se a última conexão for superior a 48 horas, indica que a máquina está "Offline",
   *   atualiza classes CSS para dot inativo e retorna uma string relativa ao tempo de inatividade (ex: "há 3 dias").
   *
   * @param curentDate - String com a data/hora da última conexão da máquina.
   * @returns string - Retorna o horário formatado ou a data relativa da última conexão.
   */
  adjustDate(curentDate: string) {
    try {
      // Converte a string para um objeto Date
      const date = new Date(curentDate);
      const now = new Date();
      var daysInactive: string = "";

      // Calcula a diferença de tempo em milissegundos entre agora e a última conexão
      const diffMs = now.getTime() - date.getTime();

      // Define se a diferença está dentro de 48 horas (em ms)
      const isWithin48Hours = diffMs <= 48 * 60 * 60 * 1000;

      if (isWithin48Hours) {
        // Formata a hora e minutos para exibição (formato HH:mm)
        const hours = String(date.getHours()).padStart(2, "0");
        const minutes = String(date.getMinutes()).padStart(2, "0");
        const time = `${hours}:${minutes}`;

        // Configura as classes CSS e textos para status online (ativo)
        this.statusDot = "status-dot-active";
        this.statusDotTitle = "Online";
        this.statusDotPing = "status-ping-active";
        this.statusDotContent =
          "Esse Status representa que o equipamento está ou ficou online nas últimas 48 Horas";

        return time;
      } else {
        // Obtém uma string com o tempo relativo da última conexão (ex: "há 3 dias")
        daysInactive = this.getRelativeDateString(curentDate);

        // Configura as classes CSS e textos para status offline (inativo)
        this.statusDot = "status-dot-inactive";
        this.statusDotPing = "status-ping-inactive";
        this.statusDotTitle = "Offline";
        this.statusDotContent =
          "Esse Status representa que o equipamento não ficou online nas ultimas 48 horas.";

        return daysInactive;
      }
    } catch (err: string | any) {
      this.errorType = "Erro de Conversão";
      this.messageError = "Erro ao converter data/dot: " + err;
      this.showMessage = true;
      return console.error(err);
    }
  }

  /**
   * Função que retorna uma string representando a diferença de tempo relativa entre a data informada e o momento atual.
   *
   * Utiliza a biblioteca dayjs com o plugin relativeTime para gerar expressões como "há um mês", "há 3 dias", etc.
   * Define o idioma para português brasileiro ("pt-br") para formatação correta da string.
   *
   * Caso ocorra algum erro na conversão ou processamento da data, captura a exceção,
   * define mensagens de erro no componente e registra o erro no console.
   *
   * @param dateString - String representando a data a ser convertida para formato relativo.
   * @returns string - String formatada com o tempo relativo, ou mensagem de erro caso ocorra exceção.
   */
  getRelativeDateString(dateString: string): string {
    try {
      // Extende o dayjs com o plugin de tempo relativo
      dayjs.extend(relativeTime);
      // Define a localidade para português brasileiro
      dayjs.locale("pt-br");
      // Retorna a diferença relativa da data para o presente (ex: "há um mês")
      return dayjs(dateString).fromNow();
    } catch (err: string | any) {
      // Em caso de erro, define o tipo e mensagem de erro, e registra no console
      this.errorType = "Erro de Conversão";
      this.messageError = "Erro ao converter a data: " + err;
      console.error(err);
      this.showMessage = true;
      return err;
    }
  }

  /**
   * Função que manipula o evento de clique com o botão direito do mouse para exibir um menu contextual customizado.
   *
   * Ao ser acionada, previne a abertura do menu padrão do navegador,
   * configura a visibilidade do menu personalizado como verdadeira,
   * define a posição do menu na tela com base na posição do cursor do mouse,
   * e ativa uma flag para indicar que o menu está aberto de forma singular.
   *
   * Caso ocorra algum erro durante o processo, captura a exceção, define mensagens de erro no componente
   * e registra o erro no console.
   *
   * @param event - Evento do mouse que dispara a função (MouseEvent).
   */
  onRightClick(event: MouseEvent) {
    try {
      // Previne o menu de contexto padrão do navegador
      event.preventDefault();

      // Exibe o menu customizado
      this.menuVisible = true;

      // converte o evento de click para obter o elemento HTML
      const target = event.target as HTMLElement;
      // seta qual equipamento foi selecionado para abrir o menu custom
      this.menuCustomSelect = target;
      const parent = this.menuCustomSelect?.parentElement;
      if (parent) {
        const tds = parent.getElementsByTagName("td");
        const secondTd = tds[1]; // índice 1 = segundo <td>
        this.machineName = secondTd.innerText;
      }

      // Define a posição do menu com base na posição do cursor no evento
      this.menuPosition = {
        x: event.clientX,
        y: event.clientY,
      };

      // Marca que o menu está aberto individualmente (sem múltiplos menus simultâneos)
      this.menuSingular = true;
    } catch (err) {
      // Tratamento de erros inesperados com mensagem e log no console
      this.errorType = "Erro Inesperado";
      this.messageError = "Erro ao abrir menu de configuração Rápida: " + err;
      this.showMessage = true;
      console.error(err);
    }
  }

  /**
   * Função responsável por ocultar o menu customizado de contexto.
   *
   * Define as flags que controlam a visibilidade do menu e sua condição singular para falso,
   * garantindo que o menu deixe de ser exibido e que o estado do componente seja atualizado corretamente.
   */
  hideMenu() {
    this.menuVisible = false;
    this.menuSingular = false;
  }

  forceConection() {
    this.canViewProcessTab = true;
    this.processExec = "Forçar Conexão";
    this.processAnimation = "process-tab-animation-maximize";
    this.hideMenu();
  }

  resizeProcessTab() {
    if (this.expenseProcessTab) {
      this.buttomSize = "/static/assets/images/maximize.png";
      this.processAnimation = "process-tab-animation-minimize";
      this.processHeader = "Processo: " + this.processExec;
      this.statusHeader = "Status: " + this.statusPorcentage;
      this.machineHeader = "Máquina: " + this.machineName;
      const process = document.getElementById("process-tb");
      if (process) {
        process.style.height = "2em";
        this.expenseProcessTab = false;
      }
    } else {
      this.buttomSize = "/static/assets/images/minimize.png";
      this.processAnimation = "process-tab-animation-maximize";
      this.processHeader = "Processo";
      this.statusHeader = "Status";
      this.machineHeader = "Máquina";
      const process = document.getElementById("process-tb");
      if (process) {
        process.style.height = "4em";
        this.expenseProcessTab = true;
      }
    }
  }

  createRange(n: number): number[] {
    return Array.from({ length: n }, (_, i) => i);
  }

  nextPageMachines(index: number) {
    this.canViewLoadingSearch = true;
    this.listMachines = this.machines[index];
    this.canViewLoadingSearch = false;
  }
}
