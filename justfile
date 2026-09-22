set dotenv-load
oracle:
    go run ./cmd/oracle
tictactoe:
    @go build ./cmd/tictactoe
    @./tictactoe
