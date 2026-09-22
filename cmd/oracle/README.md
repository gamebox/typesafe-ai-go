# typesafe-ai-go

A simple SDK to use [Typesafe AI's Jev](typesafe.ai) in your Go projects.

## Install

```sh
go install github.com/gamebox/typesafe-ai-go
```

## Usage

```go
state := "All dogs are brown.  Spot is a dog."
request := typesafe.NewRequest(state)
request.AddQuestion("brown", typesafe.Noul("Is Spot brown?"))

response, err := client.Ask(request)
if err != nil {
	fmt.Fprintf(os.Stderr, "Could not make request: %e\n", err)
	os.Exit(1)
}
a, ok := response.DecodeNoul("brown")
if a.Noul > 95 {
	fmt.Println("Spot is definitely brown.")
} else if a.Noul > 75 {
	fmt.Println("Spot is more than likely brown.")
} else if a.Noul > 50 {
	fmt.Println("Spot is probably brown.")
} else {
	fmt.Println("Spot is not brown.")
}
```

### API

Find more information on (Go Packages site)[pkg.go.dev/github.com/gamebox/typesafe-ai-go]
