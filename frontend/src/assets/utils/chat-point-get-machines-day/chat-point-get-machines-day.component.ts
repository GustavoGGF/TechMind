import { Component } from '@angular/core';
import { ChartModule } from 'primeng/chart';
import { LoadingComponent } from '../loading/loading.component';
import { HttpClientModule, HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-chat-point-get-machines-day',
  standalone: true,
  imports: [ChartModule, LoadingComponent, HttpClientModule, CommonModule],
  templateUrl: './chat-point-get-machines-day.component.html',
  styleUrl: './chat-point-get-machines-day.component.css',
})
export class ChatPointGetMachinesDayComponent {
  constructor(private http: HttpClient) {}
  canShow: boolean = false;

  data: any;
  options: any;
  status: any;
  dateInsertition: any;
  countDates: any;

  ngOnInit() {
    this.canShow = true;

    this.http.get('api/get-machines-days').subscribe((data: any) => {
      this.status = data.status;

      switch (this.status) {
        case 200:
          this.dateInsertition = data.dateInsertition;
          this.countDates = data.countDates;
          this.canShow = false;
      }
    });

    this.data = {
      labels: [this.dateInsertition],
      datasets: [
        {
          label: 'Dataset',
          data: this.countDates,
          borderColor: 'red',
          backgroundColor: 'rgba(255, 0, 0, 0.5)',
          pointStyle: 'circle',
          pointRadius: 10,
          pointHoverRadius: 15,
        },
      ],
    };

    this.options = {
      type: 'line',
      data: this.data,
      options: {
        responsive: true,
        plugins: {
          title: {
            display: true,
            text: (ctx: any) =>
              'Point Style: ' + ctx.chart.data.datasets[0].pointStyle,
          },
        },
      },
    };
  }
}
