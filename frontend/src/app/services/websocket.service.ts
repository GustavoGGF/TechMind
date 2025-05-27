import { Injectable } from "@angular/core";
import { webSocket } from "rxjs/webSocket";

@Injectable({
  providedIn: "root",
})
export class WebSocketService {
  private socket = webSocket("ws://localhost:8000/ws/server-communication/");

  getMessages() {
    return this.socket;
  }

  sendMessage(msg: any) {
    this.socket.next(msg);
  }
}
