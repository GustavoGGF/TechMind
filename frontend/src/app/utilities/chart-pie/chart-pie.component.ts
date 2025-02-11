import { Component, OnInit, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {
  ApexChart,
  // ApexDataLabels,
  ApexFill,
  ApexLegend,
  ApexNonAxisChartSeries,
  ApexResponsive,
  ApexStroke,
  ApexTheme,
  ChartComponent,
} from 'ng-apexcharts';
import { catchError, tap, throwError } from 'rxjs';

export type ChartOptions = {
  series: ApexNonAxisChartSeries;
  chart: ApexChart;
  responsive: ApexResponsive[];
  labels: any;
  stroke: ApexStroke;
  fill: ApexFill;
  legend: ApexLegend;
  // dataLabels: ApexDataLabels;
  theme: ApexTheme;
};

@Component({
  selector: 'app-chart-pie',
  templateUrl: './chart-pie.component.html',
  styleUrls: ['./chart-pie.component.css'],
})
export class ChartPieComponent implements OnInit {
  @ViewChild('chart') chart: ChartComponent | undefined;
  public chartOptions: ChartOptions;

  // Declarando Variaveis Boolean
  canViewChart: boolean = false;
  // Declarando Varaiveis Any
  status: any;
  // Declarando Variaveis Dict
  array_so_names: string[] = [];
  array_so_quantity: number[] = [];

  constructor(private http: HttpClient) {
    // Iniciando chart
    this.chartOptions = {
      series: [],
      chart: {
        type: 'polarArea',
      },
      stroke: {
        colors: ['#fff'],
      },
      fill: {
        opacity: 0.8,
      },
      responsive: [
        {
          breakpoint: 480,
          options: {
            chart: {},
            legend: {
              position: 'bottom',
            },
          },
        },
      ],
      labels: [],
      // dataLabels: {
      //   style: {
      //     colors: ['#FFFF00'],
      //   },
      // },
      legend: {
        labels: {
          colors: '#3083DC', // Cor da fonte da legenda
        },
      },
      theme: {
        mode: 'dark',
      },
    };
  }

  // Função que inicia ao entrar na pagina e faz uma requisição dos dados para o dashboard
  ngOnInit(): void {
    this.http
      .get<Array<{ system_name: string; count: number }>>(
        '/home/get-info-SO',
        {}
      )
      .pipe(
        tap((data) => {
          this.processData(data);
          this.setupChart();
        }),
        catchError((error) => {
          this.status = error.status;
          this.canViewChart = false;
          return throwError(error);
        })
      )
      .subscribe(() => {
        this.canViewChart = true;
      });
  }

  // Faz o processamento dos dados, segregando os nomes e quantidade
  private processData(data: { system_name: string; count: number }[]): void {
    this.array_so_names = data.map((item) => item.system_name);
    this.array_so_quantity = data.map((item) => item.count);
  }

  // Montando os dados do Chart
  private setupChart(): void {
    this.chartOptions = {
      series: this.array_so_quantity,
      chart: {
        type: 'polarArea',
      },
      stroke: {
        colors: ['#fff'],
      },
      fill: {
        opacity: 0.8,
      },
      responsive: [
        {
          breakpoint: 480,
          options: {
            chart: {},
            legend: {
              position: 'bottom',
            },
          },
        },
      ],
      labels: this.array_so_names,
      // dataLabels: {
      //   style: {
      //     colors: ['#FFFF00'],
      //   },
      // },
      legend: {
        labels: {
          colors: '#3083DC', // Cor da fonte da legenda
        },
      },
      theme: {
        mode: 'dark',
      },
    };
  }
}
