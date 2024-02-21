package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func main() {
    app := fiber.New()

    app.Static("/", "./build/browser")   

    app.Post("api/credential", func(c *fiber.Ctx) error {
        var data map[string]string

        if err := c.BodyParser(&data); err != nil {
            return err
        }

        user := data["user"]

        username := fmt.Sprintf("nt-lupatech\\%s", user)
        pass := data["pass"]

        l, err := ldap.DialURL("ldap://sdc01.nt-lupatech.com.br")
        if err != nil {
            log.Warn(err)
        }

        defer l.Close()

        err = l.Bind(username, pass)
        if err != nil {
            log.Warn(err)
        }

        Filter := fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", user)

        searchRequest := ldap.NewSearchRequest(
            "ou=Brasil,dc=nt-lupatech,dc=com,dc=br",
            ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0,
            false,
            Filter,
            []string{"sAMAccountName", "givenName", "memberOf"},
            nil,
        )

        searchResult, err := l.Search(searchRequest)
        if err != nil {
            log.Warn(err)
        }

        if err:= godotenv.Load("../.env"); err != nil {
            log.Warn(err)
        }

        verify := os.Getenv("VERIFY_GROUP")        

        for _, entry := range searchResult.Entries {
            for _, attr := range entry.Attributes {
                if attr.Name == "memberOf" {
                    groups := attr.Values

                    for _, value := range groups{
                        if !strings.Contains(value, verify){
                            return c.JSON(fiber.Map{"status":401})
                        } else {
                            log.Info("a")
                        }
                    }
                }
            }
        }
    
        return c.JSON(fiber.Map{"status":"ok"})
    })

    if err := app.Listen(":3000"); err != nil{
        log.Fatalf("Erro ao iniciar o Servidor %v", err)
    }
}