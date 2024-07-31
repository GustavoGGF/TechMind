import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { UtilitiesModule } from '../../utilities/utilities.module';
import { CommonModule } from '@angular/common';
import {
  HttpClient,
  HttpClientModule,
  HttpHeaders,
} from '@angular/common/http';
import { catchError, throwError } from 'rxjs';

@Component({
  selector: 'app-devices-details',
  standalone: true,
  imports: [UtilitiesModule, CommonModule, HttpClientModule],
  templateUrl: './devices-details.component.html',
  styleUrl: './devices-details.component.css',
})
export class DevicesDetailsComponent implements OnInit {
  constructor(private route: ActivatedRoute, private http: HttpClient) {}
  sn: string = '';
  name: any;
  canView: boolean = false;
  model_device: string = '';
  status: any;
  device: string[] = [];
  equip: string = '';
  brand: string = '';
  imob: string = '';
  img_url: string = '';
  img_config: string = '/static/assets/images/devices/configuracao.png';
  img_close: string = '/static/assets/images/fechar.png';
  canViewDetails: boolean = true;
  canViewModify: boolean = false;
  input_brand: string = '';
  input_model: string = '';
  canViewStatus: boolean = false;
  select_value: string = '';
  machines: string[] = [];
  machines_Names: string[] = [];
  token: string = '';
  linked: string = '';
  mac: string = '';
  selectedDevice: any;
  url_linked_device: string = '';

  ngOnInit() {
    this.name = localStorage.getItem('name');

    if (this.name.length == 0 || this.name == null) {
    } else {
      this.canView = true;
    }
    this.route.params.subscribe((params) => {
      this.sn = params['sn'];
    });
    this.getToken();
    this.getInfoDevice();
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
      });
  }

  getInfoDevice(): void {
    this.http
      .get('/home/devices/info-device/' + this.sn, {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.device = data.data[0];
          console.log(this.device);

          this.equip = this.device[0];
          this.model_device = this.device[1];
          switch (this.model_device) {
            default:
              this.img_url = '';
              break;
            case '1908FPt':
              this.img_url = '/static/assets/images/devices/13803086706.png';
              break;
            case 'FLATRON E2360V-PN':
              this.img_url = '/static/assets/images/devices/medium0.png';
              break;
            case 'SoundPoint Ip 320 SIP':
              this.img_url =
                '/static/assets/images/devices/91+5tSdU8oL._AC_SL1500_.png';
              break;
            case 'ix4-200d':
              this.img_url = '/static/assets/images/devices/i_8.png';
              break;
          }
          this.brand = this.device[4];
          this.imob = this.device[3];
          this.linked = this.device[5];
          var link = this.linked.replace(/:/g, '-');
          this.url_linked_device = '/home/computers/view-machine/' + link;
        }
      });
  }

  modifyDevice(): void {
    this.canViewDetails = false;
    this.canViewModify = true;
  }

  getBrand(event: any): void {
    this.input_brand = event.target.value;
  }

  getModel(event: any): void {
    this.input_model = event.target.value;
  }

  closeModify(): void {
    this.canViewDetails = true;
    this.canViewModify = false;
  }

  getStatus(event: any): void {
    this.select_value = event.target.value;

    switch (this.select_value) {
      default:
        this.machines = [];
        break;
      case 'None':
        this.machines = [];
        break;
      case 'InUse':
        this.http
          .get('/home/devices/get-last-machines', {})
          .pipe(
            catchError((error) => {
              this.status = error.status;

              return throwError(error);
            })
          )
          .subscribe((data: any) => {
            if (data) {
              this.machines = data.machines;

              this.machines.forEach((subArray) => {
                this.machines_Names.push(subArray[1]);
              });

              this.canViewStatus = true;
            }
          });
        break;
    }
  }

  onRowClick(index: number) {
    this.selectedDevice = this.machines[index];

    this.mac = this.selectedDevice[0];
    var device = this.sn;

    this.http
      .post(
        '/home/computers/added-device',
        {
          device: device,
          computer: this.mac,
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

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.device = [];
          this.equip = '';
          this.model_device = '';
          this.img_url = '';
          this.brand = '';
          this.imob = '';
          this.canViewDetails = true;
          this.canViewModify = false;
          this.getInfoDevice();
        }
      });
  }
}
