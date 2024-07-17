import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MessageComponent } from './message/message.component';
import { LoadingComponent } from './loading/loading.component';
import { NavbarComponent } from './navbar/navbar.component';

@NgModule({
  declarations: [MessageComponent, LoadingComponent, NavbarComponent],
  imports: [CommonModule],
  exports: [MessageComponent, LoadingComponent, NavbarComponent],
})
export class UtilitiesModule {}
