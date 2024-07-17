import { Component, ElementRef, ViewChild } from '@angular/core';
import { UtilitiesModule } from '../../utilities/utilities.module';
import { CommonModule } from '@angular/common';
import {
  HttpClient,
  HttpClientModule,
  HttpHeaders,
} from '@angular/common/http';
import { catchError, throwError } from 'rxjs';

@Component({
  selector: 'app-devices',
  standalone: true,
  imports: [UtilitiesModule, CommonModule, HttpClientModule],
  templateUrl: './devices.component.html',
  styleUrl: './devices.component.css',
})
export class DevicesComponent {
  constructor(private http: HttpClient) {}

  name: any;
  canView: boolean = false;
  imageDevices = '/static/assets/images/devices/computador.png';
  canViewNewDevices: boolean = false;
  close_button: string = '/static/assets/images/fechar.png';
  select_value: string = '';
  input_model: string = '';
  input_serial_number: string = '';
  input_imob: string = '';
  status: any;
  input_brand: string = '';
  token: any;
  devices: any;
  canViewDevices: boolean = false;

  ngOnInit(): void {
    this.name = localStorage.getItem('name');

    if (this.name.length == 0 || this.name == null) {
    } else {
      this.canView = true;
      this.getToken();
    }
  }

  getToken(): void {
    this.http
      .get('/home/devices/get-token', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          if (this.status === 0) {
          }

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.token = data.token;
        }
        this.getDevices();
      });
  }

  getDevices(): void {
    this.http
      .get('/home/devices/get-devices', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          if (this.status === 0) {
          }

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.devices = data.Dispositivos;
          this.canViewDevices = true;
          console.log(this.devices);
        }
      });
  }

  addDevices(): void {
    this.canViewNewDevices = true;
  }

  closeNewDevices(): void {
    this.canViewNewDevices = false;
  }

  device_select(event: any): void {
    this.select_value = event.target.value;
  }

  getModel(event: any): void {
    this.input_model = event.target.value;
  }

  getSerialNumber(event: any): void {
    this.input_serial_number = event.target.value;
  }

  getImob(event: any): void {
    this.input_imob = event.target.value;
  }

  getBrand(event: any): void {
    this.input_brand = event.target.value;
  }

  sumbitNewDevice(): void {
    console.log(this.select_value);

    this.http
      .post(
        '/home/devices/post-devices',
        {
          device: this.select_value,
          model: this.input_model,
          serial_number: this.input_serial_number,
          imob: this.input_imob,
          brand: this.input_brand,
        },
        {
          headers: new HttpHeaders({
            'X-CSRFToken': this.token,
            'Content-Type': 'application/json',
          }),
        }
      )
      .pipe(
        catchError((error) => {
          this.status = error.status;

          if (this.status === 200) {
          }

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.canViewNewDevices = false;
          this.devices = null;
          this.getDevices();
        }
      });
  }

  onRowClick(index: number) {
    const selectedDevice = this.devices[index];
    var sn = selectedDevice[2];
    return (window.location.href = '/home/devices/view-devices/' + sn);
  }
}
