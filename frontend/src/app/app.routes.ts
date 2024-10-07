import { RouterModule, Routes } from '@angular/router';
import { LoginComponent } from './pages/login/login.component';
import { HomeComponent } from './pages/home/home.component';
import { NgModule } from '@angular/core';
import { ComputersComponent } from './pages/computers/computers.component';
import { ComputersDetailsComponent } from './pages/computers-details/computers-details.component';
import { DevicesComponent } from './pages/devices/devices.component';
import { DevicesDetailsComponent } from './pages/devices-details/devices-details.component';

export const routes: Routes = [
  { path: '', component: LoginComponent, title: 'Login' },
  { path: 'login', component: LoginComponent, title: 'Login' },
  { path: 'home', component: HomeComponent, title: 'Dashboard' },
  {
    path: 'home/computers',
    component: ComputersComponent,
    title: 'Computadores',
  },
  {
    path: 'home/computers/view-machine/:mac',
    component: ComputersDetailsComponent,
    title: 'Detalhes do Computador',
  },
  {
    path: 'home/devices',
    component: DevicesComponent,
    title: 'Dispositivos',
  },
  {
    path: 'home/devices/view-devices/:sn',
    component: DevicesDetailsComponent,
    title: 'Detalhes do Dispositivo',
  },
];

NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
});
