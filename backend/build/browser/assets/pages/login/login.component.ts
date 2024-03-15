import {
  CUSTOM_ELEMENTS_SCHEMA,
  Component,
  ElementRef,
  ViewChild,
} from '@angular/core';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Router } from '@angular/router';
import { MessageComponent } from '../../utils/message/message.component';
import { CommonModule } from '@angular/common';
import { LoadingComponent } from '../../utils/loading/loading.component';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [HttpClientModule, MessageComponent, CommonModule, LoadingComponent],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
})
export class LoginComponent {
  urlImage = '../assets/images/logo/Logo_TechMind.png';
  canShow: boolean = false;
  status: any;
  typError: string = '';
  messageError: string = '';
  canShowMessage: boolean = false;
  name: string = '';
  pass: string = '';

  constructor(
    private elementRef: ElementRef,
    private http: HttpClient,
    private router: Router
  ) {}

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

    const body = document.body;

    if (body) {
      body.addEventListener('keyup', (event: any) => {
        if (event.keyCode === 13) {
          this.loginSubmit();
        }
      });
    }
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
      const currentUrl = window.location.href;

      this.canShow = true;

      this.http
        .post(currentUrl + 'api/credential', {
          user: this.name,
          pass: this.pass,
        })
        .subscribe((data: any) => {
          this.status = data.status;

          if (this.status === 401) {
            this.typError = 'Erro de Credencial';
            this.messageError = 'Você não possui acesso ao TechMind';
            this.canShowMessage = true;
          } else if (this.status === 404) {
            this.typError = 'Erro de Credencial';
            this.messageError = 'Senha e/ou Usuário inválido';
            this.canShowMessage = true;
            console.log(this.canShowMessage);
          } else if (this.status === 200) {
            localStorage.setItem('name', JSON.stringify(data.name));
            localStorage.setItem('csrf', JSON.stringify(data.csrf_token));
            this.router.navigateByUrl('/home');
          } else {
            console.log(this.status);
            console.log(this.status.status);
          }
        });
    }
  }

  closeMessage(): void {
    this.canShowMessage = false;
    this.canShow = false;
  }
}
