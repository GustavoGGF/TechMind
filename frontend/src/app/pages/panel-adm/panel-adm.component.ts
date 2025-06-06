import { Component, OnInit } from "@angular/core";
import { UtilitiesModule } from "../../utilities/utilities.module";
import { CommonModule } from "@angular/common";
import { HttpClient, HttpClientModule } from "@angular/common/http";
import "../../../assets/bootstrap-5.3.3-dist/js/bootstrap.js";
import "../../../assets/bootstrap-5.3.3-dist/js/bootstrap.bundle.min.js";
import { catchError, throwError } from "rxjs";
import "dayjs/locale/pt-br";
import { WebSocketService } from "../../services/websocket.service";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import "dayjs/locale/pt-br";

declare var bootstrap: any;
@Component({
  selector: "app-panel-adm",
  standalone: true,
  imports: [UtilitiesModule, CommonModule, HttpClientModule],
  templateUrl: "./panel-adm.component.html",
  styleUrl: "./panel-adm.component.css",
})
export class PanelAdmComponent implements OnInit {
  constructor(private http: HttpClient, private wsService: WebSocketService) {}

  popoverTriggerList: NodeListOf<Element> = [] as any;
  popoverList: any[] = [];

  machines: any;
  listMachines: any;
  totalData: any;
  menuCustomSelect: any;
  name: any;

  buttomSize = "/static/assets/images/minimize.png";
  computers_class: string = "";
  device_class: string = "";
  errorType: string = "";
  home_class: string = "";
  machineHeader: string = "Máquina";
  machineIP: string = "";
  machineName: string = "";
  messageError: string = "";
  panel_class: string = "active";
  processAnimation: string = "";
  processExec: string = "";
  processHeader: string = "Processo";
  statusPorcentage: string = "0%";
  statusHeader: string = "Status";
  detailFailMessage: string = "";
  arrow_up: string = "/static/assets/images/seta2.png";
  arrow_down: string = "/static/assets/images/seta.png";

  canView: boolean = false;
  canViewLoadingSearch: boolean = true;
  canViewProcessTab: boolean = false;
  expenseProcessTab: boolean = true;
  dotActive: boolean = false;
  menuSingular: boolean = false;
  menuVisible: boolean = false;
  showMessage: boolean = false;
  failServerComunication: boolean = false;

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

