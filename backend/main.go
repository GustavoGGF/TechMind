package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-ldap/ldap"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

type SystemCount struct{
    SystemName string
    Count int
}

type CityCount struct{
    CityName string
    Count int
}

func substituirValor(systemCounts []CityCount, valorBuscado, substituto string) []CityCount {
    for i, systemCount := range systemCounts {
        if strings.Contains(systemCount.CityName, valorBuscado) {
            systemCounts[i].CityName = substituto
        }
    }
    return systemCounts
}

func main() {                        
    app := fiber.New()

    // Or extend your config for customization

    // app.Use(func(c *fiber.Ctx) error {
    //     if c.Path() == "/api/credential" {
    //         return c.Next()
    //     }
    //     return csrf.New(csrf.Config{
    //         KeyLookup:      "header:X-Csrf-Token",
    //         CookieName:     "csrf_",
    //         CookieSameSite: "Lax",
    //         Expiration:     10 * time.Second,
    //         KeyGenerator:   utils.UUIDv4,
    //     })(c)
    // })
    
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

                        csrfToken := c.Locals("csrf_token")

                        return c.JSON(fiber.Map{"status":200, "name": getName, "csrf_token":csrfToken})  
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
        date := data["InsertionDate"]

        db, err := sql.Open("mysql", "mach:Lup@.CSC.!@tcp(10.1.9.19:3306)/techmindDB") //varaivel de ambiente
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
            layout := "2006-01-02 15:04"

            // Parsing da string para o tipo time.Time
            insertion_date, err := time.Parse(layout, date)
            if err != nil {
                fmt.Println("Erro ao converter string para data:", err)
                return c.Next()
            }
        
            updateQuery := "UPDATE machines SET name = ?, system_name = ?, distribution = ?, web_interface = ?, insertion_date = ? WHERE mac_address = ?"

            result, err := db.Exec(updateQuery, name, system_name, distribution, web_interface, insertion_date, mac_address)
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
                log.Fatal("Nada foi Atualizado")
                return c.Next()
            }
        } else {
            stmt, err := db.Prepare("INSERT INTO machines(name, system_name, distribution, web_interface, mac_address, insetion_date) VALUES(?, ?, ?, ?, ?, ?)")
            if err != nil{
                fmt.Println("Erro ao preparar inserção de dados: ",err)
                return c.Next()
            }
    
            defer stmt.Close()

            layout := "2006-01-02 15:04"

            // Parsing da string para o tipo time.Time
            insertion_date, err := time.Parse(layout, date)
            if err != nil {
                fmt.Println("Erro ao converter string para data:", err)
                return c.Next()
            }
    
            _, err = stmt.Exec(name, system_name, distribution, web_interface, mac_address, insertion_date)
            if err != nil {
                fmt.Println("Erro ao inserir os dados no mysql:", err)
                return c.Next()
            }
            return c.Next()
        }

    })

    app.Post("/home", func(c *fiber.Ctx) error {
        // Verificar se o token CSRF foi validado
        if csrfToken := c.Locals("csrf_token"); csrfToken != nil {
            // Token CSRF válido
            return c.SendString("Token CSRF válido")
        }
        // Token CSRF inválido ou ausente
        return c.Status(fiber.StatusForbidden).SendString("Token CSRF inválido ou ausente")
    })

    app.Get("api/home", func(c *fiber.Ctx) error {
    dataSourceName := "mach:Lup@.CSC.!@tcp(10.1.9.19:3306)/techmindDB" //env

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
    query3 := "SELECT COUNT(*) AS TotalWindowsSystems FROM machines WHERE system_name LIKE '%windows%'"
    // Executando a consulta
    var totalMachines int
    var totalLinux int
    var totalWindows int
   

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

    return c.JSON(fiber.Map{"status":200, "machines":totalMachines, "linux":totalLinux, "windows":totalWindows})
    })

    app.Get("api/machines", func(c *fiber.Ctx) error {

        dataSourceName := "mach:Lup@.CSC.!@tcp(10.1.9.19:3306)/techmindDB" //env

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

        query := "SELECT distribution, COUNT(*) AS count FROM machines GROUP BY distribution"


        var systemCounts []SystemCount

        rows, err := db.Query(query)
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
        
        return c.JSON(fiber.Map{"status":200, "systems": systemCounts})
    })

    app.Get("api/cities", func(c *fiber.Ctx) error {
        dataSourceName := "mach:Lup@.CSC.!@tcp(10.1.9.19:3306)/techmindDB" //env

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

        query := "SELECT name, COUNT(*) AS count FROM machines GROUP BY name"

        rows, err := db.Query(query)
        if err != nil {
            log.Fatal("Erro ao executar a consulta: ", err)
            return c.Next()
        }
        defer rows.Close()

        var systemCounts []CityCount
        
        for rows.Next() {
            var systemCount CityCount
            if err := rows.Scan(&systemCount.CityName, &systemCount.Count); err != nil {
                log.Fatal("Erro ao escanear a linha: ", err)
                return c.Next()
            }
            systemCounts = append(systemCounts, systemCount)
        }
        
        if err := rows.Err(); err != nil {
            log.Fatal("Erro ao iterar sobre os resultados: ", err)
            return c.Next()
        } 
        systemCounts = substituirValor(systemCounts, "lp00", "Caxias do Sul")
        systemCounts = substituirValor(systemCounts, "LP00", "Caxias do Sul")
        return c.JSON(fiber.Map{"status":200, "cities": systemCounts})        
    })

    app.Get("api/get-machines-days", func(c *fiber.Ctx) error {
        dataSourceName := "mach:Lup@.CSC.!@tcp(10.1.9.19:3306)/techmindDB" //env
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

        return c.Next()
    })

    app.Get("*", func(c *fiber.Ctx) error {
        return c.SendFile("./build/browser/index.html")
    })

    app.Use(func(c *fiber.Ctx) error{
        return c.SendFile("./build/browser/index.html") //variavel de ambiente
    })

    if err := app.Listen(":3000"); err != nil{
        log.Fatalf("Erro ao iniciar o Servidor %v", err)
    }
}