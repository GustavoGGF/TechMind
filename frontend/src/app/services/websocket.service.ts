import { Injectable } from "@angular/core";
import { webSocket } from "rxjs/webSocket";

@Injectable({
  providedIn: "root",
})
export class WebSocketService {
  protocol = window.location.protocol === "https:" ? "wss" : "ws";
  socket = webSocket(
    `${this.protocol}://${window.location.host}/ws/server-communication/`
  );

  getMessages() {
    return this.socket;
  }

  sendMessage(msg: any) {
    this.socket.next(msg);
  }
}
