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
  // Declarando as variaveis any
  devices: any;
  name: any;
  token: any;
  status: any;

  // Declarando as variaveis boolean
  canView: boolean = false;
  canViewDevices: boolean = false;
  canViewNewDevices: boolean = false;

  // Declarando as variaveis string
  close_button: string = '/static/assets/images/fechar.png';
  computers_class: string = '';
  device_class: string = 'active';
  input_brand: string = '';
  input_imob: string = '';
  input_model: string = '';
  input_serial_number: string = '';
  imageDevices: string = '/static/assets/images/devices/computador.png';
  home_class: string = '';
  select_value: string = '';

  // Função para pegar o nome do usuario logado
  ngOnInit(): void {
    this.name = localStorage.getItem('name');

    if (this.name.length == 0 || this.name == null) {
    } else {
      this.canView = true;
      // Inicia a função para pegar o token CSRF
      this.getToken();
    }
  }

  // Função para pegar o token CSRF
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

  // Função para pegar os dispositivos
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

  // Função para liberar a tela de adicionar dispositivos
  addDevices(): void {
    this.canViewNewDevices = true;
  }

  // Função para fechar a tela de adicionar dispositivos
  closeNewDevices(): void {
    this.canViewNewDevices = false;
  }

  // FUnção para pegar o dispositivo selecionado
  device_select(event: any): void {
    this.select_value = event.target.value;
  }

  // Função para pegar o modelo do dispositivo
  getModel(event: any): void {
    this.input_model = event.target.value;
  }

  // Função para pegar o numero de serie do dispositivo
  getSerialNumber(event: any): void {
    this.input_serial_number = event.target.value;
  }

  // Função para pegar a imob do dispositivo
  getImob(event: any): void {
    this.input_imob = event.target.value;
  }

  // Função para pegar a marca do dispositivo
  getBrand(event: any): void {
    this.input_brand = event.target.value;
  }

  // Função para enviar os dados do novo dispositivo
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

  // Função para pegar o index do dispositivo selecionado e direcionar para a url do dispositivo
  onRowClick(index: number) {
    const selectedDevice = this.devices[index];
    var sn = selectedDevice[2];
    return (window.location.href = '/home/devices/view-devices/' + sn);
  }
}
