import {
  Component,
  ElementRef,
  OnInit,
  Renderer2,
  ViewChild,
} from "@angular/core";
import { UtilitiesModule } from "../../utilities/utilities.module";
import { CommonModule } from "@angular/common";
import { FormsModule } from "@angular/forms";
import {
  HttpClient,
  HttpClientModule,
  HttpHeaders,
} from "@angular/common/http";
import { catchError, throwError } from "rxjs";
import { saveAs } from "file-saver";
// import { ChangeDetectorRef } from '@angular/core';
// import { saveAs } from 'file-saver';

interface Software {
  name: string;
  ids: string[];
  // Adicione outras propriedades se necessário
}

export interface ReportResponse {
  file_name: string;
  file_content: string;
}

@Component({
  selector: "app-computers",
  standalone: true,
  imports: [UtilitiesModule, CommonModule, HttpClientModule, FormsModule],
  templateUrl: "./computers.component.html",
  styleUrl: "./computers.component.css",
})
export class ComputersComponent implements OnInit {
  constructor(
    private http: HttpClient,
    private renderer: Renderer2,
    private el: ElementRef // private cdRef: ChangeDetectorRef
  ) {}
  // Declarando variaveis any
  @ViewChild("selectElement") selectElement!: ElementRef<HTMLSelectElement>;

  dataMachines: any;
  name: any;
  status: any;
  token: any;

  // Declarando variaveis string
  all_quantity: string = "";
  arrow_up: string = "/static/assets/images/seta2.png";
  arrow_down: string = "/static/assets/images/seta.png";
  computers_class: string = "active";
  device_class: string = "";
  errorType: string = "";
  fifty_quantity: string = "";
  input_name: string = "";
  input_pwd: string = "";
  input_username: string = "";
  inputUsername: string = "";
  inputPass: string = "";
  home_class: string = "";
  messageError: string = "";
  one_hundred_quantity: string = "";
  ten_quantity: string = "";
  reset_filter: string = "/static/assets/images/filtro.png";
  quantity_filter: string | null = "";
  selectedValue: string = "None";
  selectedReports: string = "None";
  closeBTN: string = "/static/assets/images/fechar.png";

  // Declarando variaveis boolean
  canView: boolean = false;
  canViewCredentials: boolean = false;
  canViewCredentialsLoading: boolean = false;
  canViewMachines: boolean = false;
  canViewMessage: boolean = false;
  checkedAll: boolean = true;
  showMessage: boolean = false;

  // Declarando variaveis list
  dis_list: string[] = [];
  so_list: string[] = [];
  soft_list: Software[] = [];

  softwares_list: any;

  // Função iniciada ao carregar a pagina
  ngOnInit(): void {
    // Pegando valores do usuario
    this.name = localStorage.getItem("name");

    // Pegando valor de quantitade do filtro
    this.quantity_filter = localStorage.getItem("quantity");
    if (this.quantity_filter == null) {
      localStorage.setItem("quantity", "10");
      this.quantity_filter = "10";
    }

    // Verificando se os dados existem
    if (this.name.length == 0 || this.name == null) {
      this.errorType = "Falta de Dados";
      this.messageError =
        "Ouve um erro ao acessar dados do LDAP, contatar a TI";
      this.showMessage = true;
    } else {
      this.canView = true;
      this.getToken();
      this.getData();
      this.getSO();
      this.getDistribution();
    }
  }

  closeMessage() {
    this.canViewMessage = false;
  }

