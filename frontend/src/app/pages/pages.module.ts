import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LoginComponent } from './login/login.component';
import { HomeComponent } from './home/home.component';
import { UtilitiesModule } from '../utilities/utilities.module';
import { ComputersComponent } from './computers/computers.component';
import { ComputersDetailsComponent } from './computers-details/computers-details.component';
import { DevicesComponent } from './devices/devices.component';
import { DevicesDetailsComponent } from './devices-details/devices-details.component';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    LoginComponent,
    HomeComponent,
    UtilitiesModule,
    ComputersComponent,
    ComputersDetailsComponent,
    DevicesComponent,
    DevicesDetailsComponent,
  ],
  providers: [],
})
export class PagesModule {}
