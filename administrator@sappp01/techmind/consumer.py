# consumers.py
import json
from channels.generic.websocket import WebsocketConsumer


class MyConsumer(WebsocketConsumer):
    async def connect(self):
        self.accept()

    async def disconnect(self, close_code):
        pass

    async def receive(self, text_data):
        data = json.loads(text_data)
        # Processo de dados ou envio para um grupo
        self.send(text_data=json.dumps({"message": "Mensagem recebida", "data": data}))
