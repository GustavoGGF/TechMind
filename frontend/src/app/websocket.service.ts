import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root', // Isso garante que o serviço esteja disponível em toda a aplicação
})
export class WebsocketService {
  private ws: WebSocket;
  private messageSubject = new BehaviorSubject<any>(null);

  message = this.messageSubject.asObservable();

  constructor() {
    this.ws = new WebSocket('ws://localhost:3000/ws/home'); // URL do WebSocket

    this.ws.onmessage = (event) => {
      this.messageSubject.next(event.data);
    };
  }

  sendMessage(msg: string) {
    this.ws.send(msg);
  }
}
