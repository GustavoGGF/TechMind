import { Component, OnInit } from '@angular/core';
import { UtilitiesModule } from '../../utilities/utilities.module';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { catchError, throwError } from 'rxjs';

@Component({
  selector: 'app-computers',
  standalone: true,
  imports: [UtilitiesModule, CommonModule, HttpClientModule],
  templateUrl: './computers.component.html',
  styleUrl: './computers.component.css',
})
export class ComputersComponent implements OnInit {
  constructor(private http: HttpClient) {}
  // Declarando variaveis any
  dataMachines: any;
  name: any;
  status: any;

  // Declarando variaveis string
  all_quantity: string = '';
  arrow_up: string = '/static/assets/images/seta2.png';
  arrow_down: string = '/static/assets/images/seta.png';
  computers_class: string = 'active';
  device_class: string = '';
  errorType: string = '';
  fifty_quantity: string = '';
  home_class: string = '';
  messageError: string = '';
  one_hundred_quantity: string = '';
  ten_quantity: string = '';
  quantity_filter: string | null = '';

  // Declarando variaveis boolean
  canView: boolean = false;
  canViewMachines: boolean = false;
  showMessage: boolean = false;

  // Declarando variaveis list
  so_list: string[] = [];

  // Função iniciada ao carregar a pagina
  ngOnInit(): void {
    // Pegando valores do usuario
    this.name = localStorage.getItem('name');

    // Pegando valor de quantitade do filtro
    this.quantity_filter = localStorage.getItem('quantity');
    if (this.quantity_filter == null) {
      localStorage.setItem('quantity', '10');
      this.quantity_filter = '10';
    }

    // Verificando se os dados existem
    if (this.name.length == 0 || this.name == null) {
      this.errorType = 'Falta de Dados';
      this.messageError =
        'Ouve um erro ao acessar dados do LDAP, contatar a TI';
      this.showMessage = true;
    } else {
      this.canView = true;

      this.getData();
    }
  }

  // Buscando as maquinas disponiveis
  getData(): void {
    this.http
      .get('/home/computers/get-data/' + this.quantity_filter, {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;
          // Ajustando filtro de quantidade
          switch (this.quantity_filter) {
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

          this.canViewMachines = true;
          this.getSO();
        }
      });
  }

  getSO(): void {
    this.http
      .get('/home/computers/get-data-SO', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.so_list = data.SO;
        }
      });
  }

  // Função para redirecionar para a pagina de visualização da maquina
  onRowClick(index: number) {
    const selectedMachine = this.dataMachines[index];

    var mac = selectedMachine[0];

    mac = mac.replace(/:/g, '-');

    return (window.location.href = '/home/computers/view-machine/' + mac);
  }

  // Função obter o valor do SO que deseja filtrar
  onRowClickSO(index: number) {
    this.canViewMachines = false;

    this.dataMachines = null;

    let so;

    if (index == 69) {
      so = 'all';
    } else {
      const selectedSO = this.so_list[index];

      so = selectedSO[0];
    }

    this.http
      .get(
        '/home/computers/get-data-SO-filter/' + this.quantity_filter + '/' + so,
        {}
      )
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  // Seta quantidade de maquinas a serem exibidas para 10
  getTen(): void {
    localStorage.setItem('quantity', '10');

    this.quantity_filter = '10';

    this.ten_quantity = 'active_filter';
    this.fifty_quantity = '';
    this.one_hundred_quantity = '';
    this.all_quantity = '';

    this.canViewMachines = false;
    this.dataMachines = null;

    this.http
      .get('/home/computers/get-quantity/10', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  // Seta quantidade de maquinas a serem exibidas para 50
  getFifty(): void {
    localStorage.setItem('quantity', '50');

    this.ten_quantity = '';
    this.fifty_quantity = 'active_filter';
    this.one_hundred_quantity = '';
    this.all_quantity = '';

    this.canViewMachines = false;
    this.dataMachines = null;

    this.http
      .get('/home/computers/get-quantity/50', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  // Seta a quantidade de maquinas a serem exibidas para 100
  getOneHundred(): void {
    localStorage.setItem('quantity', '100');

    this.ten_quantity = '';
    this.fifty_quantity = '';
    this.one_hundred_quantity = 'active_filter';
    this.all_quantity = '';

    this.canViewMachines = false;
    this.dataMachines = null;

    this.http
      .get('/home/computers/get-quantity/100', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  // Seta a quantidade de maquinas a serem exibidas para todas
  getAll(): void {
    localStorage.setItem('quantity', 'all');

    this.ten_quantity = '';
    this.fifty_quantity = '';
    this.one_hundred_quantity = '';
    this.all_quantity = 'active_filter';

    this.canViewMachines = false;
    this.dataMachines = null;

    this.http
      .get('/home/computers/get-quantity/all', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.dataMachines = data.machines;

          this.canViewMachines = true;
        }
      });
  }

  formatDate(date: string): string {
    const parsedDate = new Date(date);
    const day = String(parsedDate.getDate()).padStart(2, '0');
    const month = String(parsedDate.getMonth() + 1).padStart(2, '0'); // Meses são baseados em 0 (Janeiro é 0)
    const year = parsedDate.getFullYear();

    return `${day}/${month}/${year}`;
  }

  sortByName(): void {
    this.dataMachines.sort((a: any, b: any) => {
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

  sortDataByNameDescending(): void {
    this.dataMachines.sort((a: any, b: any) => {
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

  sortByNameSO(): void {
    this.dataMachines.sort((a: any, b: any) => {
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

  sortDataByNameDescendingSO(): void {
    this.dataMachines.sort((a: any, b: any) => {
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

  sortByNameDis(): void {
    this.dataMachines.sort((a: any, b: any) => {
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

  sortDataByNameDescendingDis(): void {
    this.dataMachines.sort((a: any, b: any) => {
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

  sortByDateDesc(): void {
    this.dataMachines.sort((a: any, b: any) => {
      const dateA = new Date(a[4]);
      const dateB = new Date(b[4]);

      return dateB.getTime() - dateA.getTime(); // Mais recente para o mais antigo
    });
  }

  sortByDateASC(): void {
    this.dataMachines.sort((a: any, b: any) => {
      const dateA = new Date(a[4]);
      const dateB = new Date(b[4]);

      return dateA.getTime() - dateB.getTime(); // Mais antigo para o mais novo
    });
  }
}
