package main

import (
	"github.com/coreos/bbolt"
	"log"
	"io/ioutil"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"net/http"
)

const (
	BOLTPATH      = "bolt.db"
	WEBSITEBUCKET = "Websites"
)

var writeChannel chan Website
var writeUserChannel chan User
var HugsCacher *HugsCache
var db *bolt.DB

func init() {
	writeChannel = make(chan Website)
	writeUserChannel = make(chan User)
}

func main() {
	fmt.Println("starting")
	database, err := bolt.Open(BOLTPATH, 0600, nil)
	err = database.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(WEBSITEBUCKET))
		return nil
	})
	if err != nil {
		log.Fatal("fuck")
	}
	db = database
	HugsCacher = HugsCache{}.New()

	go WriteDatabaseThread()
	go CacheReloader()

	defer db.Close()



	schema := graphql.MustParseSchema(getSchema(), &Resolver{})
	//for shiggles
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	http.Handle("/graphql", &relay.Handler{Schema: schema})

	fmt.Println("Ready to serve")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

//the shiggle page
var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css" rel="stylesheet" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/es6-promise/4.1.1/es6-promise.auto.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/graphql", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)

func getSchema() string {
	//you can compose or whatever else you want here, i prefer one larger schema
	file, err := ioutil.ReadFile("./schema.graphql")
	if err != nil {
		log.Fatal("schema does not exist")
	}
	return string(file)
}
