import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { RouterModule } from '@angular/router';
import { routes } from './app.routes';
import { PagesModule } from './pages/pages.module';

@NgModule({
  imports: [BrowserModule, RouterModule.forRoot(routes), PagesModule],
})
// Aplicação principal, tudo se origina daqui
export class AppModule {}
