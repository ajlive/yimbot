package main

// Import the driver
import (
	"fmt"
	"os"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

type Address struct {
	Street  string `fauna:"street"`
	City    string `fauna:"city"`
	State   string `fauna:"state"`
	ZipCode string `fauna:"zipCode"`
}
type Store struct {
	Name    string  `fauna:"name"`
	Address Address `fauna:"address"`
}

type Product struct {
	Name           string  `fauna:"name"`
	Description    string  `fauna:"description"`
	Price          float64 `fauna:"price"`
	Quantity       int64   `fauna:"quantity"`
	Store          f.RefV  `fauna:"store"`
	BackorderLimit int64   `fauna:"backorderLimit"`
	Backordered    bool    `fauna:"backordered"`
}

func run() error {
	zipcode, err := zipCodeFQL()
	if err != nil {
		return err
	}
	fmt.Printf("found zipcode: %#v\n", zipcode)
	return nil
}

func zipCodeFQL() (zipcode string, err error) {
	// set up client
	secret := os.Getenv("FAUNADB_SECRET")
	if secret == "" {
		return "", fmt.Errorf("The FAUNADB_SECRET environment variable is not set, exiting.")
	}
	endpoint := os.Getenv("FAUNADB_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://db.us.fauna.com/"
	}
	client := f.NewFaunaClient(secret, f.Endpoint(endpoint))

	// get product
	var product Product
	res, err := client.Query(f.Get(f.Ref(f.Collection("products"), "201")))
	if err != nil {
		return "", err
	}
	if err = res.At(f.ObjKey("data")).Get(&product); err != nil {
		return "", err
	}
	fmt.Printf("found product: %#v\n", product)

	// get store
	var store Store
	res, err = client.Query(f.Get(product.Store))
	if err != nil {
		return "", err
	}
	if err = res.At(f.ObjKey("data")).Get(&store); err != nil {
		return "", err
	}
	fmt.Printf("found store: %#v\n", store)

	// get zipcode
	zipcode = store.Address.ZipCode
	return zipcode, nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
