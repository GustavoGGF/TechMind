from django.contrib.auth import login
from django.contrib.auth import logout
from django.contrib.auth.decorators import login_required
from django.contrib.auth.models import User
from django.http import JsonResponse, FileResponse
from django.shortcuts import redirect, render
from django.views.decorators.csrf import csrf_exempt
from dotenv import load_dotenv
from ldap3 import ALL_ATTRIBUTES, SAFE_SYNC, Connection
from os import getenv, path
from json import loads
from django.views.decorators.http import require_POST, require_GET
from logging import basicConfig, getLogger, WARNING

basicConfig(level=WARNING)
logger = getLogger(__name__)

# Exige esta view do requerimento de proteção contra falsificação de requisições cross-site (CSRF) e  Permite apenas requisições POST nesta view.
@csrf_exempt
@require_POST
def credential(request):
    # Inicialização de variáveis para armazenar credenciais e conexões.
    username = None
    password = None
    data = None
    domain = None
    server = None
    conn = None
    base_ldap = None
    response = None

    try:
        # Carrega o corpo da requisição JSON e obtém o nome de usuário e senha fornecidos.
        data = loads(request.body)
        username = data.get("username")
        password = data.get("password")

        # Carrega variáveis de ambiente para obter informações de configuração.
        load_dotenv()

        # Obtém o nome do domínio do ambiente.
        domain = getenv("DOMAIN_NAME")

        # Obtém o servidor do Active Directory ou LDAP do ambiente.
        server = getenv("SERVER1")

        # Cria uma conexão segura com o servidor usando as credenciais fornecidas.
        conn = Connection(
            server,
            f"{domain}\\{username}",
            password,
            auto_bind=True,  # Vincula automaticamente a conexão após a criação.
            client_strategy=SAFE_SYNC,  # Define a estratégia de cliente segura e síncrona.
        )

        # Base LDAP usada para buscas.
        base_ldap = getenv("LDAP_BASE")

        # Se a conexão for bem-sucedida, executa uma busca no diretório LDAP.
        if conn.bind():
            conn.read_only = True  # Define a conexão como somente leitura.
            search_filter = f"(sAMAccountName={username})"  # Filtro de busca baseado no nome de usuário.
            ldap_base_dn = base_ldap
            response = conn.search(
                ldap_base_dn,  # DN base para a busca.
                search_filter,  # Filtro definido.
                attributes=ALL_ATTRIBUTES,  # Busca todos os atributos.
                search_scope="SUBTREE",  # Busca recursivamente em toda a árvore.
                types_only=False,  # Recupera os valores reais dos atributos.
            )

    except Exception as e:
        # Em caso de erro, imprime a exceção e retorna uma resposta JSON com status 401.
        logger.error(e)
        return JsonResponse({"status": "invalid access"}, status=401, safe=True)

    # Inicialização de variáveis para extrair informações do resultado da busca.
    extractor = None
    information = None
    name = None
    acess_user = None

    try:
        # Extrai a resposta da busca LDAP e acessa os atributos retornados.
        extractor = response[2][0]
        information = extractor.get("attributes")

        # Obtém a lista de grupos do usuário.
        groups = information.get("memberOf", [])

        # Obtém o nome de exibição, se disponível.
        if "displayName" in information:
            name = information["displayName"]

        # Nome do grupo de acesso permitido, definido no ambiente.
        acess_user = getenv("ACESS_USER")

        # Verifica se o usuário pertence ao grupo de acesso.
        for item in groups:
            if acess_user in item:
                acess = "User"
                break  # Interrompe o loop ao encontrar o grupo de acesso.

        if acess:
            # Cria ou recupera um usuário Django com o nome de usuário fornecido.
            user, created = User.objects.get_or_create(username=username)

            # Define o backend de autenticação para permitir o login manual.
            user.backend = "django.contrib.auth.backends.ModelBackend"

            if created:
                user.save()

            # Realiza o login do usuário.
            login(request, user)

        # Retorna uma resposta JSON com o nome do usuário.
        return JsonResponse({"name": name}, status=200, safe=True)

    except Exception as e:
        # Em caso de erro, imprime a exceção e retorna uma resposta JSON com status 401.
        logger.error(e)
        return JsonResponse({"status": "invalid access"}, status=401, safe=True)


# Função que realiza logout
# Requer que o usuário esteja autenticado para acessar esta view.
# Permite apenas requisições GET para esta view.
@login_required
@require_GET
def logout_func(request):
    try:
        # Executa o logout do usuário atual.
        logout(request)

        # Retorna uma resposta JSON vazia com status 200 (OK) indicando sucesso.
        return JsonResponse({}, status=200)

    except Exception as e:
        # Em caso de erro, imprime a exceção para depuração.
        logger.error(e)


# Permite apenas requisições GET para esta view.
@require_GET
def donwload_files(request):
    try:
        # Verifica se a requisição é uma requisição AJAX (XMLHttpRequest).
        if (
            "X-Requested-With" not in request.headers
            or request.headers["X-Requested-With"] != "XMLHttpRequest"
        ):
            # Se não for uma requisição AJAX, redireciona para a página home.
            return redirect("/home")

        # Caminho para o arquivo que será baixado.
        file_path = "/node/TechMind/Installers/techmind.exe"

        # Verifica se o arquivo existe no caminho especificado.
        if path.exists(file_path):
            # Se o arquivo existe, retorna o arquivo como resposta de download.
            response = FileResponse(
                open(file_path, "rb"), as_attachment=True, filename="techmind.exe"
            )
            return response
        else:
            # Se o arquivo não for encontrado, retorna um erro com status 400.
            return JsonResponse({}, status=400)

    except Exception as e:
        # Em caso de erro, imprime a exceção para depuração.
        logger.error(e)

        # Retorna uma resposta de erro com status 300 (não convencional, seria mais apropriado usar outro código de erro).
        return JsonResponse({}, status=300)

def donwload_techMind(request):
    return logger.error("foi")