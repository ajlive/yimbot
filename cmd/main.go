package main

// Import the driver
import (
	"fmt"
	"os"
	"strings"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type HelloWorld struct {
	Greeting string `fauna:"greeting"`
	Name     string `fauna:"name"`
}

func run() error {
	secret := os.Getenv("FAUNADB_SECRET")
	if secret == "" {
		return fmt.Errorf("The FAUNADB_SECRET environment variable is not set, exiting.")
	}

	endpoint := os.Getenv("FAUNADB_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://db.us.fauna.com/"
	}

	client := f.NewFaunaClient(secret, f.Endpoint(endpoint))

	res, err := client.Query(f.Get(f.Ref(f.Collection("helloworld"), "313296063476793410")))
	if err != nil {
		return err
	}

	var hello HelloWorld
	if err := res.At(f.ObjKey("data")).Get(&hello); err != nil {
		return err
	}

	fmt.Printf("%v\n", strings.Join([]string{hello.Greeting, hello.Name}, ""))

	return nil
}
