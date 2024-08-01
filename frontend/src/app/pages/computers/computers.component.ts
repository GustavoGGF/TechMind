import { Component } from '@angular/core';
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
export class ComputersComponent {
  constructor(private http: HttpClient) {}
  // Declarando variaveis any
  dataMachines: any;
  name: any;
  status: any;

  // Declarando variaveis string
  computers_class: string = 'active';
  device_class: string = '';
  errorType: string = '';
  home_class: string = '';
  messageError: string = '';

  // Declarando variaveis boolean
  canView: boolean = false;
  canViewMachines: boolean = false;
  showMessage: boolean = false;

  // Função iniciada ao carregar a pagina
  ngOnInit(): void {
    // Pegando valores do usuario
    this.name = localStorage.getItem('name');
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
      .get('/home/computers/get-data', {})
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

  // Função para redirecionar para a pagina de visualização da maquina
  onRowClick(index: number) {
    const selectedMachine = this.dataMachines[index];

    var mac = selectedMachine[0];

    mac = mac.replace(/:/g, '-');

    return (window.location.href = '/home/computers/view-machine/' + mac);
  }
}
