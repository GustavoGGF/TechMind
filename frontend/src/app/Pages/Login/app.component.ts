import { Component, ViewChild, ElementRef } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
})
export class AppComponent {
  elements: any;
  constructor(private elementRef: ElementRef) {}
  title = 'TechMind';

  urlImage: string = '../../../assets/images/logo/Logo_TechMind.png';

  nextAnimate: boolean = false;

  letters: any;

  @ViewChild('logo') logo!: ElementRef | undefined;

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
  }
}
