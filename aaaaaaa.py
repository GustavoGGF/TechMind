from ldap3 import Server, Connection, ALL, SUBTREE
import dns.resolver


def get_ip_from_dns(hostname):
    try:
        # Consulta DNS para registros A (IPv4)
        answers = dns.resolver.resolve(hostname, "A")
        ips = [rdata.address for rdata in answers]
        return ips
    except Exception as e:
        print(f"An error occurred while resolving DNS for {hostname}: {e}")
        return []


def query_ad_computers(
    ldap_server, ldap_user, ldap_password, ldap_base_dn, search_filter
):
    try:
        # Conectar ao servidor LDAP
        server = Server(ldap_server, get_info=ALL)
        conn = Connection(
            server,
            user=ldap_user,
            password=ldap_password,
            auto_bind=True,
            read_only=True,
        )

        # Realizar a pesquisa
        conn.search(
            search_base=ldap_base_dn,
            search_filter=search_filter,
            attributes=["dnsHostName"],
            search_scope=SUBTREE,
            types_only=False,
        )

        # Processar e imprimir os resultados
        for entry in conn.entries:
            if "dnsHostName" in entry:
                hostname = entry.dnsHostName.value
                if hostname:
                    ips = get_ip_from_dns(hostname)
                    print(f"Hostname: {hostname}, IPs: {', '.join(ips)}")
                else:
                    print(f"Hostname is empty for entry: {entry.entry_dn}")
            else:
                print(f"No dnsHostName attribute found in entry: {entry.entry_dn}")

    except Exception as e:
        print(f"An error occurred: {e}")


if __name__ == "__main__":
    ldap_server = "sdc01.nt-lupatech.com.br"  # Substitua pelo servidor LDAP
    ldap_user = "nt-lupatech\gfreitas"  # Substitua pelo seu usuário LDAP
    ldap_password = "GFLup@mae!02"  # Substitua pela sua senha LDAP
    ldap_base_dn = "dc=nt-lupatech,dc=com,dc=br"  # Substitua pelo DN base apropriado
    search_filter = "(objectClass=computer)"  # Filtro de pesquisa para computadores

    query_ad_computers(
        ldap_server, ldap_user, ldap_password, ldap_base_dn, search_filter
    )
