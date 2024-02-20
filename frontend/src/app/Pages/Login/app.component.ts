import { Component, ViewChild, ElementRef } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { LoadingComponent } from '../../utils/loading/loading.component';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { HttpClient } from '@angular/common/http';
// import { Location } from '@angular/common';
// import { Observable } from 'rxjs';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, LoadingComponent, CommonModule, HttpClientModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
})
export class AppComponent {
  title = 'TechMind';
  elements: any;
  canShow: boolean = false;

  constructor(private elementRef: ElementRef, private http: HttpClient) {}

  // private http: HttpClient

  urlImage: string = '../../../assets/images/logo/Logo_TechMind.png';

  nextAnimate: boolean = false;

  letters: any;

  currentUrl: any;

  @ViewChild('logo') logo!: ElementRef | undefined;
  @ViewChild('main') main!: ElementRef | undefined;

  ngAfterViewInit(): void {
    this.currentUrl = window.location.href;
    this.http.get(this.currentUrl + 'api').subscribe((result) => {
      console.log(result);
    });

    if (this.logo) {
      this.logo.nativeElement.addEventListener('animationend', () => {
        const letters =
          this.elementRef.nativeElement.querySelectorAll('.letter');

        letters.forEach((letter: any) => {
          (letter as HTMLElement).classList.add('animate1');
        });
      });
    }
    if (this.main) {
      this.main.nativeElement.addEventListener('keyup', (event: any) => {
        if (event.keyCode === 13) {
          this.loginInput();
        }
      });
    }
  }

  name: string = '';
  pass: string = '';
  // credential: any;

  getUser(event: any): void {
    this.name = event.target.value;
  }
  getPass(event: any): void {
    this.pass = event.target.value;
  }

  data = { user: this.name, pass: this.pass };

  loginInput(): void {
    if (this.name.length > 6 && this.pass.length > 10) {
      this.canShow = true;

      this.http
        .post(this.currentUrl + 'api/' + 'credentials', this.data)
        .subscribe((result) => {
          console.warn(result);
        });
    }
  }
}
