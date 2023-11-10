package main

import (
	"Vectra/src/controller"
	"Vectra/src/model/i18n"
	. "Vectra/src/model/service"
	. "Vectra/src/model/storage"
	view "Vectra/src/view/go"
	"bytes"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"log"
	"time"
)

type Host struct {
	Fiber *fiber.App
}

func main() {

	err := i18n.GetInstance().SetUp(GetApiV1().GetStore().Config.Langs...)
	if err != nil {
		panic(err)
	}

	isDev := GetApiV1().GetStore().Config.IsDev
	currentDomain := getDomainName(isDev)
	hosts := map[string]*Host{}
	store := session.New(session.Config{
		KeyLookup:      "cookie:" + CookieNameForSession,
		CookieSecure:   !isDev,
		CookieHTTPOnly: true,
	})

	// Static
	makeStatic(store, hosts, currentDomain)

	// Website
	makeWebsite(isDev, store, hosts, currentDomain)

	// Start the app
	log.Fatal(createApp(hosts).Listen(currentDomain))
}

func makeStatic(store *session.Store, hosts map[string]*Host, currentDomain string) {
	static := fiber.New()
	static.Use(cors.New())
	static.Use(compress.New())
	static.Use(logger.New(logger.Config{Format: "STATIC [${ip}]:${port} ${status} - ${method} ${path}\n"}))
	static.Use(cache.New(cache.Config{
		Expiration:   time.Hour * 3 * 24,
		Storage:      store.Storage,
		CacheControl: true,
	}))
	static.Get("sprite", func(ctx *fiber.Ctx) error {
		var buf = new(bytes.Buffer)
		view.Jade_sprite(buf)
		return ctx.Send(buf.Bytes())
	})
	static.Static("/", "./static", fiber.Static{
		Compress: true,
	})
	hosts["static."+currentDomain] = &Host{static}
}

func makeWebsite(isDev bool, store *session.Store, hosts map[string]*Host, currentDomain string) {

	firstLaunchHandler := func(ctx *fiber.Ctx) error {
		isFirstLaunch := GetApiV1().IsFirstLaunch()
		if ctx.Method() == "GET" && ctx.Path() != "/init" && isFirstLaunch {
			return ctx.Redirect("/init")
		}
		if ctx.Method() == "GET" && ctx.Path() == "/init" && !isFirstLaunch {
			return ctx.Redirect("/", fiber.StatusServiceUnavailable)
		}
		return ctx.Next()
	}

	fillDefaultSession := func(ctx *fiber.Ctx) error {
		sess, err := store.Get(ctx)
		if err != nil {
			return err
		}
		if sess.Get(SessionKeyForUserId) == nil {
			sess.Set(SessionKeyForUserId, "")
		}
		sess.Save()
		return ctx.Next()
	}
	csrfHandler := csrf.New(csrf.Config{
		CookieName:     CookieNameForCSRF,
		CookieSecure:   !isDev,
		Expiration:     time.Minute * 30,
		CookieSameSite: "Strict",
		Storage:        store.Storage,
	})

	site := fiber.New()
	site.Use(logger.New(logger.Config{Format: "[${ip}]:${port} ${status} - ${method} ${path}\n"}))
	site.Use(compress.New())
	site.Use(favicon.New(favicon.Config{File: "./static/favicon.ico", URL: "/favicon.ico"}))
	site.Use(firstLaunchHandler, fillDefaultSession, csrfHandler)
	hosts[currentDomain] = &Host{site}

	controller.NewApiV1Controller(site.Group("/api/v1"), store)
	controller.NewViewController(site.Group("/"), store)
}

func createApp(hosts map[string]*Host) *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
			if err != nil {

				if ctx.Method() == "POST" {
					return ctx.Status(fiber.StatusInternalServerError).JSON(
						ReasonExch{Reason: "Unknown"},
					)
				}
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}
			return nil
		},
	})
	app.Use(recover.New())

	app.Use(func(c *fiber.Ctx) error {
		host := hosts[c.Hostname()]
		if host == nil {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			host.Fiber.Handler()(c.Context())
			return nil
		}
	})

	return app
}

func getDomainName(isDev bool) string {
	if isDev {
		return GetApiV1().GetStore().Config.DevDomain
	}
	return GetApiV1().GetStore().Config.Domain
}
