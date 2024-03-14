import { Component } from '@angular/core';
import { ChartModule } from 'primeng/chart';
import { LoadingComponent } from '../loading/loading.component';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-chart-doughnut',
  standalone: true,
  imports: [ChartModule, LoadingComponent, CommonModule],
  templateUrl: './chart-doughnut.component.html',
  styleUrl: './chart-doughnut.component.css',
  schemas: [],
})
export class ChartDoughnutComponent {
  canShow: boolean = false;

  data: any;

  options: any;

  ngOnInit() {
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
