import { Component, OnInit, ViewChild } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { ChartComponent } from 'ng-apexcharts';
import {
  ApexAxisChartSeries,
  ApexChart,
  ApexDataLabels,
  ApexGrid,
  ApexStroke,
  ApexTitleSubtitle,
  ApexXAxis,
} from 'ng-apexcharts';
import { catchError, tap, throwError } from 'rxjs';

export type ChartOptions = {
  chart: ApexChart;
  dataLabels: ApexDataLabels;
  grid: ApexGrid;
  series: ApexAxisChartSeries;
  stroke: ApexStroke;
  title: ApexTitleSubtitle;
  xaxis: ApexXAxis;
};

@Component({
  selector: 'app-chart-spline-line',
  templateUrl: './chart-spline-line.component.html',
  styleUrls: ['./chart-spline-line.component.css'],
})
export class ChartSplineLineComponent implements OnInit {
  @ViewChild('chart') chart: ChartComponent | undefined;
  public chartOptions: ChartOptions;
  // Declarando varaiveis Dict
  dates: string[] = [];
  quantity: number[] = [];

  // Declarando varaiveis Boolean
  canViewChart: boolean = false;

  // Construtor monta o Chart Inicialmente
  constructor(private http: HttpClient) {
    this.chartOptions = {
      series: [
        {
          name: '',
          data: [],
        },
      ],
      chart: {
        height: 350,
        type: 'line',
        zoom: {
          enabled: false,
        },
      },
      dataLabels: {
        enabled: false,
        enabledOnSeries: [0],
        textAnchor: 'start',
        distributed: true,
      },
      stroke: {
        curve: 'stepline',
      },
      title: {
        text: '',
        align: 'left',
      },
      grid: {
        row: {
          colors: ['#f3f3f3', 'transparent'], // takes an array which will be repeated on columns
          opacity: 0.5,
        },
      },
      xaxis: {
        categories: [],
      },
    };
  }

  ngOnInit(): void {
    this.http
      .get<Array<{ date: string; count: number }>>(
        '/home/get-info-last-update',
        {
          headers: new HttpHeaders({
            Accept: 'application/json',
          }),
        }
      )
      .pipe(
        tap((data) => {
          this.processData(data);
          this.setupChart();
        }),
        catchError((error) => {
          // this.status = error.status;
          // this.canViewChart = false;
          return throwError(error);
        })
      )
      .subscribe(() => {
        this.canViewChart = true;
      });
  }

  // Faz o processamento dos dados, segregando os nomes e quantidade
  private processData(data: { date: string; count: number }[]): void {
    this.dates = data.map((item) => item.date);
    this.quantity = data.map((item) => item.count);
  }

  private setupChart(): void {
    this.chartOptions = {
      series: [
        {
          name: 'Desktops',
          data: this.quantity,
          type: 'line',
        },
      ],
      chart: {
        height: 350,
        type: 'line',
        zoom: {
          enabled: false,
        },
      },
      dataLabels: {
        enabled: false,
        enabledOnSeries: [0],
        textAnchor: 'start',
        distributed: true,
      },
      stroke: {
        curve: 'stepline',
      },
      title: {
        text: 'Entradas de computadores',
        align: 'left',
        style: {
          color: '#FFFFFF',
          fontWeight: 'bold',
        },
      },
      grid: {
        row: {
          colors: ['#f3f3f3', 'transparent'], // takes an array which will be repeated on columns
          opacity: 0.5,
        },
      },
      xaxis: {
        categories: this.dates,
        labels: {
          style: {
            colors: '#FFFF00',
          },
        },
      },
    };
  }
}
