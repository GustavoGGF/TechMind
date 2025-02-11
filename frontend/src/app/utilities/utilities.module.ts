import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { NgApexchartsModule } from 'ng-apexcharts';
import { ChartPieComponent } from './chart-pie/chart-pie.component';
import { LoadingComponent } from './loading/loading.component';
import { LoadingSearchComponent } from './loading-search/loading-search.component';
import { MessageComponent } from './message/message.component';
import { NavbarComponent } from './navbar/navbar.component';
import { ChartSplineLineComponent } from './chart-spline-line/chart-spline-line.component';
import { LoadingPerfectApeComponent } from './loading-perfect-ape/loading-perfect-ape.component';

// Modulo que gerencia os utilitarios Disponibilizando eles onde o Modulo for importado
@NgModule({
  declarations: [
    ChartPieComponent,
    LoadingComponent,
    LoadingSearchComponent,
    MessageComponent,
    NavbarComponent,
    ChartSplineLineComponent,
    LoadingPerfectApeComponent,
  ],
  imports: [CommonModule, NgApexchartsModule],
  exports: [
    ChartPieComponent,
    LoadingComponent,
    LoadingSearchComponent,
    MessageComponent,
    NavbarComponent,
    ChartSplineLineComponent,
    LoadingPerfectApeComponent,
  ],
})
export class UtilitiesModule {}
