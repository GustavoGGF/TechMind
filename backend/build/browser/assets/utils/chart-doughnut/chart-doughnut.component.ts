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
  counts: number[] = [];
  systemNames: string[] = [];

  ngOnInit() {
    this.canShow = true;

    const currentUrl = window.location.href;

    this.http.get(currentUrl).subscribe((data: any) => {
      this.status = data.status;

      if (this.status === 401) {
      } else if (this.status === 404) {
      } else if (this.status === 200) {
        for (let item of data.systems) {
          this.systemNames.push(item.SystemName);
          this.counts.push(item.Count);
          this.canShow = false;
        }
      } else {
        console.log(this.status);
        console.log(this.status.status);
      }
    });

    const documentStyle = getComputedStyle(document.documentElement);

    const color: String = 'var(--black)';

    this.data = {
      labels: this.systemNames,
      datasets: [
        {
          data: this.counts,
          backgroundColor: [
            documentStyle.getPropertyValue('--blue-500'),
            documentStyle.getPropertyValue('--yellow-500'),
          ],
          hoverBackgroundColor: [
            documentStyle.getPropertyValue('--blue-400'),
            documentStyle.getPropertyValue('--yellow-400'),
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
