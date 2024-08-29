import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { RouterModule } from '@angular/router';
import { routes } from './app.routes';
import { PagesModule } from './pages/pages.module';
import { SocketIoModule, SocketIoConfig } from 'ngx-socket-io';

const config: SocketIoConfig = {
  url: 'http://localhost:3000/home',
  options: {},
};

@NgModule({
  imports: [
    BrowserModule,
    RouterModule.forRoot(routes),
    PagesModule,
    SocketIoModule.forRoot(config),
  ],
  providers: [],
})
// Aplicação principal, tudo se origina daqui
export class AppModule {}
