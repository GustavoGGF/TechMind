from pathlib import Path
from decouple import config
from os import path

# Build paths inside the project like this: BASE_DIR / 'subdir'.
BASE_DIR = Path(__file__).resolve().parent.parent

SECRET_KEY = "django-insecure--0q*51qtu#ja%#(%&x1us^gmsj&*jb&!u*@oh60#=k@67xwx!j"

DEBUG = True
NAME = None
USER = None
PASSWORD = None
HOST = None
PORT = None

ALLOWED_HOSTS = ["localhost", "techmind","techmind-backup", "techmind.lupatech.com.br", "127.0.0.1"]

CORS_ALLOWED_ORIGINS = [
    "https://techmind.lupatech.com.br",
    "https://sappp01.lupatech.com.br/",
]

CSRF_TRUSTED_ORIGINS = [
    "https://techmind.lupatech.com.br",
    "https://sappp01.lupatech.com.br/",
]

# Application definition

INSTALLED_APPS = [
    "django.contrib.admin",
    "django.contrib.auth",
    "django.contrib.contenttypes",
    "django.contrib.sessions",
    "django.contrib.messages",
    "django.contrib.staticfiles",
    "techmind",
    "infoapp",
]

MIDDLEWARE = [
    "django.middleware.security.SecurityMiddleware",
    "django.contrib.sessions.middleware.SessionMiddleware",
    "corsheaders.middleware.CorsMiddleware",
    "django.middleware.common.CommonMiddleware",
    "django.middleware.csrf.CsrfViewMiddleware",
    "django.contrib.auth.middleware.AuthenticationMiddleware",
    "django.contrib.messages.middleware.MessageMiddleware",
    "django.middleware.clickjacking.XFrameOptionsMiddleware",
    "corsheaders.middleware.CorsMiddleware",
]

ROOT_URLCONF = "techmind.urls"

TEMPLATES = [
    {
        "BACKEND": "django.template.backends.django.DjangoTemplates",
        "DIRS": [BASE_DIR / "static"],
        "APP_DIRS": True,
        "OPTIONS": {
            "context_processors": [
                "django.template.context_processors.debug",
                "django.template.context_processors.request",
                "django.contrib.auth.context_processors.auth",
                "django.contrib.messages.context_processors.messages",
            ],
        },
    },
]

WSGI_APPLICATION = "techmind.wsgi.application"
ASGI_APPLICATION = "techmind.asgi.application"

DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.mysql",
        "NAME": config("DB_NAME"),
        "USER": config("DB_USER"),
        "PASSWORD": config("DB_PASSWORD"),
        "HOST": config("DB_HOST"),
        "PORT": config("DB_PORT", cast=int),
    }
}

# Channels
CHANNEL_LAYERS = {
    "default": {
        "BACKEND": "channels_redis.core.RedisChannelLayer",
        "CONFIG": {
            "hosts": [("127.0.0.1", 6379)],
        },
    },
}


AUTH_PASSWORD_VALIDATORS = [
    {
        "NAME": "django.contrib.auth.password_validation.UserAttributeSimilarityValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.MinimumLengthValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.CommonPasswordValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.NumericPasswordValidator",
    },
]

# settings.py

LOGGING = {
    "version": 1,
    "disable_existing_loggers": False,
    "handlers": {
        "console": {
            "level": "DEBUG",
            "class": "logging.StreamHandler",
        },
    },
    "root": {
        "handlers": ["console"],
        "level": "DEBUG",
    },
}


LANGUAGE_CODE = "pt-br"

TIME_ZONE = "Brazil/East"

USE_I18N = True

USE_TZ = True

STATIC_URL = "/static/"
STATIC_ROOT = (
    BASE_DIR / "/var/www/hmt/static/"
)  # Diretório onde os arquivos estáticos serão coletados para produção
STATICFILES_DIRS = [
    path.join(BASE_DIR, "static"),
]
DEFAULT_AUTO_FIELD = "django.db.models.BigAutoField"


USE_X_FORWARDED_HOST = True
USE_X_FORWARDED_PROTO = True