package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-ldap/ldap"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

type SystemCount struct{
    SystemName string
    Count int
}

func main() {                        
    app := fiber.New()
    
    app.Static("/", "./build/browser")   

    app.Use("/api", func(c *fiber.Ctx) error{
        return c.Next()
    })

    app.Post("api/credential", func(c *fiber.Ctx) error {
        var data map[string]string

        if err := c.BodyParser(&data); err != nil {
            fmt.Println("Erro ao converter o body: ", err)
            return c.Next()
        }
    
        user := data["user"]

        username := fmt.Sprintf("nt-lupatech\\%s", user) //variavel de ambiente
        pass := data["pass"]

        l, err := ldap.DialURL("ldap://sdc01.nt-lupatech.com.br") //env
        if err != nil {
            fmt.Println("Erro ao pegar o servidor ldap: ",err)
        }

        defer l.Close()

        err = l.Bind(username, pass)
        if err != nil {
            if ldapErr, ok := err.(*ldap.Error); ok {
                if ldapErr.ResultCode == ldap.LDAPResultInvalidCredentials {
                    return c.JSON(fiber.Map{"status":404})
                } else {
                    fmt.Printf("Erro LDAP: %v", ldapErr)
                }
            } else {
                fmt.Printf("Erro não identificado: %v", err)
            }
        }

        Filter := fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", user)

        searchRequest := ldap.NewSearchRequest(
            "ou=Brasil,dc=nt-lupatech,dc=com,dc=br",
            ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0,
            false,
            Filter,
            []string{"sAMAccountName", "givenName", "memberOf", "displayName"},
            nil,
        )

        searchResult, err := l.Search(searchRequest)
        if err != nil {
            fmt.Println("Erro ao pesquisar dados no ldap",err)
        }

        verify := os.Getenv("VERIFY_GROUP")       

        found := false
        for _, entry := range searchResult.Entries {
            for _, attr := range entry.Attributes {
                if attr.Name == "memberOf" {
                    groups := attr.Values

                    for _, value := range groups{
                        if strings.Contains(value, verify){
                            found = true
                            break                            
                        }                  
                    } 
                    if !found{
                        return c.JSON(fiber.Map{"status":401})
                    } else{
                        getName := entry.GetAttributeValue("displayName")

                        

                        return c.JSON(fiber.Map{"status":200, "name": getName})  
                    }    
                }
            }
        }   

        return c.JSON(fiber.Map{"status":200})
    })

    app.Post("api/machines", func(c *fiber.Ctx) error {
        var data map[string]string

        if err := c.BodyParser(&data); err != nil {
            fmt.Println("Erro ao converter dados do body: ", err)
            return c.Next()
        }

        name := data["Name"]
        system_name := data["System"]
        distribution := data["Distribution"]    
        web_interface := data["InterfaceInternet"]  
        mac_address := data["MacAddress"]

        db, err := sql.Open("mysql", "mach:Lup@.CSC.!@tcp(10.1.9.0:3306)/techmindDB") //varaivel de ambiente
        if err != nil{
            fmt.Println("Erro conectar no mysql: ", err)
            return c.Next()
        }

        defer db.Close()

        err = db.Ping()
        if err != nil {
            fmt.Println("Erro ao pingar o servidor: ", err)
            return c.Next()
        }

        verify_query := fmt.Sprintf("SELECT COUNT(ID) FROM machines WHERE mac_address = '%s'", mac_address)

        var repet int
        err = db.QueryRow(verify_query).Scan(&repet)
        if err != nil {
            log.Fatal("Erro ao pegar número total de máquinas:", err)
        }

        if repet >= 1{
            updateQuery := "UPDATE machines SET name = ?, system_name = ?, distribution = ?, web_interface = ? WHERE mac_address = ?"

            result, err := db.Exec(updateQuery, name, system_name, distribution, web_interface, mac_address)
            if err != nil {
                log.Fatal("Erro ao atualizar os valores:", err)
            }

            rowsAffected, err := result.RowsAffected()
            if err != nil {
                log.Fatal("Erro ao obter o número de linhas afetadas:", err)
            }

            if rowsAffected > 0 {
                return c.Next()
            } else {
                return c.Next()
            }
        } else {
            stmt, err := db.Prepare("INSERT INTO machines(name, system_name, distribution, web_interface, mac_address) VALUES(?, ?, ?, ?, ?)")
            if err != nil{
                fmt.Println("Erro ao preparar inserção de dados: ",err)
                return c.Next()
            }
    
            defer stmt.Close()
    
            _, err = stmt.Exec(name, system_name, distribution, web_interface, mac_address)
            if err != nil {
                fmt.Println("Erro ao inserir os dados no mysql:", err)
                return c.Next()
            }
            return c.Next()
        }

    })

    app.Get("home", func(c *fiber.Ctx) error {
    dataSourceName := "mach:Lup@.CSC.!@tcp(10.1.9.0:3306)/techmindDB" //env

    // Abrindo uma conexão com o banco de dados MySQL
    db, err := sql.Open("mysql", dataSourceName)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Verificar se a conexão com o banco de dados está ativa
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    // Consulta SQL para contar o número de IDs na tabela
    query := "SELECT COUNT(ID) FROM machines" //env
    query2 := "SELECT COUNT(*) AS TotalLinuxSystems FROM machines WHERE system_name LIKE '%linux%'"
    query3 := "SELECT COUNT(*) AS TotalLinuxSystems FROM machines WHERE system_name LIKE '%windows%'"
    query4 := "SELECT distribution, COUNT(*) AS count FROM machines GROUP BY distribution"

    // Executando a consulta
    var totalMachines int
    var totalLinux int
    var totalWindows int
    var systemCounts []SystemCount

    err = db.QueryRow(query).Scan(&totalMachines)
    if err != nil {
        log.Fatal("Erro ao pegar número total de máquinas: ",err)
        return c.Next()
    }

    err = db.QueryRow(query2).Scan(&totalLinux)
    if err != nil {
        log.Fatal("Erro ao pegar o número total de linux: ",err)
        return c.Next()
    }

    err = db.QueryRow(query3).Scan(&totalWindows)
    if err != nil {
        log.Fatal("Erro ao pegar o número total de Windows: ",err)
        return c.Next()
    }

    rows, err := db.Query(query4)
    if err != nil {
        log.Fatal("Erro ao executar a consulta: ", err)
        return c.Next()
    }
    defer rows.Close()
    
    for rows.Next() {
        var systemCount SystemCount
        if err := rows.Scan(&systemCount.SystemName, &systemCount.Count); err != nil {
            log.Fatal("Erro ao escanear a linha: ", err)
            return c.Next()
        }
        systemCounts = append(systemCounts, systemCount)
    }
    
    if err := rows.Err(); err != nil {
        log.Fatal("Erro ao iterar sobre os resultados: ", err)
        return c.Next()
    }

    return c.JSON(fiber.Map{"status":200, "machines":totalMachines, "linux":totalLinux, "windows":totalWindows, "systems": systemCounts})
    })


    app.Use(func(c *fiber.Ctx) error{
        return c.SendFile("./build/browser/index.html") //variavel de ambiente
    })

    if err := app.Listen(":3000"); err != nil{
        log.Fatalf("Erro ao iniciar o Servidor %v", err)
    }
}