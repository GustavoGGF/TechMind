import {
  HttpClient,
  HttpClientModule,
  HttpHeaders,
} from '@angular/common/http';
import {
  AfterViewInit,
  Component,
  ElementRef,
  OnInit,
  Renderer2,
  ViewChild,
} from '@angular/core';
import { catchError, throwError } from 'rxjs';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { UtilitiesModule } from '../../utilities/utilities.module';

@Component({
  selector: 'app-computers-details',
  standalone: true,
  imports: [CommonModule, HttpClientModule, UtilitiesModule],
  templateUrl: './computers-details.component.html',
  styleUrl: './computers-details.component.css',
})
export class ComputersDetailsComponent implements OnInit, AfterViewInit {
  private clickListener: (() => void) | undefined;

  constructor(
    private route: ActivatedRoute,
    private http: HttpClient,
    private renderer: Renderer2
  ) {}

  @ViewChild('main') main!: ElementRef;

  divs: string[] = [];

  currentUser: string = '';
  macAddress: string = '';
  name: any;
  name_pc: string = '';
  operational_System: string = '';
  system_version: string = '';
  ip: string = '';
  url_logo: string = '';
  manufacturer: string = '';
  url_manufacturer: string = '';
  model: string = '';
  url_model: string = '';
  canView: boolean = false;
  canViewHardWare: boolean = false;
  canViewDataAdmin: boolean = true;
  serial_number: string = '';
  max_capacity_memory: string = '';
  number_of_slot: string = '';
  hard_disk_model = '';
  hard_disk_serial_number = '';
  hard_disk_user_capacity = '';
  hard_disk_sata_version = '';
  cpu_architecture = '';
  cpu_operation_mode = '';
  cpus = '';
  cpu_vendor_id = '';
  cpu_model_name = '';
  cpu_thread = '';
  cpu_core = '';
  cpu_socket = '';
  cpu_max_mhz = '';
  cpu_min_mhz = '';
  gpu_product = '';
  gpu_vendor_id = '';
  gpu_bus_info = '';
  gpu_logical_name = '';
  gpu_clock = '';
  gpu_configuration = '';
  audio_device_product = '';
  audio_device_model = '';
  bios_version = '';
  present: string = '';
  motherboard_manufacturer: string = '';
  motherboard_product_name: string = '';
  motherboard_version: string = '';
  motherboard_serial_name: string = '';
  motherboard_asset_tag: string = '';
  canViewSoftWare: boolean = false;
  list_softwares: string = '';
  softwares: string[] = [];
  canViewDevices: boolean = false;
  devices: string[] = [];
  domain: string = '';
  memories: any;
  memory_windows: boolean = false;
  softwareList: { name: string; version: string; vendor: string }[] = [];
  canViewOthers: boolean = false;
  imob: string = '';
  img_config: string = '/static/assets/images/devices/configuracao.png';
  modifyOther: boolean = false;
  input_imob: string = '';
  token: any;
  location: string = '';
  select_value: string = '';
  note: string = '';
  input_note: string = '';

  urlResize = '/static/assets/images/expandir-setas.png';

  info_PC: any;
  status: any;

  showBar: boolean = false;

  ngAfterViewInit() {
    // Adiciona o event listener global para cliques no documento
    this.clickListener = this.renderer.listen(
      'document',
      'click',
      this.handleClick.bind(this)
    );
  }

  handleClick(event: MouseEvent): void {
    const target = event.target as HTMLElement;

    if (target) {
      // Verifique o id do elemento clicado e execute a lógica desejada
      if (
        target.id !== 'nvbar' &&
        target.id !== 'resize' &&
        target.id !== 'hard_data'
      ) {
        this.showBar = false;
      }
    }
  }

