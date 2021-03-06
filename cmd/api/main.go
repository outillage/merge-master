package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aevea/knit/api"
	"github.com/aevea/knit/api/generated"
	cfg "github.com/aevea/knit/internal/config"
	"github.com/pacedotdev/oto/otohttp"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func main() {
	cfg.InitEnv()

	if !viper.IsSet("GITHUB_TOKEN") {
		log.Println("Missing GITHUB_TOKEN requests to Github will fail")
	}

	githubToken := viper.GetString("GITHUB_TOKEN")

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	pullRequestService := api.PullRequestService{HTTPClient: httpClient}

	g := api.HealthcheckService{}

	server := otohttp.NewServer()
	generated.RegisterHealthcheck(server, g)
	generated.RegisterPullRequestService(server, pullRequestService)
	http.Handle("/oto/", server)

	fs := http.FileServer(http.Dir("./webapp/public"))
	http.Handle("/", fs)

	log.Println("Server starting on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
