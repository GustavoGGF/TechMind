import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css',
})
export class NavbarComponent {
  @Input() name: string = '';

  logoTechMind: string = '/static/assets/images/Logo_TechMind.png';
}
