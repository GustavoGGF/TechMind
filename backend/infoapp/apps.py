from django.apps import AppConfig
from os import path

class InfoappConfig(AppConfig):
    default_auto_field = 'django.db.models.BigAutoField'
    name = 'infoapp'
    path = path.dirname(path.abspath(__file__))
