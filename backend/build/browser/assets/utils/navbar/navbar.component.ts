import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import '../../bootstrap-5.3.3-dist/js/bootstrap.js';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css',
})
export class NavbarComponent {
  img: string = '../assets/images/logo/Logo_TechMind.png';
  Name: any;

  ngAfterViewInit() {
    const User = localStorage.getItem('name');
    if (User !== null) {
      this.Name = JSON.parse(User);
    }
  }
}
