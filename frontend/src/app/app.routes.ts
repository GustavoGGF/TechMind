import { RouterModule, Routes } from '@angular/router';
import { LoginComponent } from '../assets/pages/login/login.component';
import { HomeComponent } from '../assets/pages/home/home.component';
import { NgModule } from '@angular/core';

export const routes: Routes = [
  { path: '', component: LoginComponent, title: 'Login' },
  { path: 'home', component: HomeComponent, title: 'Dashboard' },
];

NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
});
