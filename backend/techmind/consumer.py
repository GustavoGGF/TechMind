# consumers.py
import json
from channels.generic.websocket import WebsocketConsumer


class MyConsumer(WebsocketConsumer):
    def connect(self):
        self.accept()

    def receive(self, text_data):
        data = json.loads(text_data)
        # Processo de dados ou envio para um grupo
        self.send(text_data=json.dumps({"message": "Mensagem recebida", "data": data}))
