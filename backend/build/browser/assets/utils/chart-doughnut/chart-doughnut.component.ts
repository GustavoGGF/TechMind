import { Component } from '@angular/core';
import { ChartModule } from 'primeng/chart';
import { LoadingComponent } from '../loading/loading.component';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';

@Component({
  selector: 'app-chart-doughnut',
  standalone: true,
  imports: [ChartModule, LoadingComponent, CommonModule, HttpClientModule],
  templateUrl: './chart-doughnut.component.html',
  styleUrl: './chart-doughnut.component.css',
  schemas: [],
})
export class ChartDoughnutComponent {
  constructor(private http: HttpClient) {}
  canShow: boolean = false;

  data: any;
  status: any;
  options: any;

  ngOnInit() {
    this.canShow = true;

    const currentUrl = window.location.href;

    this.http.get(currentUrl).subscribe((data: any) => {
      this.status = data.status;

      if (this.status === 401) {
      } else if (this.status === 404) {
      } else if (this.status === 200) {
        console.log(data.systems);
      } else {
        console.log(this.status);
        console.log(this.status.status);
      }
    });

    const documentStyle = getComputedStyle(document.documentElement);

    const color: String = 'var(--black)';

    this.data = {
      labels: ['A', 'B', 'C'],
      datasets: [
        {
          data: [300, 50, 100],
          backgroundColor: [
            documentStyle.getPropertyValue('--blue-500'),
            documentStyle.getPropertyValue('--yellow-500'),
            documentStyle.getPropertyValue('--green-500'),
          ],
          hoverBackgroundColor: [
            documentStyle.getPropertyValue('--blue-400'),
            documentStyle.getPropertyValue('--yellow-400'),
            documentStyle.getPropertyValue('--green-400'),
          ],
        },
      ],
    };

    this.options = {
      cutout: '60%',
      plugins: {
        legend: {
          labels: {
            color: color,
          },
        },
      },
    };
  }
}
