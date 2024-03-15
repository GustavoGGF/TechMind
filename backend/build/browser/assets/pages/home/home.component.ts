import { Component } from '@angular/core';
import { NavbarComponent } from '../../utils/navbar/navbar.component';
import { ChartDoughnutSOComponent } from '../../utils/chart-doughnut-SO/chart-doughnut.component';
import { CommonModule } from '@angular/common';
import {
  HttpClientModule,
  HttpClient,
  HttpHeaders,
} from '@angular/common/http';
import { ChartDoughnutCitiesComponent } from '../../utils/chart-doughnut-cityes/chart-doughnut-cities.component';
import { Router } from '@angular/router';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    NavbarComponent,
    ChartDoughnutSOComponent,
    CommonModule,
    HttpClientModule,
    ChartDoughnutCitiesComponent,
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
  constructor(private http: HttpClient, private router: Router) {}

  Csrf: any;

  ngAfterViewInit(): void {
    const Token = localStorage.getItem('csrf');
    if (Token !== null) {
      this.Csrf = JSON.parse(Token);
    } else {
      this.router.navigateByUrl('/');
    }

    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'X-CSRF-Token': this.Csrf,
    });

    this.http.post('home', { data: '' }, { headers: headers });

    console.log('depois do post');

    this.http.get('api/home').subscribe((data: any) => {
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
