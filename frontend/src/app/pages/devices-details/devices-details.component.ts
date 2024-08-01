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

  // Variaveis String
  brand: string = '';
  computers_class: string = '';
  device_class: string = 'active';
  equip: string = '';
  home_class: string = '';
  img_close: string = '/static/assets/images/fechar.png';
  img_config: string = '/static/assets/images/devices/configuracao.png';
  img_url: string = '';
  imob: string = '';
  input_brand: string = '';
  input_model: string = '';
  linked: string = '';
  mac: string = '';
  model_device: string = '';
  select_value: string = '';
  sn: string = '';
  token: string = '';
  url_linked_device: string = '';

  // Variaveis Any
  name: any;
  selectedDevice: any;
  status: any;

  // Varaivei Boolean
  canView: boolean = false;
  canViewDetails: boolean = true;
  canViewModify: boolean = false;
  canViewStatus: boolean = false;

  // Varaiveis Array
  device: string[] = [];
  machines: string[] = [];
  machines_Names: string[] = [];

  // Função inicia assim que a pagina carrega
  ngOnInit() {
    // Pega os dados do usuario
    this.name = localStorage.getItem('name');

    // Verificar se os dados existem
    if (this.name.length == 0 || this.name == null) {
      this.canView = false;
    } else {
      this.canView = true;
    }

    // Pega o SN do dispositivo
    this.route.params.subscribe((params) => {
      this.sn = params['sn'];
    });
    // Inicia a função para pegar o token CSRF
    this.getToken();
    // Inica a função para pegar dados dos dispositivos
    this.getInfoDevice();
  }

  // Função para pegar o Token CSRF
  getToken(): void {
    this.http
      .get('/home/get-token', {})
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

  // Função que pega as informações dos dispositivos
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
          this.equip = this.device[0];
          this.model_device = this.device[1];
          let model_string = this.model_device
            .toLowerCase()
            .replace(/\s+/g, '');
          // Define a imagem do dispositivo
          switch (model_string) {
            default:
              this.img_url = '';
              break;
            case '1908fpt':
              this.img_url = '/static/assets/images/devices/13803086706.png';
              break;
            case 'flatrone2360v-pn':
              this.img_url = '/static/assets/images/devices/medium0.png';
              break;
            case 'soundpointip320sip':
              this.img_url =
                '/static/assets/images/devices/91+5tSdU8oL._AC_SL1500_.png';
              break;
            case 'soundpointip331':
              this.img_url =
                '/static/assets/images/devices/91+5tSdU8oL._AC_SL1500_.png';
              break;
            case 'ix4-200d':
              this.img_url = '/static/assets/images/devices/1_8.png';
              break;
            case 'sg300-52':
              this.img_url =
                '/static/assets/images/devices/cisco-sg300-52mp-k9-eu-small-business-sg300-52mp-switch-10126213.png';
              break;
            case '2007FPb':
              this.img_url = '/static/assets/images/devices/3282872674_1.png';
              break;
            case 'e157fpc':
              this.img_url = '/static/assets/images/devices/3282872674_1.png';
              break;
            case '3580l33/l3h':
              this.img_url = '/static/assets/images/devices/img_1042.png';
              break;
            case 'i5056':
              this.img_url = '/static/assets/images/devices/152247512_1.png';
              break;
            case 'e170sc':
              this.img_url =
                '/static/assets/images/devices/monitor_lcd_dell_e170sc_17_polegadas_8617_1_fae8340e93d667349f89992e883b59c2.png';
              break;
            case 'syncmastert190':
              this.img_url =
                '/static/assets/images/devices/usado_monitor_samsung_syncmaster_19_polegadas_wide_t190_339_1_16c02a2372deadc946cfc58aa9c4b942.png';
              break;
            case 'hpl1906':
              this.img_url =
                '/static/assets/images/devices/monitor_19_hp_l1906_vga_1280_x_1024px_quadrado_11745_1_1f65dcf7df12237460d02b7b51d372a7.png';
              break;
            case 'hpcompaqle1711':
              this.img_url = '/static/assets/images/devices/c02738428.png';
              break;
            case 'hple1711':
              this.img_url = '/static/assets/images/devices/c02738428.png';
              break;
            case 'e178fpc':
              this.img_url =
                '/static/assets/images/devices/monitor_dell_e178fpc_17_polegadas_8613_1_46f514f1de0f425b7270b4f6435780a1.png';
              break;
            case '740nc(r)':
              this.img_url =
                '/static/assets/images/devices/D_NQ_NP_691506-MLB74041561731_012024-O.png';
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

  // Libera a tela para modificação de dispositivo
  modifyDevice(): void {
    this.canViewDetails = false;
    this.canViewModify = true;
  }

  // Função que pega a Marca do dispositivo
  getBrand(event: any): void {
    this.input_brand = event.target.value;
  }

  // Função que pega o modelo do dispositivo
  getModel(event: any): void {
    this.input_model = event.target.value;
  }

  // Função que fecha a tela de modificação
  closeModify(): void {
    this.canViewDetails = true;
    this.canViewModify = false;
  }

  // Função que pega o Status selecionado
  getStatus(event: any): void {
    this.select_value = event.target.value;

    switch (this.select_value) {
      default:
        this.machines = [];
        break;
      case 'None':
        this.machines = [];
        break;
      // Caso o status seja 'InUse' ele pega os ultimos 10 computadores
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

  // Função que linka o dispositivo com um computador
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
