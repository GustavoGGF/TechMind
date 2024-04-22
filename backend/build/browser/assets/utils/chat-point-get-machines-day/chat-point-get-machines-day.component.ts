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

  ngOnInit() {
    this.canShow = true;

    this.http.get('api/get-machines-days').subscribe((data: any) => {});
  }
}
