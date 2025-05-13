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

  name: any;

  errorType: string = "";
  messageError: string = "";
  home_class: string = "";
  computers_class: string = "";
  device_class: string = "";
  panel_class: string = "active";
  statusDot: string = "";
  statusDotTitle: string = "";
  statusDotContent: string = "";
  statusDotPing: string = "";

  canView: boolean = false;
  showMessage: boolean = false;

  status: number = 0;

  machines: any;

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
      this.getMachines();
    }

    document.addEventListener("click", (event: MouseEvent) => {
      return this.closePoPOver(event);
    });
  }

  // Função para fechar a mensagem de erro
  hideMessage() {
    this.showMessage = false;
  }

  activePopOver(element: HTMLElement) {
    const popover = bootstrap.Popover.getOrCreateInstance(element);
    popover.show();
  }

  closePoPOver(event: MouseEvent) {
    const target = event.target as HTMLElement;

    if (
      target.id !== "status_dot" &&
      !(
        target.classList.contains("popover-header") ||
        target.classList.contains("popover-body")
      )
    ) {
      const popoverElement = document.getElementById("status_dot");
      if (popoverElement) {
        const popover = bootstrap.Popover.getInstance(popoverElement);
        if (popover) {
          popover.hide();
        }
      }
    }
  }

  getMachines() {
    this.http
      .get("/home/panel-adm/get-machines", {})
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
          this.machines = this.groupSplitter(data.machines[0], 100);
        }
      });
  }

  groupSplitter<T>(array: T[], size: number): T[][] {
    const groups: T[][] = [];
    for (let i = 0; i < array.length; i += size) {
      groups.push(array.slice(i, i + size));
    }
    return groups;
  }

  adjustDate(curentDate: string) {
    const date = new Date(curentDate);
    const now = new Date();
    var daysInactive: string = "";

    // diferença em milissegundos
    const diffMs = now.getTime() - date.getTime();

    // 48 horas em milissegundos = 48 * 60 * 60 * 1000
    const isWithin48Hours = diffMs <= 48 * 60 * 60 * 1000;
    if (isWithin48Hours) {
      const hours = String(date.getHours()).padStart(2, "0");
      const minutes = String(date.getMinutes()).padStart(2, "0");
      const time = `${hours}:${minutes}`;
      this.statusDot = "status-dot-active";
      this.statusDotTitle = "Online";
      this.statusDotPing = "status-ping-active";
      this.statusDotContent =
        "Esse Status representa que o equipamento está ou ficou online nas últimas 48 Horas";
      return time;
    } else {
      // // Calcular a diferença em milissegundos
      // const differenceInMillis = now.getTime() - date.getTime();

      // // Converter de milissegundos para dias
      // const differenceInDays = Math.floor(
      //   differenceInMillis / (1000 * 3600 * 24)
      // );
      daysInactive = this.getRelativeDateString(curentDate);
      // daysInactive = this.getDateInactive(differenceInDays);
      this.statusDot = "status-dot-inactive";
      this.statusDotPing = "status-ping-inactive";
      this.statusDotTitle = "Offline";
      this.statusDotContent =
        "Esse Status representa que o equipamento não ficou online nas ultimas 48 horas.";
      return daysInactive;
    }
  }

  getRelativeDateString(dateString: string): string {
    dayjs.extend(relativeTime);
    dayjs.locale("pt-br");
    return dayjs(dateString).fromNow(); // Ex: "há um mês"
  }

  getDateInactive(date: number) {
    if (date >= 7 && date < 14) {
      return "1 semana atrás";
    } else if (date >= 14 && date < 21) {
      return "2 semanas atrás";
    } else if (date >= 21 && date < 30) {
      return "3 semanas atrás";
    } else if (date == 30 || date == 31) {
      return "1 mês atrás";
    } else if (date < 60) {
      return "1 mês atrás";
    } else if (date >= 60 && date < 90) {
      return "2 meses atrás";
    } else if (date >= 90 && date < 120) {
      return "3 meses atrás";
    } else if (date >= 120 && date < 150) {
      return "4 meses atrás";
    } else if (date >= 150 && date < 180) {
      return "5 meses atrás";
    } else if (date >= 180 && date < 210) {
      return "6 meses atrás";
    } else if (date >= 210 && 240) {
      return "7 meses atrás";
    } else if (date >= 240 && date < 270) {
      return "8 meses atrás";
    } else if (date >= 270 && date < 300) {
      return "9 meses atrás";
    } else if (date >= 300 && date < 330) {
      return "10 meses atrás";
    } else if (date >= 330 && date < 360) {
      return "11 meses atrás";
    } else if (date >= 360 && date < 720) {
      return "1 ano atrás";
    } else {
      return "";
    }
  }
}
