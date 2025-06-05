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
    path("logout/", views.logout_func, name="central-logout"),
    path('download-files/<str:file>/<str:version>', views.donwload_files, name='central-donwload-files'),
    path("get-current-version/<str:os>", views.get_current_version, name="central-get-current-version")
]

urlpatterns += staticfiles_urlpatterns()
