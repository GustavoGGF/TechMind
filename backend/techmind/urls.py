from django.contrib import admin
from django.urls import path, include
from django.views.generic import TemplateView
from django.contrib.staticfiles.urls import staticfiles_urlpatterns
from . import views

urlpatterns = [
    path("admin/", admin.site.urls),
    path("", TemplateView.as_view(template_name="index.html"), name="login"),
    path("login/", TemplateView.as_view(template_name="index.html"), name="login"),
    path("credential/", views.credential, name="central-credential"),
    path("home/", include("infoapp.urls")),
    # URL para fazer logout
    path("logout/", views.logoutFunc, name="central-logout"),
    path('donwload-files/', views.donwloadFiles, name='central-donwload-files')
]

urlpatterns += staticfiles_urlpatterns()
