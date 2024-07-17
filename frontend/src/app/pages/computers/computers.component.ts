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

  name: any;
  status: any;

  errorType: string = '';
  messageError: string = '';

  canView: boolean = false;
  canViewMachines: boolean = false;
  showMessage: boolean = false;

  dataMachines: any;

  ngOnInit(): void {
    this.name = localStorage.getItem('name');

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

  onRowClick(index: number) {
    const selectedMachine = this.dataMachines[index];

    var mac = selectedMachine[0];

    mac = mac.replace(/:/g, '-');

    return (window.location.href = '/home/computers/view-machine/' + mac);
  }
}