      this.wsService.getMessages().subscribe((msg: any) => {
        const codeInterface = msg.code;
        const statusInterface = msg.status;

        switch (codeInterface) {
          default:
            break;
          case "cnmh":
            this.statusPorcentage = "30%";
            this.detailFailMessage =
              "Problema Com conexão com a maquina " +
              this.machineName +
              ": verificar se TechMind está rodando";
            break;
          case "vldvr":
            this.statusPorcentage = "60%";
            this.detailFailMessage =
              "Não foi possivel verificar a versão do TechMind da máquina: " +
              this.machineName +
              ", atualizar manualmente.";
            break;
          case "crtvs":
            this.statusPorcentage = "100%";
            this.detailFailMessage = "";
            break;
        }

        if (statusInterface !== 200) {
          this.statusPorcentage = "Fail";
          return;
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
          this.machines = this.groupSplitter(data.machines, 100);
          this.listMachines = this.machines[0];

          this.ordemList("increasing", "insertion_date");

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
   * Atualiza o status visual de uma máquina com base na data da última atividade.
   *
   * A função verifica se a data fornecida (`curentDate`) está dentro do intervalo das últimas 48 horas.
   * Com base nisso, aplica estilos visuais (CSS) nos elementos informados para indicar se o equipamento
   * está "Online" (ativo recentemente) ou "Offline" (inativo há mais de 48 horas).
   *
   * - Se estiver dentro de 48 horas, retorna o horário formatado da última atividade (HH:mm).
   * - Caso contrário, retorna uma string relativa com o tempo de inatividade (ex: "há 3 dias").
   *
   * Os elementos DOM são identificados pelos seus respectivos IDs (`elementId` e `secondElementId`) e têm
   * suas classes CSS e atributos de tooltip atualizados dinamicamente.
   *
   * @param curentDate - String representando a data/hora da última conexão da máquina.
   * @param elementId - ID do elemento principal onde será aplicado o status visual (dot).
   * @param secondElementId - ID do segundo elemento relacionado (por exemplo, indicador de ping).
   * @returns string | void - O horário formatado ou a data relativa da última atividade.
   */
  adjustDate(
    curentDate: string,
    elementId: string,
    secondElementId: string
  ): string | void {
    try {
      // Converte a string para um objeto Date
      const date = new Date(curentDate);
      const now = new Date();

      // Calcula a diferença de tempo em milissegundos entre agora e a última conexão
      const diffMs = now.getTime() - date.getTime();

      // Define se a diferença está dentro de 48 horas (em ms)
      const isWithin48Hours = diffMs <= 48 * 60 * 60 * 1000;

      var element = document.getElementById(elementId);
      var secondElement = document.getElementById(secondElementId);

      if (isWithin48Hours) {
        // Formata a hora e minutos para exibição (formato HH:mm)
        const formattedTime = this.getRelativeDateString(curentDate);

        this.setStatus(
          element,
          secondElement,
          "status-dot-active",
          "status-ping-active",
          "Online",
          "Esse Status representa que o equipamento está ou ficou online nas últimas 48 Horas."
        );

        return formattedTime;
      } else {
        // Obtém uma string com o tempo relativo da última conexão (ex: "há 3 dias")
        const daysInactive = this.getRelativeDateString(curentDate);

        this.setStatus(
          element,
          secondElement,
          "status-dot-inactive",
          "status-ping-inactive",
          "Offline",
          "Esse Status representa que o equipamento não ficou online nas últimas 48 horas."
        );

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
   * Função auxiliar que aplica classes CSS e atributos de dados (data attributes) em elementos HTML,
   * com o objetivo de indicar visualmente um determinado status e exibir informações adicionais por tooltip (ou popover).
   *
   * Essa função é utilizada para alterar dinamicamente o estado visual de elementos na interface,
   * como por exemplo, sinalizar que um usuário está online, offline, inativo, etc.
   *
   * Aplica a classe CSS de status no primeiro elemento, junto com os atributos "data-bs-title" e "data-bs-content",
   * utilizados por bibliotecas como Bootstrap para exibição de tooltips ou popovers.
   * Também aplica uma classe CSS de animação (como "ping" ou similar) em um segundo elemento, se fornecido.
   *
   * O uso do modificador `private` indica que essa função é interna à classe e não deve ser acessada externamente.
   *
   * @param element - O elemento HTML principal que receberá a classe de status e os atributos de dados.
   * @param secondElement - Um segundo elemento HTML opcional que receberá a classe de animação (ping).
   * @param statusClass - Classe CSS a ser aplicada no primeiro elemento para representar o status visual.
   * @param pingClass - Classe CSS a ser aplicada no segundo elemento para animações visuais (ex: efeito de "ping").
   * @param title - Texto que será definido no atributo `data-bs-title`, utilizado para exibição de tooltip.
   * @param content - Texto que será definido no atributo `data-bs-content`, utilizado para exibição de tooltip.
   */
  private setStatus(
    element: HTMLElement | null,
    secondElement: HTMLElement | null,
    statusClass: string,
    pingClass: string,
    title: string,
    content: string
  ): void {
    // Verifica se o primeiro elemento foi fornecido
    if (element) {
      // Aplica a classe de status (ex: 'text-success') ao elemento
      element.classList.add(statusClass);

      // Define os atributos de dados utilizados por tooltips/popovers do Bootstrap
      element.setAttribute("data-bs-title", title);
      element.setAttribute("data-bs-content", content);
    }

    // Verifica se o segundo elemento foi fornecido
    if (secondElement) {
      // Aplica a classe de animação (ex: 'ping') ao segundo elemento
      secondElement.classList.add(pingClass);
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
      dayjs.extend(relativeTime);
      dayjs.locale("pt-br");

      return dayjs(dateString).fromNow();
    } catch (err: any) {
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

  /**
   * Força a atualização de uma máquina.
   *
   * A função extrai o nome e o IP da máquina a partir dos elementos `<td>` dentro do componente
   * customizado de menu (`menuCustomSelect`). Em seguida:
   *
   * - Atualiza os estados da interface para exibir a aba de processo.
   * - Define o nome do processo e a animação correspondente.
   * - Oculta o menu de contexto.
   * - Inicia o processo de conexão com o comando `"force-update"`.
   *
   * @remarks
   * Essa função depende da estrutura do DOM conter ao menos três células `<td>` na linha associada ao menu.
   *
   * @returns void
   */
  forceUpdate() {
    const parent = this.menuCustomSelect?.parentElement;
    if (parent) {
      const tds = parent.getElementsByTagName("td");
      const secondTd = tds[1];
      const tirdTd = tds[2];

      this.machineName = secondTd.innerText;
      this.machineIP = tirdTd.innerText;
    }
    this.canViewProcessTab = true;
    this.processExec = "Atualizar Versão";
    this.processAnimation = "process-tab-animation-maximize";
    this.hideMenu();
    this.contactMachine("force-update");
  }

  /**
   * Alterna entre os modos expandido e minimizado da aba de processo.
   *
   * Quando expandida, a aba exibe informações completas sobre o processo em execução,
   * como nome, status e máquina. Quando minimizada, reduz a altura do componente visual
   * e exibe um resumo mais compacto das informações.
   *
   * Também atualiza dinamicamente:
   * - Ícone de botão (maximizar/minimizar),
   * - Classe de animação (para transição suave),
   * - Títulos exibidos na interface.
   *
   * A manipulação é feita diretamente no elemento DOM com ID `process-tb`.
   *
   * @returns void
   */

  resizeProcessTab(): void {
    const processElement = document.getElementById("process-tb");

    if (!processElement) return;

    const isExpanded = this.expenseProcessTab;

    // Atualiza propriedades visuais
    this.buttomSize = isExpanded
      ? "/static/assets/images/maximize.png"
      : "/static/assets/images/minimize.png";

    this.processAnimation = isExpanded
      ? "process-tab-animation-minimize"
      : "process-tab-animation-maximize";

    this.processHeader = isExpanded
      ? "Processo: " + this.processExec
      : "Processo";

    this.statusHeader = isExpanded
      ? "Status: " + this.statusPorcentage
      : "Status";

    this.machineHeader = isExpanded
      ? "Máquina: " + this.machineName
      : "Máquina";

    // Ajusta altura visual da aba
    processElement.style.height = isExpanded ? "2em" : "4em";

    // Alterna estado de expansão
    this.expenseProcessTab = !isExpanded;
  }

  /**
   * Gera um array numérico de comprimento `n`, contendo valores inteiros sequenciais a partir de 0.
   *
   * Exemplo:
   * ```ts
   * createRange(5); // [0, 1, 2, 3, 4]
   * ```
   *
   * Pode ser utilizado, por exemplo, para criar iterações em templates ou gerar índices de forma dinâmica.
   *
   * @param n - Número de elementos no array gerado.
   * @returns number[] - Array contendo números de 0 até n-1.
   */
  createRange(n: number): number[] {
    return Array.from({ length: n }, (_, i) => i);
  }

  /**
   * Atualiza a lista visível de máquinas com base na página selecionada.
   *
   * Essa função é geralmente utilizada junto à função `createRange`, onde `index`
   * representa o índice da página atual dentro de uma estrutura paginada (`machines`).
   *
   * a lista de máquinas (`listMachines`) com os dados correspondentes à página selecionada.
   *
   * @param index - Índice da página de máquinas a ser exibida.
   */
  nextPageMachines(index: number) {
    this.canViewLoadingSearch = true;
    this.listMachines = this.machines[index];
    this.canViewLoadingSearch = false;
  }

  /**
   * Envia uma requisição GET para a API a fim de executar uma ação específica em uma máquina remota.
   *
   * A ação a ser realizada é definida pelo parâmetro `event` (ex: "force-update"),
   * e a requisição é direcionada para o IP da máquina armazenado em `this.machineIP`.
   *
   * Durante a execução, o progresso é inicialmente setado para "25%" e, em caso de erro,
   * o status HTTP é armazenado e o erro é logado no console.
   *
   * @param event - Identificador da ação a ser solicitada na máquina (ex: "force-update").
   */
  contactMachine(event: string) {
    this.statusPorcentage = "25%";
    this.http
      .get(
        "/home/panel-adm/contact-machine/" + event + "/" + this.machineIP,
        {}
      )
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe({
        next: (response) => {
          // Lógica para tratar resposta, se necessário
          return response;
        },
        error: (err) => {
          console.error("Erro na requisição", err);
        },
      });
  }

  /**
   * Ordena a lista de máquinas com base na coluna e direção especificadas.
   *
   * A direção da ordenação é determinada pelo parâmetro `order`, que pode ser "increasing"
   * (ordem crescente) ou "decreasing" (ordem decrescente). A função mapeia as colunas para
   * funções específicas de ordenação correspondentes e executa a ordenação apropriada
   * na lista `listMachines`.
   *
   * - Para colunas de data, nome, IP, usuário logado ou versão, chama a função de ordenação
   *   correspondente, passando o sinal de direção (`+` para crescente e `-` para decrescente).
   * - Caso a coluna não seja reconhecida, exibe um aviso no console.
   *
   * @param order - Direção da ordenação: "increasing" para crescente, "decreasing" para decrescente.
   * @param column - Nome da coluna pela qual a lista deve ser ordenada (ex: "name", "ip").
   */
  ordemList(order: "increasing" | "decreasing", column: string) {
    // Define o sinal de direção da ordenação baseado no parâmetro order
    // "+" indica crescente, "-" indica decrescente
    const direction = order === "increasing" ? "+" : "-";

    // Mapeamento das colunas para as funções de ordenação correspondentes
    // Cada função chama o método correto passando a lista e o sentido da ordenação
    const sortActions: { [key: string]: () => void } = {
      insertion_date: () =>
        this.sortMachinesByDate(this.listMachines, direction),
      name: () => this.sortMachinesByName(this.listMachines, direction, "name"),
      ip: () => this.sortMachinesByIp(this.listMachines, direction),
      logged_user: () =>
        this.sortMachinesByName(this.listMachines, direction, "logged_user"),
      version: () => this.sortMachinesByVersion(this.listMachines, direction),
    };

    // Recupera a função de ordenação correspondente à coluna solicitada
    const sortFn = sortActions[column];

    if (sortFn) {
      // Se a coluna for válida, executa a função de ordenação
      sortFn();
    } else {
      // Caso a coluna seja desconhecida, exibe um aviso no console para facilitar debug
      console.warn(`Coluna desconhecida: ${column}`);
    }
  }

  /**
   * Ordena uma lista de máquinas pela data de inserção.
   *
   * A ordenação pode ser crescente ou decrescente com base no parâmetro `signal`.
   * Caso alguma máquina não tenha a data (`insertion_date`) definida, ela não afeta a ordenação.
   *
   * @param machines - Array de objetos representando as máquinas, cada um com a propriedade `insertion_date`.
   * @param signal - Direção da ordenação: "+" para decrescente (mais recentes primeiro), "-" para crescente (mais antigos primeiro).
   * @returns Array ordenado de máquinas conforme a data.
   */
  sortMachinesByDate(machines: any[], signal: "+" | "-"): any[] {
    return machines.sort((a, b) => {
      // Se algum dos objetos não possuir data de inserção, não altera a ordem entre eles
      if (!a.insertion_date || !b.insertion_date) return 0;

      // Converte as datas de inserção para timestamps (milissegundos desde 1970)
      const dateA = new Date(a.insertion_date).getTime();
      const dateB = new Date(b.insertion_date).getTime();

      // Ordena conforme o sinal:
      // "+" (decrescente): máquinas com datas mais recentes vêm primeiro (dateB - dateA)
      // "-" (crescente): máquinas com datas mais antigas vêm primeiro (dateA - dateB)
      return signal === "+" ? dateB - dateA : dateA - dateB;
    });
  }

  /**
   * Ordena uma lista de máquinas alfabeticamente com base em uma coluna especificada.
   *
   * A ordenação pode ser crescente ou decrescente, definida pelo parâmetro `order`.
   * A comparação é feita de forma case-insensitive para garantir consistência.
   * Caso algum item não possua valor para a coluna especificada, a ordenação mantém a ordem original desses itens.
   *
   * @param machines - Array de objetos representando as máquinas.
   * @param order - Direção da ordenação: "+" para crescente (A-Z), "-" para decrescente (Z-A).
   * @param column - Nome da propriedade do objeto pelo qual será feita a ordenação.
   * @returns Array ordenado de máquinas conforme a coluna e ordem especificadas.
   */
  sortMachinesByName(machines: any[], order: "+" | "-", column: string): any[] {
    return machines.sort((a, b) => {
      // Se algum dos objetos não possuir valor para a coluna, mantém a ordem entre eles
      if (!a[column] || !b[column]) return 0;

      // Normaliza para string em minúsculo para comparação case-insensitive
      const nameA = a[column]?.toString().toLowerCase();
      const nameB = b[column]?.toString().toLowerCase();

      // Compara os valores para ordenar
      if (nameA < nameB) return order === "+" ? -1 : 1; // Crescente ou decrescente
      if (nameA > nameB) return order === "+" ? 1 : -1;
      return 0; // valores iguais
    });
  }

  /**
   * Ordena uma lista de máquinas com base no endereço IP.
   *
   * A ordenação considera cada octeto do IP como um número para garantir a ordem correta
   * (ex: "192.168.1.10" vem depois de "192.168.1.2").
   * A direção da ordenação pode ser crescente ("+") ou decrescente ("-").
   * Caso algum IP seja nulo ou indefinido, esses itens mantêm a posição relativa (retorna 0).
   *
   * @param machines - Array de objetos representando as máquinas.
   * @param order - Direção da ordenação: "+" para crescente, "-" para decrescente.
   * @returns Array ordenado de máquinas conforme o endereço IP e direção especificada.
   */
  sortMachinesByIp(machines: any[], order: "+" | "-"): any[] {
    return machines.sort((a, b) => {
      // Ignora itens sem IP válido para manter a ordem relativa deles
      if (!a.ip || !b.ip) return 0;

      // Converte IP para array de números (octetos)
      const ipA = a.ip.split(".").map(Number);
      const ipB = b.ip.split(".").map(Number);

      // Compara cada octeto na ordem, até encontrar diferença
      for (let i = 0; i < 4; i++) {
        if (ipA[i] !== ipB[i]) {
          // Retorna a diferença baseada na direção da ordenação
          return order === "+" ? ipA[i] - ipB[i] : ipB[i] - ipA[i];
        }
      }

      // IPs iguais mantêm a ordem
      return 0;
    });
  }

  /**
   * Ordena uma lista de máquinas com base na versão do software.
   *
   * A versão é esperada no formato "x.y.z" (ex: "1.2.10").
   * A função compara cada parte numérica da versão (major, minor, patch) para definir a ordem correta.
   * Se a versão estiver ausente, assume "0.0.0" para garantir a comparação.
   * A ordenação pode ser crescente ("+") ou decrescente ("-").
   *
   * @param machines - Array de objetos representando as máquinas.
   * @param order - Direção da ordenação: "+" para crescente, "-" para decrescente.
   * @returns Array ordenado de máquinas conforme a versão e direção especificada.
   */
  sortMachinesByVersion(machines: any[], order: "+" | "-"): any[] {
    return machines.sort((a, b) => {
      // Se uma das versões for indefinida, mantém a ordem relativa
      if (!a.version || !b.version) return 0;

      // Divide as versões em arrays numéricos [major, minor, patch]
      const v1 = (a.version || "0.0.0").split(".").map(Number);
      const v2 = (b.version || "0.0.0").split(".").map(Number);

      // Compara cada parte da versão sequencialmente
      for (let i = 0; i < 3; i++) {
        if (v1[i] !== v2[i]) {
          // Retorna a diferença com base na ordem desejada
          return order === "+" ? v1[i] - v2[i] : v2[i] - v1[i];
        }
      }

      // Versões iguais mantêm a ordem
      return 0;
    });
  }
}
