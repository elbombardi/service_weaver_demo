package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	weaver "github.com/ServiceWeaver/weaver"
)

func main() {
	if err := weaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main] `weaver:"main"`
	reverser                       weaver.Ref[Reverser]
	hello                          weaver.Listener `weaver:"hello"`
}

func serve(context context.Context, app *app) error {
	fmt.Printf("hello listener is listening on %s\n", app.hello)
	http.HandleFunc("/hello", func(resp http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		if name == "" {
			name = "world"
		}
		reversed, err := app.reverser.Get().Reverse(context, name)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(resp, "Hello, %s!\n", reversed)
	})
	return http.Serve(app.hello, nil)
}
