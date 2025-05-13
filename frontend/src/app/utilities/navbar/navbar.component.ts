import { HttpClient } from "@angular/common/http";
import { Component, Input } from "@angular/core";
import { catchError, throwError } from "rxjs";

@Component({
  selector: "app-navbar",
  templateUrl: "./navbar.component.html",
  styleUrl: "./navbar.component.css",
})
export class NavbarComponent {
  // Declarando e exportando as variaveis que serão usadas no componente
  @Input() name: string = "";
  @Input() home_class: string = "";
  @Input() computers_class: string = "";
  @Input() device_class: string = "";
  @Input() panel_class: string = "";

  constructor(private http: HttpClient) {}

  logoTechMind: string = "/static/assets/images/Logo_TechMind.png";
  logoutIMG: string = "/static/assets/images/logout.png";

  status: number = 0;

  logoutApp(): void {
    // Função para realizar o logout do usuário
    this.http
      .get("/logout/", {})
      .pipe(
        catchError((error) => {
          return throwError(error);
        })
      )
      .subscribe({
        next: () => {
          // Sucesso no logout, você pode redirecionar ou fazer outra ação
          return (window.location.href = "/login");
        },
        error: (error) => {
          console.error("Erro ao realizar logout", error);
        },
      });
  }
}
