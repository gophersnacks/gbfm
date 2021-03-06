package snacks

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gophersnacks/gbfm/pkg/render"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gophersnacks/gbfm/pkg/web"
)

var r = render.New("snacks/application.html", map[string]interface{}{})

// App is where all routes and middleware for gophersnacks.com are defined
func App() *buffalo.App {
	app := buffalo.New(buffalo.Options{
		Addr:        "0.0.0.0:3000",
		Env:         web.ENV,
		SessionName: "_snacks_session",
	})
	// Automatically redirect to SSL
	app.Use(ssl.ForceSSL(secure.Options{
		SSLRedirect:     web.ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	}))

	if web.ENV == "development" {
		app.Use(middleware.ParameterLogger)
	}

	// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
	// Remove to disable this.
	app.Use(csrf.New)

	app.Use(web.LayoutMiddleware(
		"GopherSnacks",
		"snacks/partials/nav.html",
		"snacks/partials/footer.html",
	))

	// TODO: create gorm transaction wrapper middleware
	//
	// Wraps each request in a transaction.
	//  c.Value("tx").(*pop.Connection)
	// Remove to disable this.
	//app.Use(middleware.PopTransaction(models.DB))

	app.Use(web.Translator.Middleware())

	app.GET("/", homeHandler)
	app.GET("/snack/{snack_slug}", snackHandler)

	app.GET("/topics", topicHandler)
	app.GET("/topic/{topic_slug}", topicHandler)
	app.ServeFiles("/", render.AssetsBox) // serve files from the public directory

	return app
}
