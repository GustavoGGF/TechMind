import { Injectable } from '@angular/core';
import {
  HttpInterceptor,
  HttpRequest,
  HttpHandler,
  HttpEvent,
} from '@angular/common/http';
import { Observable } from 'rxjs';
import { LoginComponent } from './pages/login/login.component';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  constructor(private loginService: LoginComponent) {}

  intercept(
    request: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    const token = this.loginService.getToken();

    console.log('token intercept:', token);

    if (token) {
      // Adiciona o token de autenticação ao cabeçalho Authorization
      request = request.clone({
        setHeaders: {
          Authorization: `Bearer ${token}`,
        },
      });
    } else {
      console.log('sem token');
    }

    return next.handle(request);
  }
}
