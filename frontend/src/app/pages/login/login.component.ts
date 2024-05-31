import { CommonModule, DOCUMENT } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import {
  CUSTOM_ELEMENTS_SCHEMA,
  Component,
  ElementRef,
  Inject,
  ViewChild,
} from '@angular/core';
import { MessageComponent } from '../../shared/message/message.component';
import { LoadingComponent } from '../../shared/loading/loading.component';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, HttpClientModule, MessageComponent, LoadingComponent],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
})
export class LoginComponent {
  urlImage = '../../../assets/images/Logo_TechMind.png';
  canShow: boolean = false;
  canShowMessage: boolean = false;
  typError: string = '';
  messageError: string = '';
  name: string = '';
  pass: string = '';
  status: any;
  authTokenKey: string = 'authToken';

  constructor(
    private elementRef: ElementRef,
    private http: HttpClient,
    @Inject(DOCUMENT) private document: Document
  ) {}

  @ViewChild('logo') logo: ElementRef | undefined;
  @ViewChild('main') main: ElementRef | undefined;

  ngOninit(): void {
    localStorage.setItem('authToken', JSON.stringify(''));
    // Após o login bem-sucedido, armazene o token em um cookie
    this.document.cookie = `authToken=${''};max-age=0; path=/;`;
  }

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

  getUser(event: any) {
    this.name = event.target.value;
    return this.name;
  }

  getPass(event: any) {
    this.pass = event.target.value;
    return this.pass;
  }

  closeMessage(): void {
    this.canShowMessage = false;
    this.canShow = false;
  }

  loginSubmit(): void {
    if (this.name.length > 1 && this.pass.length > 1) {
      const currentUrl = window.location.href;

      this.canShow = true;

      this.http
        .post(currentUrl + 'api/credential', {
          username: this.name,
          password: this.pass,
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
            const newData = data.data;
            localStorage.setItem('name', JSON.stringify(newData.name));
            localStorage.setItem('authToken', JSON.stringify(newData.token));
            // Após o login bem-sucedido, armazene o token em um cookie
            this.document.cookie = `authToken=${newData.token};max-age=3600; path=/;`;
            window.location.href = '/home';
          } else {
            console.log(this.status);
            console.log(this.status.status);
          }
        });
    }
  }

  // Recupera o token de autenticação do cookie
  getToken(): string | null {
    return (
      this.document.cookie
        .split('; ')
        .find((row) => row.startsWith(`${this.authTokenKey}=`))
        ?.split('=')[1] || null
    );
  }
}
