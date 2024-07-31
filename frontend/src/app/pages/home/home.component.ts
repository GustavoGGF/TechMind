import { Component } from '@angular/core';
import { UtilitiesModule } from '../../utilities/utilities.module';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { catchError, throwError } from 'rxjs';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [UtilitiesModule, CommonModule, HttpClientModule],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css',
})
export class HomeComponent {
  constructor(private http: HttpClient) {}

  name: any;
  status: any;

  errorType: string = '';
  messageError: string = '';

  canView: boolean = false;
  notData: boolean = true;
  showMessage: boolean = false;

  totalMachines: number = 0;
  totalWindows: number = 0;
  totalUnix: number = 0;

  ngOnInit() {
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

  hideMessage() {
    this.showMessage = false;
  }

  getData() {
    this.http
      .get('/home/get-Info-Main-Panel/', {})
      .pipe(
        catchError((error) => {
          this.status = error.status;

          return throwError(error);
        })
      )
      .subscribe((data: any) => {
        if (data) {
          this.totalWindows = data.windows;
          this.totalMachines = data.total;
          this.totalUnix = data.unix;
          this.notData = false;
        }
      });
  }
}
