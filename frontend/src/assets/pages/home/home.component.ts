import { Component } from '@angular/core';
import { NavbarComponent } from '../../utils/navbar/navbar.component';
import { ChartDoughnutComponent } from '../../utils/chart-doughnut/chart-doughnut.component';
import { CommonModule } from '@angular/common';
import { HttpClientModule, HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    NavbarComponent,
    ChartDoughnutComponent,
    CommonModule,
    HttpClientModule,
  ],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css',
  schemas: [],
})
export class HomeComponent {
  status: any;
  machinesLinux: number = 0;
  machinesWindows: number = 0;
  machinesAll: number = 0;
  constructor(private http: HttpClient) {}

  ngAfterViewInit(): void {
    const currentUrl = window.location.href;
    this.http.get(currentUrl).subscribe((data: any) => {
      this.status = data.status;

      if (this.status === 401) {
      } else if (this.status === 404) {
      } else if (this.status === 200) {
        this.machinesLinux = data.linux;
        this.machinesAll = data.machines;
        this.machinesWindows = data.windows;
      } else {
        console.log(this.status);
        console.log(this.status.status);
      }
    });
  }
}
