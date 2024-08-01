from django.urls import path
from . import views

urlpatterns = [
    path("", views.home, name="central-home"),
    path("get-Info-Main-Panel/", views.getInfoMainPanel, name="Get-Info-Main-Panel"),
    path("computers/", views.computers, name="computers"),
    # Função da url /computers que solicita dados dos computadores conforme quantidade solicitada
    path(
        "computers/get-data/<str:quantity>",
        views.getDataComputers,
        name="central-get-data-computers",
    ),
    path("computers/post-machines", views.postMachines, name="central-post-machines"),
    path(
        "computers/view-machine/<str:mac_address>",
        views.postMachinesWithMac,
        name="central-post-machines",
    ),
    path(
        "computers/info-machine/<str:mac_address>",
        views.infoMachine,
        name="central-post-machines",
    ),
    path(
        "devices/",
        views.devices,
        name="central-devices",
    ),
    path(
        "devices/post-devices",
        views.devices_post,
        name="central-devices",
    ),
    path(
        "devices/get-devices/<str:quantity>",
        views.devices_get,
        name="central-devices-get",
    ),
    path(
        "devices/view-devices/<str:sn>",
        views.devices_details,
        name="central-devices-details",
    ),
    path(
        "devices/info-device/<str:sn>",
        views.infoDevice,
        name="central-info-device",
    ),
    path("devices/get-last-machines", views.lastMachines, name="central-last-machines"),
    path("computers/added-device", views.addedDevices, name="central-added-device"),
    path(
        "computers/added-devices/<str:mac_address>",
        views.computersDevices,
        name="central-computers-device",
    ),
    # url de modificação da aba outros
    path(
        "computers/modify-others/<str:mac_address>",
        views.computersModify,
        name="centra-computers-modify",
    ),
    # url que libera o token csrf
    path("get-token", views.getToken, name="central-get-token"),
    # URL que muda a quantidade de computadores conforme solicitação
    path(
        "computers/get-quantity/<str:quantity>",
        views.getQuantity,
        name="central-get-quantity",
    ),
]
