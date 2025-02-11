import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-message',
  templateUrl: './message.component.html',
  styleUrl: './message.component.css',
})
export class MessageComponent {
  @Input() errorType: string = 'Error';
  @Input() messageError: string = 'Message Error';

  closeIMG = '/static/assets/images/fechar.png';
}