  ngOnInit(): void {
    this.getToken();
    this.name = localStorage.getItem('name');

    if (this.name.length == 0 || this.name == null) {
    } else {
      this.canView = true;
    }
    this.route.params.subscribe((params) => {
      this.macAddress = params['mac'];
    });

    this.http
      .get('/home/computers/info-machine/' + this.macAddress, {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.info_PC = data.data[0];
          this.name_pc = this.info_PC[1];
          this.currentUser = this.info_PC[5];
          this.operational_System = this.info_PC[3];
          switch (this.operational_System) {
            default:
              this.url_logo = '';
              break;
            case 'ubuntu':
              this.url_logo = '/static/assets/images/brands/ubuntu.png';
              break;
            case 'Windows10':
              this.url_logo = '/static/assets/images/brands/windows10.png';
              break;
          }
          this.system_version = this.info_PC[6];
          this.domain = this.info_PC[7];
          this.ip = this.info_PC[8];
          this.manufacturer = this.info_PC[9];
          switch (this.manufacturer) {
            default:
              this.url_manufacturer = '';
              break;
            case 'Dell Inc.':
              this.url_manufacturer = '/static/assets/images/brands/dell.png';
              break;
            case 'VMware, Inc.':
              this.url_manufacturer = '/static/assets/images/brands/VMware.png';
              break;
          }
          this.model = this.info_PC[10];
          switch (this.model) {
            default:
              this.url_model = '';
              break;
            case 'Precision M4600':
              this.url_model = '/static/assets/images/models/_OVR.webp';
              break;
            case 'VMware Virtual Platform':
              this.url_model = '/static/assets/images/brands/feature-image.png';
              break;
            case 'OptiPlex GX620':
              this.url_model =
                '/static/assets/images/models/81p7NifF3RL._AC_SL1500_.jpg';
              break;
            case 'Precision WorkStation T3400 ':
              this.url_model = '/static/assets/images/models/OriginalJPG-.png';
              break;
          }
          this.serial_number = this.info_PC[11];
          this.max_capacity_memory = this.info_PC[12];
          this.number_of_slot = this.info_PC[13];
          this.hard_disk_model = this.info_PC[14];
          this.hard_disk_serial_number = this.info_PC[15];
          this.hard_disk_user_capacity = this.info_PC[16];
          this.hard_disk_sata_version = this.info_PC[17];
          this.cpu_architecture = this.info_PC[18];
          this.cpu_operation_mode = this.info_PC[19];
          this.cpus = this.info_PC[20];
          this.cpu_vendor_id = this.info_PC[21];
          this.cpu_model_name = this.info_PC[22];
          this.cpu_thread = this.info_PC[23];
          this.cpu_socket = this.info_PC[24];
          this.cpu_max_mhz = this.info_PC[25];
          this.cpu_min_mhz = this.info_PC[26];
          this.cpu_core = this.info_PC[27];
          this.gpu_product = this.info_PC[28];
          this.gpu_vendor_id = this.info_PC[29];
          this.gpu_bus_info = this.info_PC[30];
          this.gpu_logical_name = this.info_PC[31];
          this.gpu_clock = this.info_PC[32];
          this.gpu_configuration = this.info_PC[33];
          this.audio_device_product = this.info_PC[34];
          this.audio_device_model = this.info_PC[35];
          if (this.info_PC[36].includes('present')) {
            this.present = 'Present';
          } else {
            this.present = 'Not found';
          }
          let regex = /(.{2})\.(.{2})/;
          let matches = this.bios_version.match(regex);
          if (matches) {
            let part_1 = matches[1];
            let part_2 = matches[2];
            this.bios_version = part_1 + '.' + part_2;
          }
          this.motherboard_manufacturer = this.info_PC[37];
          this.motherboard_product_name = this.info_PC[38];
          this.motherboard_version = this.info_PC[39];
          this.motherboard_serial_name = this.info_PC[40];
          this.motherboard_asset_tag = this.info_PC[41];
          let softwares_list = this.info_PC[42];
          if (softwares_list) {
            if (this.operational_System == 'Windows10') {
              let jsonString = softwares_list.replace(/'/g, '"');
              this.softwareList = JSON.parse(jsonString);
            } else {
              let names = softwares_list.split(',');

              for (let i = 0; i < names.length; i++) {
                this.softwares.push(names[i]);
              }
            }
          }
          this.memories = this.info_PC[43];
          if (this.memories) {
            let valid = this.memories.replace(/'/g, '"');

            this.memories = JSON.parse(valid);

            switch (this.operational_System) {
              default:
                break;
              case 'Windows10':
                this.memory_windows = true;
            }
          }
          this.imob = this.info_PC[44];
          this.location = this.info_PC[45];
          this.note = this.info_PC[46];
        }
      });
  }

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

  async resizeBar(): Promise<void> {
    if (this.showBar) {
      this.showBar = false;
    } else {
      this.showBar = true;
    }
  }

  showHardware(): void {
    this.canViewDataAdmin = false;
    this.canViewHardWare = true;
    this.canViewSoftWare = false;
    this.canViewDevices = false;
    this.canViewOthers = false;
  }

  showDataAdmin(): void {
    this.canViewDataAdmin = true;
    this.canViewHardWare = false;
    this.canViewSoftWare = false;
    this.canViewDevices = false;
    this.canViewOthers = false;
  }

  showSoftWare(): void {
    this.canViewDataAdmin = false;
    this.canViewHardWare = false;
    this.canViewSoftWare = true;
    this.canViewDevices = false;
    this.canViewOthers = false;
  }

  showDevices(): void {
    this.canViewDataAdmin = false;
    this.canViewHardWare = false;
    this.canViewSoftWare = false;
    this.canViewDevices = true;
    this.canViewOthers = false;
    var mac = this.macAddress.replace(/-/g, '');
    this.http
      .get('/home/computers/added-devices/' + mac, {})
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
          this.devices = data.data;
        }
      });
  }

  showOthers(): void {
    this.canViewDataAdmin = false;
    this.canViewHardWare = false;
    this.canViewSoftWare = false;
    this.canViewDevices = false;
    this.canViewOthers = true;
  }

  onRowClick(index: number) {
    const selectedDevice = this.devices[index];
    var sn = selectedDevice[2];
    return (window.location.href = '/home/devices/view-devices/' + sn);
  }

  convertBytesToGB(bytes: number): number {
    return bytes / 1024 ** 3;
  }

  convertBytesToGB2(capacity: string): number {
    // Remove qualquer texto adicional e converte para número
    const numericValue = parseFloat(capacity.replace(/[^0-9]/g, ''));
    // Assumindo que o valor é em GB
    return numericValue;
  }

  modifyDevice(): void {
    this.modifyOther = true;
  }

  getImob(event: any): void {
    this.input_imob = event.target.value;
  }

  device_select(event: any): void {
    this.select_value = event.target.value;
  }

  getNote(event: any): void {
    this.input_note = event.target.value;
  }

  submitOthers(): void {
    var mac = this.macAddress.replace(/-/g, '');
    this.http
      .post(
        '/home/computers/modify-others/' + mac,
        {
          imob: this.input_imob,
          location: this.select_value,
          note: this.input_note,
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

          if (this.status === 0) {
          }

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          if (data.imob) {
            this.imob = data.imob;
          } else if (data.location) {
            this.location = data.location;
          }
        }
      });

    this.modifyOther = false;
  }
}
