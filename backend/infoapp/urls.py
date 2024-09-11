from django.urls import path
from . import views

urlpatterns = [
    path("", views.home, name="central-home"),
    path("get-Info-Main-Panel/", views.getInfoMainPanel, name="Get-Info-Main-Panel"),
    # URL para pegar os sistemas operacionais para o gráfico
    path("get-info-SO", views.getInfoSOChart, name="central-get-info-so-chart"),
    # URL para pegar as entradas dos computadores
    path(
        "get-info-last-update",
        views.getInfoLastUpdate,
        name="central-get-info-last-update",
    ),
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
    # URL que busca os SO para o filtro
    path("computers/get-data-SO", views.getDataSO, name="central-get-data-SO"),
    # URL que busca os distribution para o filtro
    path("computers/get-data-DIS", views.getDataDIS, name="central-get-data-DIS"),
    # URL que busca os computadores pelo fitlro do SO
    path(
        "computers/get-data-SO-filter/<str:quantity>/<str:so>",
        views.getDataSoFilter,
        name="central-get-data-so-filter",
    ),
    # URL que busca os computadores pelo fitlro de distribution
    path(
        "computers/get-data-DIS-filter/<str:quantity>/<str:dis>",
        views.getDataDisFilter,
        name="central-get-data-so-filter",
    ),
    # URL que busca os computadores pelo filtro de nome
    path(
        "computers/get-machine-varchar/<str:quantity>/<str:name>",
        views.getDataVarchar,
        name="central-get-data-varchar",
    ),
    # URL que gera relatorio de DNS que visa mostrar ip identicos
    path("computers/report-dns", views.getReportDNS, name="central-get-report-dns"),
    # URL para gerar um relatorio com as maquinas em XLS
    path(
        "computers/get-report/xls/", views.getReportXLS, name="central-get-report-xls"
    ),
]
