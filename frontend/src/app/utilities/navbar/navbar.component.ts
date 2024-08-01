import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css',
})
export class NavbarComponent {
  // Declarando e exportando as variaveis que serão usadas no componente
  @Input() name: string = '';
  @Input() home_class: string = '';
  @Input() computers_class: string = '';
  @Input() device_class: string = '';

  logoTechMind: string = '/static/assets/images/Logo_TechMind.png';
}
