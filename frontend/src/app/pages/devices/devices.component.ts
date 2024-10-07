import { Component } from '@angular/core';
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
  all_quantity: string = '';
  arrow_down: string = '/static/assets/images/seta.png';
  arrow_up: string = '/static/assets/images/seta2.png';
  close_button: string = '/static/assets/images/fechar.png';
  computers_class: string = '';
  device_class: string = 'active';
  fifty_quantity: string = '';
  input_brand: string = '';
  input_imob: string = '';
  input_model: string = '';
  input_serial_number: string = '';
  imageDevices: string = '/static/assets/images/devices/computador.png';
  home_class: string = '';
  select_value: string = '';
  ten_quantity: string = '';
  one_hundred_quantity: string = '';
  quantity_devices: string | null = '';

  // Função para pegar o nome do usuario logado
  ngOnInit(): void {
    this.name = localStorage.getItem('name');

    this.quantity_devices = localStorage.getItem('quantity_devices');
    if (this.quantity_devices == null) {
      localStorage.setItem('quantity_devices', '10');
      this.quantity_devices = '10';
    }

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
        this.getDevices();
      });
  }

  // Função para pegar os dispositivos
  getDevices(): void {
    this.http
      .get('/home/devices/get-devices/' + this.quantity_devices, {})
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
          this.devices = data.devices;
          this.canViewDevices = true;
          switch (this.quantity_devices) {
            default:
              break;
            case '10':
              this.ten_quantity = 'active_filter';
              this.fifty_quantity = '';
              this.one_hundred_quantity = '';
              this.all_quantity = '';
              break;
            case '50':
              this.ten_quantity = '';
              this.fifty_quantity = 'active_filter';
              this.one_hundred_quantity = '';
              this.all_quantity = '';
              break;
            case '100':
              this.ten_quantity = '';
              this.fifty_quantity = '';
              this.one_hundred_quantity = 'active_filter';
              this.all_quantity = '';
              break;
            case 'all':
              this.ten_quantity = '';
              this.fifty_quantity = '';
              this.one_hundred_quantity = '';
              this.all_quantity = 'active_filter';
          }
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

  // Seta quantidade de maquinas a serem exibidas para 10
  getTen(): void {
    localStorage.setItem('quantity_devices', '10');

    this.quantity_devices = '10';

    this.ten_quantity = 'active_filter';
    this.fifty_quantity = '';
    this.one_hundred_quantity = '';
    this.all_quantity = '';

    this.canViewDevices = false;
    this.devices = null;

    this.http
      .get('/home/devices/get-devices/10', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.devices = data.devices;

          this.canViewDevices = true;
        }
      });
  }

  // Seta quantidade de maquinas a serem exibidas para 10
  getFifty(): void {
    localStorage.setItem('quantity_devices', '50');

    this.quantity_devices = '50';

    this.ten_quantity = '';
    this.fifty_quantity = 'active_filter';
    this.one_hundred_quantity = '';
    this.all_quantity = '';

    this.canViewDevices = false;
    this.devices = null;

    this.http
      .get('/home/devices/get-devices/50', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.devices = data.devices;

          this.canViewDevices = true;
        }
      });
  }

  getOneHundred(): void {
    localStorage.setItem('quantity_devices', '100');

    this.quantity_devices = '100';

    this.ten_quantity = '';
    this.fifty_quantity = '';
    this.one_hundred_quantity = 'active_filter';
    this.all_quantity = '';

    this.canViewDevices = false;
    this.devices = null;

    this.http
      .get('/home/devices/get-devices/100', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.devices = data.devices;

          this.canViewDevices = true;
        }
      });
  }

  getAll(): void {
    localStorage.setItem('quantity_devices', 'all');

    this.quantity_devices = '100';

    this.ten_quantity = '';
    this.fifty_quantity = '';
    this.one_hundred_quantity = '';
    this.all_quantity = 'active_filter';

    this.canViewDevices = false;
    this.devices = null;

    this.http
      .get('/home/devices/get-devices/all', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.devices = data.devices;

          this.canViewDevices = true;
        }
      });
  }

  // Reorganiza os dispositivos pelo equipamento em ordem alfabetica
  sortByEquip(): void {
    this.devices.sort((a: any, b: any) => {
      const nameA = a[0].toUpperCase(); // Ignore case
      const nameB = b[0].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }
      return 0;
    });
  }

  // Reorganiza os computadores pelo equipamento em ordem alfabetica invertido
  sortDataByEquipDescending(): void {
    this.devices.sort((a: any, b: any) => {
      const nameA = a[0].toUpperCase(); // Ignore case
      const nameB = b[0].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return 1;
      }
      if (nameA > nameB) {
        return -1;
      }
      return 0;
    });
  }

  // Reorganiza os dispositivos pela Marca em ordem alfabetica
  sortByBrand(): void {
    this.devices.sort((a: any, b: any) => {
      const nameA = a[4].toUpperCase(); // Ignore case
      const nameB = b[4].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }
      return 0;
    });
  }

  // Reorganiza os computadores pela Marca em ordem alfabetica invertido
  sortDataByBrandDescending(): void {
    this.devices.sort((a: any, b: any) => {
      const nameA = a[4].toUpperCase(); // Ignore case
      const nameB = b[4].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return 1;
      }
      if (nameA > nameB) {
        return -1;
      }
      return 0;
    });
  }

  // Reorganiza os dispositivos pelo equipamento em ordem alfabetica
  sortByModel(): void {
    this.devices.sort((a: any, b: any) => {
      const nameA = a[1].toUpperCase(); // Ignore case
      const nameB = b[1].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }
      return 0;
    });
  }

  // Reorganiza os computadores pelo equipamento em ordem alfabetica invertido
  sortDataByModelDescending(): void {
    this.devices.sort((a: any, b: any) => {
      const nameA = a[1].toUpperCase(); // Ignore case
      const nameB = b[1].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return 1;
      }
      if (nameA > nameB) {
        return -1;
      }
      return 0;
    });
  }

  // Reorganiza os dispositivos pelo SN em ordem alfabetica
  sortBySN(): void {
    this.devices.sort((a: any, b: any) => {
      const nameA = a[2].toUpperCase(); // Ignore case
      const nameB = b[2].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }
      return 0;
    });
  }

  // Reorganiza os computadores pelo SN em ordem alfabetica invertido
  sortDataBySNDescending(): void {
    this.devices.sort((a: any, b: any) => {
      const nameA = a[2].toUpperCase(); // Ignore case
      const nameB = b[2].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return 1;
      }
      if (nameA > nameB) {
        return -1;
      }
      return 0;
    });
  }

  // Reorganiza os dispositivos pelo SN em ordem alfabetica
  sortByImob(): void {
    this.devices.sort((a: any, b: any) => {
      const nameA = a[3].toUpperCase(); // Ignore case
      const nameB = b[3].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }
      return 0;
    });
  }

  // Reorganiza os computadores pelo SN em ordem alfabetica invertido
  sortDataByImobDescending(): void {
    this.devices.sort((a: any, b: any) => {
      const nameA = a[3].toUpperCase(); // Ignore case
      const nameB = b[3].toUpperCase(); // Ignore case
      if (nameA < nameB) {
        return 1;
      }
      if (nameA > nameB) {
        return -1;
      }
      return 0;
    });
  }
}
