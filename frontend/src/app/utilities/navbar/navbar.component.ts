import { HttpClient } from '@angular/common/http';
import { Component, Input } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css',
})
export class NavbarComponent {
  // Declarando e exportando as variaveis que serão usadas no componente
  @Input() name: string = '';
  @Input() home_class: string = '';
  @Input() computers_class: string = '';
  @Input() device_class: string = '';

  constructor(private http: HttpClient, private router: Router) {}

  logoTechMind: string = '/static/assets/images/Logo_TechMind.png';
  logoutIMG: string = '/static/assets/images/logout.png';

  status: number = 0;

  logoutApp(): void {
    // Função para realizar o logout do usuário
    this.http.get('/logout', {}).pipe(
      catchError((error) => {
        this.status = error.status;

        if (this.status === 200) {
          this.router.navigate(['/login']);
        }

        return throwError(error);
      })
    );
  }
}
