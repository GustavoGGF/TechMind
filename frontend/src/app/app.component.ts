import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Component, ElementRef, ViewChild } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { LoadingComponent } from '../assets/components/loading/loading.component';
import { CommonModule } from '@angular/common';
import 'animate.css';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, HttpClientModule, LoadingComponent, CommonModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
})
export class AppComponent {
  title = 'frontend';

  urlImage = '../assets/images/logo/Logo_TechMind.png';

  constructor(private elementRef: ElementRef, private http: HttpClient) {}

  name: string = '';
  pass: string = '';
  canShow: boolean = false;

  @ViewChild('logo') logo: ElementRef | undefined;
  @ViewChild('main') main: ElementRef | undefined;

  ngAfterViewInit(): void {
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

  getUser(event: any): void {
    this.name = event.target.value;
  }
  getPass(event: any): void {
    this.pass = event.target.value;
  }

  loginInput(): void {
    if (this.name && this.pass) {
      this.canShow = true;

      const currentUrl = window.location.href;

      this.http.post(currentUrl + 'api/credential', {
        user: this.name,
        pass: this.pass,
      });
      // .subscribe(() => (this.canShow = false));
    }
  }
}
