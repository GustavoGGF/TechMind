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

func main() {
    app := fiber.New()

    app.Static("/", "./build/browser")   

    app.Use("/api", func(c *fiber.Ctx) error{
        return c.Next()
    })

    app.Post("api/credential", func(c *fiber.Ctx) error {
        var data map[string]string

        if err := c.BodyParser(&data); err != nil {
            return err
        }
    
        user := data["user"]

        username := fmt.Sprintf("nt-lupatech\\%s", user) //variavel de ambiente
        pass := data["pass"]

        l, err := ldap.DialURL("ldap://sdc01.nt-lupatech.com.br")
        if err != nil {
            fmt.Println(err)
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
            fmt.Println(err)
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
            return err
        }

        name := data["Name"]
        system := data["System"]

        fmt.Println("Antes de conectar")
        

        db, err := sql.Open("mysql", "mach:Lup@.CSC.!@tcp(10.1.9.0:3306)/techmindDB") //varaivel de ambiente
        if err != nil{
            fmt.Println("Erro conectar")
            fmt.Println(err)
            return c.Next()
        }

        fmt.Println("Depois de conectar")

        defer db.Close()

        err = db.Ping()
        if err != nil {
            fmt.Println(err)
            return c.Next()
        }

        stmt, err := db.Prepare("INSERT INTO machines(name, system_name) VALUES(?, ?)")
        if err != nil{
            fmt.Println(err)
            return c.Next()
        }

        defer stmt.Close()

        _, err = stmt.Exec(name, system)
        if err != nil {
            fmt.Println(err)
            return c.Next()
        }
        return c.Next()
    })


    app.Use(func(c *fiber.Ctx) error{
        return c.SendFile("./build/browser/index.html") //variavel de ambiente
    })

    if err := app.Listen(":3000"); err != nil{
        log.Fatalf("Erro ao iniciar o Servidor %v", err)
    }
}