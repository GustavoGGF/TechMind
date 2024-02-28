import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-message',
  standalone: true,
  imports: [],
  templateUrl: './message.component.html',
  styleUrl: './message.component.css',
})
export class MessageComponent {
  @Input() errorType: string = 'Error';
  @Input() messageError: string = 'Message Error';

  closeIMG = '../assets/images/utils/fechar.png';
}