  // Função que obtem o token CSRF
  getToken(): void {
    console.log("pegou token");

    this.http
      .get("/home/get-token", {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          if (this.status === 0) {
          }

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.token = data.token;
        }
      });
  }

  // Buscando as maquinas disponiveis
  getData(): void {
    this.http
      .get("/home/computers/get-data/" + this.quantity_filter, {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;
          // Ajustando filtro de quantidade
          switch (this.quantity_filter) {
            default:
              break;
            case "10":
              this.ten_quantity = "active_filter";
              this.fifty_quantity = "";
              this.one_hundred_quantity = "";
              this.all_quantity = "";
              break;
            case "50":
              this.ten_quantity = "";
              this.fifty_quantity = "active_filter";
              this.one_hundred_quantity = "";
              this.all_quantity = "";
              break;
            case "100":
              this.ten_quantity = "";
              this.fifty_quantity = "";
              this.one_hundred_quantity = "active_filter";
              this.all_quantity = "";
              break;
            case "all":
              this.ten_quantity = "";
              this.fifty_quantity = "";
              this.one_hundred_quantity = "";
              this.all_quantity = "active_filter";
          }

          this.canViewMachines = true;
          // Apos pegar os dados principais chama a função para preencher o filtro de SO
          return this.mountSoftwares();
        }
      });
  }

  mountSoftwares(): void {
    const machineNames = this.dataMachines.map((machine: any[]) => machine[40]);

    const so = this.dataMachines.map((machine: any[]) => machine[3]);
    const ids = this.dataMachines.map((machine: any[]) => machine[0]);

    // Encontrar todos os índices onde o valor de `so` é "Microsoft Windows 10 Pro"
    const windows10ProIndexes = so
      .map((value: string, index: number) =>
        value === "Microsoft Windows 10 Pro" ? index : -1
      )
      .filter((index: number) => index !== -1);

    // Usar esses índices para pegar os valores correspondentes de `machineNames` e `ids`
    const result = windows10ProIndexes.map((index: number) => ({
      name: machineNames[index],
      id: ids[index],
    }));

    // Processar cada resultado com a função `stringToSortedArray`
    result.forEach(({ name, id }: { name: string; id: string }) => {
      const names = this.stringToSortedArray(name);
      return this.updateSoftwareList(names, id);
    });
  }

  stringToSortedArray(array: string): string[] {
    if (!array || array.trim().length <= 1) {
      // Se a string estiver vazia ou muito curta, retornar o erro
      console.error("Formato de string inválido.");
      return [];
    }
    try {
      let soft_list: any[] = [];
      const trimmedData = array.trim();
      // Verifica se a string é um array válido
      if (trimmedData.startsWith("[") && trimmedData.endsWith("]")) {
        soft_list = JSON.parse(trimmedData.replace(/'/g, '"'));
      }

      // Verifica se a string é um objeto único
      else if (trimmedData.startsWith("{") && trimmedData.endsWith("}")) {
        // Adiciona colchetes para transformar em um array com um único item
        const arrayString = `[${trimmedData}]`;
        soft_list = JSON.parse(arrayString.replace(/'/g, '"'));
      }
      // Verifica se soft_list é um array de objetos
      if (
        Array.isArray(soft_list) &&
        soft_list.every((item) => typeof item === "object" && item !== null)
      ) {
        // Extrai o valor da propriedade `name` de cada objeto
        return soft_list.map((item) => item.name);
      } else {
        console.error("A string não contém um array de objetos válidos.");
        return [];
      }
    } catch (error) {
      console.error("Erro ao converter a string para JSON:", error);
      return [];
    }
  }

  updateSoftwareList(names: string[], id: string): void {
    names.forEach((name: string) => {
      // Verifica se já existe um objeto com o mesmo nome em soft_list
      const software = this.soft_list.find(
        (software) => software.name === name
      );

      if (software) {
        // Adiciona o ID ao array de IDs se ainda não estiver presente
        if (!software.ids.includes(id)) {
          software.ids.push(id);
        }
      } else {
        // Adiciona um novo objeto com o nome e o ID
        this.soft_list.push({ name, ids: [id] });
        this.soft_list.sort((a, b) =>
          a.name.toLowerCase().localeCompare(b.name.toLowerCase())
        );
      }
    });
  }

  generateMachineReport(index: number): void {
    // Obtém o software selecionado usando o índice
    const selected_Soft = this.soft_list[index];

    // Filtra o array machines para encontrar máquinas cujo ID corresponde ao ID do software selecionado
    const filteredMachines = this.dataMachines.filter(
      (machine: string[]) => selected_Soft.ids.includes(machine[0]) // Verifica se o ID da máquina está na lista de IDs do software
    );

    // Exibe o relatório ou faz algo com as máquinas filtradas
    this.dataMachines = filteredMachines;
  }

  resetSoft(): void {
    this.selectedValue = "None";

    this.http
      .get("/home/computers/get-data/" + this.quantity_filter, {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;
          // Ajustando filtro de quantidade
          switch (this.quantity_filter) {
            default:
              break;
            case "10":
              this.ten_quantity = "active_filter";
              this.fifty_quantity = "";
              this.one_hundred_quantity = "";
              this.all_quantity = "";
              break;
            case "50":
              this.ten_quantity = "";
              this.fifty_quantity = "active_filter";
              this.one_hundred_quantity = "";
              this.all_quantity = "";
              break;
            case "100":
              this.ten_quantity = "";
              this.fifty_quantity = "";
              this.one_hundred_quantity = "active_filter";
              this.all_quantity = "";
              break;
            case "all":
              this.ten_quantity = "";
              this.fifty_quantity = "";
              this.one_hundred_quantity = "";
              this.all_quantity = "active_filter";
          }
          this.canViewMachines = true;
        }
      });
  }

  // Função para pegar os valores do filto de SO
  getDistribution(): void {
    this.http
      .get("/home/computers/get-data-DIS", {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dis_list = data.DIS;
        }
      });
  }

  // Função para pegar os valores do filto de Distribuição
  getSO(): string[] | any {
    this.http
      .get("/home/computers/get-data-SO", {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.so_list = data.SO;
          // Após pegar os dados do filtro chama a função para pegar os dados do filtro de Distribuição
        }
      });
  }

  // Função para redirecionar para a pagina de visualização da maquina
  onRowClick(index: number) {
    const selectedMachine = this.dataMachines[index];
    let mac = selectedMachine[0];
    mac = mac.replace(/:/g, "-");

    return window.open("/home/computers/view-machine/" + mac, "_blank");
  }

  // Função obter o valor do SO que deseja filtrar
  onRowClickSO(index: number) {
    this.canViewMachines = false;

    this.dataMachines = null;

    let so;

    if (index == 69) {
      so = "all";
    } else {
      const selectedSO = this.so_list[index];

      so = selectedSO[0];
    }

    this.http
      .get(
        "/home/computers/get-data-SO-filter/" + this.quantity_filter + "/" + so,
        {}
      )
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  // Função obter o valor de distribution que deseja filtrar
  onRowClickDIS(index: number): void {
    this.canViewMachines = false;

    this.dataMachines = null;

    let so;

    if (index == 69) {
      so = "all";
    } else {
      const selectedDis = this.dis_list[index];

      so = selectedDis[0];
    }

    this.http
      .get(
        "/home/computers/get-data-DIS-filter/" +
          this.quantity_filter +
          "/" +
          so,
        {}
      )
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  // Seta quantidade de maquinas a serem exibidas para 10
  getTen(): void {
    localStorage.setItem("quantity", "10");

    this.quantity_filter = "10";

    this.ten_quantity = "active_filter";
    this.fifty_quantity = "";
    this.one_hundred_quantity = "";
    this.all_quantity = "";

    this.canViewMachines = false;
    this.dataMachines = null;

    this.http
      .get("/home/computers/get-quantity/10", {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  // Seta quantidade de maquinas a serem exibidas para 50
  getFifty(): void {
    localStorage.setItem("quantity", "50");

    this.quantity_filter = "50";
    this.ten_quantity = "";
    this.fifty_quantity = "active_filter";
    this.one_hundred_quantity = "";
    this.all_quantity = "";

    this.canViewMachines = false;
    this.dataMachines = null;

    this.http
      .get("/home/computers/get-quantity/50", {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  // Seta a quantidade de maquinas a serem exibidas para 100
  getOneHundred(): void {
    localStorage.setItem("quantity", "100");
    this.quantity_filter = "100";
    this.ten_quantity = "";
    this.fifty_quantity = "";
    this.one_hundred_quantity = "active_filter";
    this.all_quantity = "";

    this.canViewMachines = false;
    this.dataMachines = null;

    this.http
      .get("/home/computers/get-quantity/100", {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  // Seta a quantidade de maquinas a serem exibidas para todas
  getAll(): void {
    localStorage.setItem("quantity", "all");
    this.quantity_filter = "all";
    this.ten_quantity = "";
    this.fifty_quantity = "";
    this.one_hundred_quantity = "";
    this.all_quantity = "active_filter";

    this.canViewMachines = false;
    this.dataMachines = null;

    this.http
      .get("/home/computers/get-quantity/all", {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  // Função que formata as datas que aparecem na tabela dos computadores
  formatDate(date: string): string {
    const parsedDate = new Date(date);
    const day = String(parsedDate.getDate()).padStart(2, "0");
    const month = String(parsedDate.getMonth() + 1).padStart(2, "0"); // Meses são baseados em 0 (Janeiro é 0)
    const year = parsedDate.getFullYear();

    return `${day}/${month}/${year}`;
  }

  // Reorganiza os os computadores pelo nome em ordem alfabetica
  sortByName(): void {
    this.dataMachines.sort((a: any, b: any) => {
      const nameA = a[1].toUpperCase(); // Ignore case
      const nameB = b[1].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }
      return 0;
    });
  }

  // Reorganiza os os computadores pelo nome em ordem alfabetica invertido
  sortDataByNameDescending(): void {
    this.dataMachines.sort((a: any, b: any) => {
      const nameA = a[1].toUpperCase(); // Ignore case
      const nameB = b[1].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return 1;
      }
      if (nameA > nameB) {
        return -1;
      }
      return 0;
    });
  }

  // Reorganiza os os computadores pelo SO em ordem alfabetica
  sortByNameSO(): void {
    this.dataMachines.sort((a: any, b: any) => {
      const nameA = a[2].toUpperCase(); // Ignore case
      const nameB = b[2].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }
      return 0;
    });
  }

  // Reorganiza os os computadores pelo SO em ordem alfabetica invertido
  sortDataByNameDescendingSO(): void {
    this.dataMachines.sort((a: any, b: any) => {
      const nameA = a[2].toUpperCase(); // Ignore case
      const nameB = b[2].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return 1;
      }
      if (nameA > nameB) {
        return -1;
      }
      return 0;
    });
  }

  // Reorganiza os os computadores pela Distribuição em ordem alfabetica
  sortByNameDis(): void {
    this.dataMachines.sort((a: any, b: any) => {
      const nameA = a[3].toUpperCase(); // Ignore case
      const nameB = b[3].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }
      return 0;
    });
  }

  // Reorganiza os os computadores pela Distribuição em ordem alfabetica invertido
  sortDataByNameDescendingDis(): void {
    this.dataMachines.sort((a: any, b: any) => {
      const nameA = a[3].toUpperCase(); // Ignore case
      const nameB = b[3].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return 1;
      }
      if (nameA > nameB) {
        return -1;
      }
      return 0;
    });
  }

  // Reorganiza os os computadores pela data em ordem decrecente
  sortByDateDesc(): void {
    this.dataMachines.sort((a: any, b: any) => {
      const dateA = new Date(a[4]);
      const dateB = new Date(b[4]);

      return dateB.getTime() - dateA.getTime(); // Mais recente para o mais antigo
    });
  }

  // Reorganiza os os computadores pela data em ordem crescente
  sortByDateASC(): void {
    this.dataMachines.sort((a: any, b: any) => {
      const dateA = new Date(a[4]);
      const dateB = new Date(b[4]);

      return dateA.getTime() - dateB.getTime(); // Mais antigo para o mais novo
    });
  }

  // Obtenendo o nome do dispositivo que o usuario deseja buscar
  getName(event: any): void {
    this.input_name = event.target.value;

    this.canViewMachines = false;

    this.dataMachines = null;

    if (this.input_name.length === 0) {
      return this.getData();
    }

    this.http
      .get(
        "/home/computers/get-machine-varchar/" +
          this.quantity_filter +
          "/" +
          this.input_name,
        {}
      )
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  // Pega o usuario inserido
  getUserName(event: any): void {
    this.input_username = event.target.value;
  }

  // Pega a senha inserida
  getPwd(event: any): void {
    this.input_pwd = event.target.value;
  }

  onReportChange(value: string): void {
    if (value === "DNS") {
      this.reportDNS();
    } else if (value === "reportxls") {
      this.exportMachineReport();
    }
  }
  // Habilita a tela de credencial
  reportDNS(): void {
    this.canViewCredentials = true;
  }

  closeCredential(): void {
    this.canViewCredentials = false;
    this.selectedReports = "None";
    this.inputUsername = "";
    this.inputPass = "";
    this.canViewCredentialsLoading = false;
  }

  exportMachineReport(): void {
    // Seleciona todos os checkboxes com a classe 'ckip'
    const checkboxes = this.el.nativeElement.querySelectorAll(".ckip");

    // Array para armazenar os valores dos checkboxes marcados
    const selectedValues: string[] = [];

    // Itera sobre os checkboxes
    checkboxes.forEach((checkbox: HTMLInputElement) => {
      if (checkbox.checked) {
        // Se o checkbox está marcado, adiciona seu valor ao array
        selectedValues.push(checkbox.value);
      }
    });

    if (selectedValues.length == 0) {
      const selectT = document.getElementById(
        "selectElement"
      ) as HTMLSelectElement;
      selectT.value = "None";
      this.errorType = "Falta de Dados";
      this.messageError = "Necessário selecionar ao menos um computador";
      this.canViewMessage = true;
      return;
    }

    // Converte o array em uma string separada por vírgulas
    const selectedValuesString = selectedValues.join(",");

    this.http
      .post<ReportResponse>(
        "/home/computers/get-report/xls/",
        {
          selectedValues: selectedValuesString,
        },
        {
          headers: new HttpHeaders({
            "Content-Type": "application/json",
          }),
        }
      )
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((response) => {
        const fileName = response.file_name || "report.xlsx";
        const fileContent = response.file_content;

        // Converte a string base64 em um Blob
        const byteCharacters = atob(fileContent);
        const byteNumbers = new Array(byteCharacters.length);
        for (let i = 0; i < byteCharacters.length; i++) {
          byteNumbers[i] = byteCharacters.charCodeAt(i);
        }

        const byteArray = new Uint8Array(byteNumbers);
        const blob = new Blob([byteArray], {
          type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        });

        // Usa a biblioteca file-saver para salvar o arquivo
        saveAs(blob, fileName);
      }),
      (error: any) => {
        console.error("Download error:", error);
      };
  }

  // Envia o usuario e senha para o backend
  submitReportDNS(): void {
    this.canViewCredentialsLoading = true;
    this.http
      .post(
        "/home/computers/report-dns",
        {
          username: this.input_username,
          pwd: this.input_pwd,
        },
        {
          headers: new HttpHeaders({
            "X-CSRFToken": this.token,
            "Content-Type": "application/json",
          }),
        }
      )
      .pipe(
        catchError((error) => {
          this.status = error.status;

          if (this.status !== 200) {
            this.canViewCredentialsLoading = false;
          }
          if (this.status === 401) {
            this.errorType = "Invalid Credentials";
            this.messageError = "Usuario e/ou Senha Incorreto";
            this.canViewMessage = true;
          }
          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          const link = document.createElement("a");
          link.href = `data:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;base64,${data.filedata}`;
          link.download = data.filename;
          link.click();
          this.canViewCredentials = false;
          this.canViewCredentialsLoading = false;
        }
      });
  }

  selectAll(): void {
    const checkboxes = this.el.nativeElement.querySelectorAll(".ckip");
    checkboxes.forEach((checkbox: HTMLInputElement) => {
      if (this.checkedAll) {
        this.renderer.setProperty(checkbox, "checked", true);
      } else {
        this.renderer.setProperty(checkbox, "checked", false);
      }
    });
    if (this.checkedAll) {
      this.checkedAll = false;
    } else {
      this.checkedAll = true;
    }
  }
}
