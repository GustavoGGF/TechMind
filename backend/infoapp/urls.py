from django.urls import path
from . import views

urlpatterns = [
    path("", views.home, name="central-home"),
    path("get-Info-Main-Panel/", views.getInfoMainPanel, name="Get-Info-Main-Panel"),
    path("computers/", views.computers, name="computers"),
    path(
        "computers/get-data", views.getDataComputers, name="central-get-data-computers"
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
    path("devices/get-token", views.devices_token, name="central-devices-token"),
    path("devices/get-devices", views.devices_get, name="central-devices-get"),
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
]
