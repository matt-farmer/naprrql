// sifql-server.go
//
// simple web server to support gql queries & web ui (graphiql)
//
package naprrql

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/playlyfe/go-graphql"
)

var executor *graphql.Executor

//
// wrapper type to capture graphql input
//
type GQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

//
// the core graphql handler routine
//
func graphQLHandler(c echo.Context) error {

	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return err
	}

	query := grq.Query
	variables := grq.Variables
	gqlContext := map[string]interface{}{}

	result, err := executor.Execute(gqlContext, query, variables, "")
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, result)

}

//
// launches the server
//
func RunQLServer() {

	executor = buildExecutor()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// routes to html/css/javascript resources
	e.Static("/", "public")
	e.File("/sifql", "public/ql_index.html")
	e.File("/ui", "public/ui_index.html")

	// the graphql handler
	e.POST("/graphql", graphQLHandler)

	fmt.Printf("\n\nBrowse to follwing locations:\n")
	fmt.Printf("\n\thttp://localhost:1329/ui\n\n for qa report user interface\n")
	fmt.Printf("\n\thttp://localhost:1329/sifql\n\n for data explorer\n\n")

	e.Logger.Fatal(e.Start(":1329"))
}
