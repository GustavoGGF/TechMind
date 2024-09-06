import os
from django.core.asgi import get_asgi_application
from channels.routing import ProtocolTypeRouter
from .routing import application as websocket_application

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "techmind.settings")

techmind_asgi_app = get_asgi_application()

application = ProtocolTypeRouter(
    {
        "http": techmind_asgi_app,
        "websocket": websocket_application,
    }
)
