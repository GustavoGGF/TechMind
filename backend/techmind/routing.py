from django.urls import re_path
from techmind.consumer import MyConsumer

websocket_urlpatterns = [
    re_path(r"ws/home/", MyConsumer.as_asgi()),
]
