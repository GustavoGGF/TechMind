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
  // Variaveis Array
  devices: string[] = [];
  divs: string[] = [];
  softwares: string[] = [];

  // Variaveis String
  audio_device_model: string = '';
  audio_device_product: string = '';
  bios_version: string = '';
  cpu_architecture: string = '';
  cpu_core: string = '';
  cpu_max_mhz: string = '';
  cpu_min_mhz: string = '';
  cpu_model_name: string = '';
  cpu_operation_mode: string = '';
  cpu_thread: string = '';
  cpu_socket: string = '';
  cpu_vendor_id: string = '';
  cpus: string = '';
  currentUser: string = '';
  computers_class: string = 'active';
  device_class: string = '';
  domain: string = '';
  gpu_bus_info: string = '';
  gpu_clock: string = '';
  gpu_configuration: string = '';
  gpu_logical_name: string = '';
  gpu_product: string = '';
  gpu_vendor_id: string = '';
  hard_disk_model: string = '';
  hard_disk_sata_version: string = '';
  hard_disk_serial_number: string = '';
  hard_disk_user_capacity: string = '';
  home_class: string = '';
  img_config: string = '/static/assets/images/devices/configuracao.png';
  imob: string = '';
  input_imob: string = '';
  input_note: string = '';
  ip: string = '';
  list_softwares: string = '';
  location: string = '';
  macAddress: string = '';
  manufacturer: string = '';
  max_capacity_memory: string = '';
  model: string = '';
  motherboard_asset_tag: string = '';
  motherboard_manufacturer: string = '';
  motherboard_serial_name: string = '';
  motherboard_product_name: string = '';
  motherboard_version: string = '';
  name_pc: string = '';
  note: string = '';
  number_of_slot: string = '';
  operational_System: string = '';
  present: string = '';
  select_value: string = '';
  serial_number: string = '';
  system_version: string = '';
  url_logo: string = '';
  url_manufacturer: string = '';
  url_model: string = '';
  urlResize = '/static/assets/images/expandir-setas.png';

  // Varaiveis any
  info_PC: any;
  memories: any;
  name: any;
  status: any;
  token: any;

  // Variaveis Boolean
  canView: boolean = false;
  canViewDataAdmin: boolean = true;
  canViewDevices: boolean = false;
  canViewHardWare: boolean = false;
  canViewOthers: boolean = false;
  canViewSoftWare: boolean = false;
  memory_windows: boolean = false;
  modifyOther: boolean = false;
  showBar: boolean = false;

  // Variaveis Object
  softwareList: { name: string; version: string; vendor: string }[] = [];

  // Setando função para verificar o click na pagina
  ngAfterViewInit() {
    // Adiciona o event listener global para cliques no documento
    this.clickListener = this.renderer.listen(
      'document',
      'click',
      this.handleClick.bind(this)
    );
  }

  // Caso o click não seja na aba de guia ou no botão de resize, a barra de informações é escondida
  handleClick(event: MouseEvent): void {
    const target = event.target as HTMLElement;

    if (target) {
      if (
        target.id !== 'nvbar' &&
        target.id !== 'resize' &&
        target.id !== 'hard_data'
      ) {
        this.showBar = false;
      }
    }
  }

  // Função inicia ao entrar na pagina
  ngOnInit(): void {
    // Pegando o token CSRF
    this.getToken();
    // Pegando os dados do usuario
    this.name = localStorage.getItem('name');

    // Verificando se o nome foi obtido
    if (this.name.length == 0 || this.name == null) {
    } else {
      this.canView = true;
    }

    // Pegando o mac_address
    this.route.params.subscribe((params) => {
      this.macAddress = params['mac'];
    });

    // Obtendo dados do equipamento
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
          // Selecionando a logo do sistema operacional
          let operational_System_string = this.operational_System
            .toLowerCase()
            .replace(/\s+/g, '');
          switch (operational_System_string) {
            default:
              this.url_logo = '';
              break;
            case 'ubuntu':
              this.url_logo = '/static/assets/images/brands/ubuntu.png';
              break;
            case 'windows10':
              this.url_logo = '/static/assets/images/brands/windows10.png';
              break;
          }
          this.system_version = this.info_PC[6];
          this.domain = this.info_PC[7];
          this.ip = this.info_PC[8];
          this.manufacturer = this.info_PC[9];
          // Selecionando a logo da Marca do equipamento
          let manufacturer_string = this.manufacturer
            .toLowerCase()
            .replace(/\s+/g, '');
          switch (manufacturer_string) {
            default:
              this.url_manufacturer = '';
              break;
            case 'dellinc.':
              this.url_manufacturer = '/static/assets/images/brands/dell.png';
              break;
            case 'vmware,Inc.':
              this.url_manufacturer = '/static/assets/images/brands/VMware.png';
              break;
            case 'hp':
              this.url_manufacturer =
                '/static/assets/images/brands/100px-HP_logo_2012.svg.png';
              break;
            case 'hewlett-packard':
              this.url_manufacturer =
                '/static/assets/images/brands/100px-HP_logo_2012.svg.png';
              break;
            case 'pcware':
              this.url_manufacturer =
                '/static/assets/images/brands/logo-pcware-.png';
              break;
            case 'lenovo':
              this.url_manufacturer =
                '/static/assets/images/brands/logo-lenovo-1024x295.png';
              break;
          }
          this.model = this.info_PC[10];
          // Selecionando a imagem do equipamento
          let model_string = this.model.toLowerCase().replace(/\s+/g, '');
          console.log(model_string);

          switch (model_string) {
            default:
              this.url_model = '';
              break;
            case 'precisionm4600':
              this.url_model = '/static/assets/images/models/_OVR.webp';
              break;
            case 'vmwarevirtualplatform':
              this.url_model = '/static/assets/images/brands/feature-image.png';
              break;
            case 'optiplexgx620':
              this.url_model =
                '/static/assets/images/models/81p7NifF3RL._AC_SL1500_.jpg';
              break;
            case 'precisionworkstationt3400':
              this.url_model = '/static/assets/images/models/OriginalJPG-.png';
              break;
            case 'hpcompaqelite8300sff':
              this.url_model = '/static/assets/images/models/c02753259.jpg';
              break;
            case 'ipmh310gpro':
              this.url_model = '/static/assets/images/models/IPMH310G_PRO.png';
              break;
            case 'hpprodesk400g4sff':
              this.url_model = '/static/assets/images/models/c05924778.png';
              break;
            case '':
              this.url_model =
                '/static/assets/images/models/24165952843_LenovoV14Gen3ABABlackforTextureIMG_202201050201591696357227289.png';
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
          // Verificando se o SMBIOS está presente
          if (this.info_PC[36].includes('present')) {
            this.present = 'Present';
          } else {
            this.present = 'Not found';
          }

          // Ajustando versão da bios
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

          // Ajustando a lsita de softwares
          let softwares_list = this.info_PC[42];
          if (softwares_list) {
            if (
              this.operational_System == 'Windows10' ||
              this.operational_System == 'Windows8.1'
            ) {
              let jsonString = softwares_list.replace(/'/g, '"');
              this.softwareList = JSON.parse(jsonString);
            } else {
              let names = softwares_list.split(',');

              for (let i = 0; i < names.length; i++) {
                this.softwares.push(names[i]);
              }
            }
          }

          // Ajustando as memorias
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

  // FUnção que obtem o token CSRF
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

  // Função que expande e retrai as abas
  async resizeBar(): Promise<void> {
    if (this.showBar) {
      this.showBar = false;
    } else {
      this.showBar = true;
    }
  }

  // Função que mostra a aba de HardWare
  showHardware(): void {
    this.canViewDataAdmin = false;
    this.canViewHardWare = true;
    this.canViewSoftWare = false;
    this.canViewDevices = false;
    this.canViewOthers = false;
  }

  // Função que mostra a aba de Dados Administrativos
  showDataAdmin(): void {
    this.canViewDataAdmin = true;
    this.canViewHardWare = false;
    this.canViewSoftWare = false;
    this.canViewDevices = false;
    this.canViewOthers = false;
  }

  // Função que mostra os Softwares
  showSoftWare(): void {
    this.canViewDataAdmin = false;
    this.canViewHardWare = false;
    this.canViewSoftWare = true;
    this.canViewDevices = false;
    this.canViewOthers = false;
  }

  // Função que mostra os dispositivos atrelado ao equipamento
  showDevices(): void {
    this.canViewDataAdmin = false;
    this.canViewHardWare = false;
    this.canViewSoftWare = false;
    this.canViewDevices = true;
    this.canViewOthers = false;
    var mac = this.macAddress.replace(/-/g, '');
    // Obtem os dispositivos atrelados
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

  // Função que mostra a ba outros
  showOthers(): void {
    this.canViewDataAdmin = false;
    this.canViewHardWare = false;
    this.canViewSoftWare = false;
    this.canViewDevices = false;
    this.canViewOthers = true;
  }

  // Função que manda para a url do dispositivos selecionado
  onRowClick(index: number) {
    const selectedDevice = this.devices[index];
    var sn = selectedDevice[2];
    return (window.location.href = '/home/devices/view-devices/' + sn);
  }

  // Função que converte bytes em GB
  convertBytesToGB(bytes: number): number {
    return bytes / 1024 ** 3;
  }

  // Função que converte bytes em GB
  convertBytesToGB2(capacity: string): number {
    // Remove qualquer texto adicional e converte para número
    const numericValue = parseFloat(capacity.replace(/[^0-9]/g, ''));
    // Assumindo que o valor é em GB
    return numericValue;
  }

  // FUnção que libera a modificação na aba outros
  modifyDevice(): void {
    this.modifyOther = true;
  }

  // Função que obtem o valor do imob
  getImob(event: any): void {
    this.input_imob = event.target.value;
  }

  // Função que obtem o valor da localização selecionado
  device_select(event: any): void {
    this.select_value = event.target.value;
  }

  //Função que obtem as observações
  getNote(event: any): void {
    this.input_note = event.target.value;
  }

  // Função que salva os dados modificados e ja mostra eles atualizados na tela
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
          } else if (data.note) {
            this.note = data.note;
          }
        }
      });

    this.modifyOther = false;
  }
}
