package actions

import (
	"api/locales"
	"sync"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/contenttype"
	"github.com/gobuffalo/middleware/forcessl"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"
	"github.com/gobuffalo/x/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

var ENV = envy.Get("GO_ENV", "development")

var (
	app     *buffalo.App
	appOnce sync.Once
	T       *i18n.Translator
)

func init() {
	goth.UseProviders(
		github.New(
			envy.Get("GITHUB_ClientID", ""),
			envy.Get("GITHUB_ClientSecrets", ""),
			envy.Get("GITHUB_CALLBACK_URL", "http://localhost:3000/auth/github/callback"),
		),
	)
}

func App() *buffalo.App {
	appOnce.Do(func() {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: envy.Get("SESSION_SECRET", "raccoon-mh"),
		})

		// app.Use(forceSSL())
		app.Use(paramlogger.ParameterLogger)
		app.Use(contenttype.Set("application/json"))

		auth := app.Group("/auth")
		auth.GET("/{provider}/callback", AuthCallback)
		auth.GET("/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))

		api := app.Group("/api")
		// api.Use(AuthMiddleware)
		api.GET("/{targetController}", GetRouteController)
		api.POST("/{targetController}", PostRouteController)

	})

	return app
}

func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
