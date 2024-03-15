import { Component } from '@angular/core';
import { ChartModule } from 'primeng/chart';
import { LoadingComponent } from '../loading/loading.component';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';

@Component({
  selector: 'app-chart-doughnut-cities',
  standalone: true,
  imports: [ChartModule, LoadingComponent, CommonModule, HttpClientModule],
  templateUrl: './chart-doughnut-cities.component.html',
  styleUrl: './chart-doughnut-cities.component.css',
})
export class ChartDoughnutCitiesComponent {
  constructor(private http: HttpClient) {}
  canShow: boolean = false;

  data: any;
  status: any;
  options: any;
  counts: number[] = [];
  systemNames: string[] = [];

  ngOnInit() {
    this.canShow = true;

    this.http.get('api/cities').subscribe((data: any) => {
      this.status = data.status;

      if (this.status === 401) {
      } else if (this.status === 404) {
      } else if (this.status === 200) {
        for (let item of data.cities) {
          console.log(data);

          this.systemNames.push(item.CityName);
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
      cutout: '40%',
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
