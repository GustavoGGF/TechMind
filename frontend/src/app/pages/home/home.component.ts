import { Component, OnInit } from "@angular/core";
import { UtilitiesModule } from "../../utilities/utilities.module";
import { CommonModule } from "@angular/common";
import { HttpClient, HttpClientModule } from "@angular/common/http";
import { catchError, throwError } from "rxjs";

@Component({
  selector: "app-home",
  standalone: true,
  imports: [UtilitiesModule, CommonModule, HttpClientModule],
  templateUrl: "./home.component.html",
  styleUrl: "./home.component.css",
})
export class HomeComponent implements OnInit {
  constructor(private http: HttpClient) {}
  // Declarando variavel any
  name: any;
  status: any;

  // Declarando variaveis string
  computers_class: string = "";
  device_class: string = "";
  errorType: string = "";
  home_class: string = "active";
  messageError: string = "";
  message: string = "";

  // Declarando variaveis boolean
  canView: boolean = false;
  notData: boolean = true;
  showMessage: boolean = false;

  // Variaveis number
  totalMachines: number = 0;
  totalWindows: number = 0;
  totalUnix: number = 0;

  // Função inicia ao iniciar o componente
  ngOnInit() {
    // Pegando os dados do usuario
    this.name = localStorage.getItem("name");
    // Verificando se os dados existem
    if (this.name.length == 0 || this.name == null) {
      this.errorType = "Falta de Dados";
      this.messageError =
        "Ouve um erro ao acessar dados do LDAP, contatar a TI";
      this.showMessage = true;
    } else {
      this.canView = true;
      // chamando a função para obter os dados
      this.startPolling();
    }
  }

  ngAfterViewInit() {}

  // Função para fechar a mensagem de erro
  hideMessage() {
    this.showMessage = false;
  }

  // Função que os dados do dashboard
  getData() {
    this.http
      .get("/home/get-Info-Main-Panel/", {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.totalWindows = data.windows;
          this.totalMachines = data.total;
          this.totalUnix = data.unix;
          this.notData = false;
          this.verifyColorPie();
        }
      });
  }

  // Polling a cada 10 segundos (10000 ms)
  startPolling() {
    this.getData(); // Chamada inicial
    setInterval(() => {
      this.getData();
    }, 60000); // 10 segundos
  }

  verifyColorPie() {
    const legendTextElements = document.getElementsByClassName(
      "apexcharts-legend-text"
    );

    Array.from(legendTextElements).forEach((element) => {
      // Verifica se o elemento é um HTMLElement
      if (element instanceof HTMLElement) {
        // Aplica a cor amarelo com !important
        element.style.setProperty("color", "yellow", "important");
      }
    });

    Array.from(legendTextElements).forEach((element) => {
      if (element instanceof HTMLElement) {
        // Obtém o estilo computado do elemento
        const color = window.getComputedStyle(element).color;

        if (color !== "rgb(255, 255, 0)") {
          this.verifyColorPie();
        } else {
        }
      }
    });
  }
}
