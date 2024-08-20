from django.contrib.auth import login
from django.contrib.auth import logout
from django.contrib.auth.decorators import login_required
from django.contrib.auth.models import User
from django.http import JsonResponse
from django.shortcuts import redirect
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_http_methods
from dotenv import load_dotenv
from ldap3 import ALL_ATTRIBUTES, SAFE_SYNC, Connection
from os import getenv
import json


@csrf_exempt
def credential(request):
    if request.method == "POST":
        username = None
        password = None
        data = None
        domain = None
        domain = None
        server = None
        server = None
        conn = None
        base_ldap = None
        response = None
        try:
            data = json.loads(request.body)
            username = data.get("username")
            password = data.get("password")

            load_dotenv()

            domain = getenv("DOMAIN_NAME")

            server = getenv("SERVER1")

            conn = Connection(
                server,
                f"{domain}\{username}",
                password,
                auto_bind=True,
                client_strategy=SAFE_SYNC,
            )

            base_ldap = getenv("LDAP_BASE")

            if conn.bind():
                conn.read_only = True
                search_filter = f"(sAMAccountName={username})"
                ldap_base_dn = base_ldap
                response = conn.search(
                    ldap_base_dn,
                    search_filter,
                    attributes=ALL_ATTRIBUTES,
                    search_scope="SUBTREE",
                    types_only=False,
                )

        except Exception as e:
            print(e)
            return JsonResponse({"status": "invalid access"}, status=401, safe=True)

        extractor = None
        information = None
        name = None
        acess_user = None
        try:
            extractor = response[2][0]
            information = extractor.get("attributes")

            groups = information.get("memberOf", [])

            if "displayName" in information:
                name = information["displayName"]

            acess_user = getenv("ACESS_USER")

            for item in groups:
                if acess_user in item:
                    acess = "User"
                    break  # Se encontrou, não precisa continuar procurandoç

            if acess:
                user, created = User.objects.get_or_create(username=username)

                # Autentica o usuário (mesmo sem senha, para fins de exemplo).
                user.backend = "django.contrib.auth.backends.ModelBackend"  # Define o backend de autenticação.
                login(request, user)  # Loga o usuário.
            return JsonResponse({"name": name}, status=200, safe=True)

        except Exception as e:
            print(e)
            return JsonResponse({"status": "invalid access"}, status=401, safe=True)

    if request.method == "GET":
        return redirect("login")


# Função que realiza logout
@login_required
@csrf_exempt
@require_http_methods(["GET"])
def logoutFunc(request):
    logout(request)
    return redirect("/login")
