from channels.generic.websocket import AsyncWebsocketConsumer
import json

class MonitoringConsumer(AsyncWebsocketConsumer):
    async def connect(self):
        self.group_name = "monitoring"
        await self.channel_layer.group_add(self.group_name, self.channel_name)
        await self.accept()
        await self.send(text_data=json.dumps({"status": "conectado"}))
        
    async def disconnect(self, close_code):
       await self.channel_layer.group_discard(self.group_name, self.channel_name)
        
    async def send_mensagem(self, event):
        await self.send(text_data=json.dumps({
            "message": event["message"],
            "code":event["code"],
            "status":event["status"]
        }))