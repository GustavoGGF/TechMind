import {
  CUSTOM_ELEMENTS_SCHEMA,
  Component,
  ElementRef,
  Inject,
  ViewChild,
} from '@angular/core';
import { UtilitiesModule } from '../../utilities/utilities.module';
import { CommonModule, DOCUMENT } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { catchError, throwError } from 'rxjs';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [UtilitiesModule, CommonModule, HttpClientModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
})
export class LoginComponent {
  constructor(
    private elementRef: ElementRef,
    private http: HttpClient,
    @Inject(DOCUMENT) private document: Document
  ) {}
  messageError: string = '';
  name: string = '';
  pass: string = '';
  typError: string = '';
  urlImage = '/static/assets/images/Logo_TechMind.png';

  canShow: boolean = false;
  canShowMessage: boolean = false;

  status: any;

  @ViewChild('logo') logo: ElementRef | undefined;
  @ViewChild('main') main: ElementRef | undefined;

  ngAfterViewInit(): void {
    if (this.logo) {
      this.logo.nativeElement.addEventListener('animationend', () => {
        const letters =
          this.elementRef.nativeElement.querySelectorAll('.letter');

        letters.forEach((letter: any) => {
          (letter as HTMLElement).classList.add('animate1');
        });
      });
    }

    if (this.document) {
      this.document.addEventListener('keyup', (event: any) => {
        if (event.keyCode === 13) {
          this.loginSubmit();
        }
      });
    }
  }

  closeMessage(): void {
    this.canShowMessage = false;
    this.canShow = false;
  }

  getUser(event: any) {
    this.name = event.target.value;
    return this.name;
  }

  getPass(event: any) {
    this.pass = event.target.value;
    return this.pass;
  }

  loginSubmit(): void {
    if (this.name.length > 1 && this.pass.length > 1) {
      this.canShow = true;

      this.http
        .post('/credential/', {
          username: this.name,
          password: this.pass,
        })
        .pipe(
          catchError((error) => {
            this.status = error.status;

            if (this.status === 401) {
              this.typError = 'Erro de Credencial';
              this.messageError = 'Senha e/ou Usuário inválido';
            } else if (this.status === 404) {
              this.typError = 'Erro de Recurso';
              this.messageError = 'Recurso não encontrado';
            } else if (this.status === 500) {
              this.typError = 'Erro de Servidor';
              this.messageError = 'Erro interno do servidor';
            } else {
              this.typError = 'Erro Desconhecido';
              this.messageError = 'Ocorreu um erro desconhecido';
            }

            this.canShowMessage = true;
            this.canShow = false;
            return throwError(error);
          })
        )
        .subscribe((data: any) => {
          if (data.name) {
            localStorage.setItem('name', data.name);

            window.location.href = '/home';
          }
        });
    }
  }
}
