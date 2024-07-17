import { HttpClient, HttpClientModule } from '@angular/common/http';
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
  first_slot_dim: string = '';
  second_slot_dim: string = '';
  third_slot_dim: string = '';
  fourth_slot_dim: string = '';
  first_size: string = '';
  second_size: string = '';
  third_size: string = '';
  fourth_size: string = '';
  first_type: string = '';
  second_type: string = '';
  third_type: string = '';
  fourth_type: string = '';
  first_type_details = '';
  second_type_details = '';
  third_type_details = '';
  fourth_type_details = '';
  first_speed_memory = '';
  second_speed_memory = '';
  third_speed_memory = '';
  fourth_speed_memory = '';
  first_serial_number = '';
  second_serial_number = '';
  third_serial_number = '';
  fourth_serial_number = '';
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
              break;
            case 'ubuntu':
              this.url_logo = '/static/assets/images/ubuntu.png';
          }
          this.system_version = this.info_PC[6];
          this.ip = this.info_PC[8];
          this.manufacturer = this.info_PC[9];
          switch (this.manufacturer) {
            default:
              break;
            case 'Dell Inc.':
              this.url_manufacturer = '/static/assets/images/dell.png';
          }
          this.model = this.info_PC[10];
          switch (this.model) {
            default:
              break;
            case 'Precision M4600':
              this.url_model = '/static/assets/images/models/_OVR.webp';
          }
          this.serial_number = this.info_PC[11];
          this.max_capacity_memory = this.info_PC[12];
          this.number_of_slot = this.info_PC[13];
          this.first_slot_dim = this.info_PC[14];
          this.second_slot_dim = this.info_PC[15];
          this.third_slot_dim = this.info_PC[16];
          this.fourth_slot_dim = this.info_PC[17];
          this.first_size = this.info_PC[18];
          this.second_size = this.info_PC[19];
          this.third_size = this.info_PC[20];
          this.fourth_size = this.info_PC[21];
          this.first_type = this.info_PC[22];
          this.second_type = this.info_PC[23];
          this.third_type = this.info_PC[24];
          this.fourth_type = this.info_PC[25];
          this.first_type_details = this.info_PC[26];
          this.second_type_details = this.info_PC[27];
          this.third_type_details = this.info_PC[28];
          this.fourth_type_details = this.info_PC[29];
          this.first_speed_memory = this.info_PC[30];
          this.second_speed_memory = this.info_PC[31];
          this.third_speed_memory = this.info_PC[32];
          this.fourth_speed_memory = this.info_PC[33];
          this.first_serial_number = this.info_PC[34];
          this.second_serial_number = this.info_PC[35];
          this.third_serial_number = this.info_PC[36];
          this.fourth_serial_number = this.info_PC[37];
          this.hard_disk_model = this.info_PC[38];
          this.hard_disk_serial_number = this.info_PC[39];
          this.hard_disk_user_capacity = this.info_PC[40];
          this.hard_disk_sata_version = this.info_PC[41];
          this.cpu_architecture = this.info_PC[42];
          this.cpu_operation_mode = this.info_PC[43];
          this.cpus = this.info_PC[44];
          this.cpu_vendor_id = this.info_PC[45];
          this.cpu_model_name = this.info_PC[46];
          this.cpu_thread = this.info_PC[47];
          this.cpu_socket = this.info_PC[48];
          this.cpu_max_mhz = this.info_PC[49];
          this.cpu_min_mhz = this.info_PC[50];
          this.cpu_core = this.info_PC[51];
          this.gpu_product = this.info_PC[52];
          this.gpu_vendor_id = this.info_PC[53];
          this.gpu_bus_info = this.info_PC[54];
          this.gpu_logical_name = this.info_PC[55];
          this.gpu_clock = this.info_PC[56];
          this.gpu_configuration = this.info_PC[57];
          this.audio_device_product = this.info_PC[58];
          this.audio_device_model = this.info_PC[59];
          if (this.info_PC[60].includes('present')) {
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
          this.motherboard_manufacturer = this.info_PC[61];
          this.motherboard_product_name = this.info_PC[62];
          this.motherboard_version = this.info_PC[63];
          this.motherboard_serial_name = this.info_PC[64];
          this.motherboard_asset_tag = this.info_PC[65];
          let softwares_list = this.info_PC[66];
          let names = softwares_list.split(',');

          for (let i = 0; i < names.length; i++) {
            this.softwares.push(names[i]);
          }
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
  }

  showDataAdmin(): void {
    this.canViewDataAdmin = true;
    this.canViewHardWare = false;
    this.canViewSoftWare = false;
    this.canViewDevices = false;
  }

  showSoftWare(): void {
    this.canViewDataAdmin = false;
    this.canViewHardWare = false;
    this.canViewSoftWare = true;
    this.canViewDevices = false;
  }

  showDevices(): void {
    this.canViewDataAdmin = false;
    this.canViewHardWare = false;
    this.canViewSoftWare = false;
    this.canViewDevices = true;
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

  onRowClick(index: number) {
    const selectedDevice = this.devices[index];
    var sn = selectedDevice[2];
    return (window.location.href = '/home/devices/view-devices/' + sn);
  }
}
