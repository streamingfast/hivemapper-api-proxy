# hivemapper-api-proxy

## Run main
```bash
go run ./cmd/api-proxy/main.go 
```
## How to use
1. Run the main with the above command.
2. If you want to switch members, you need to kill the proxy process and re-run it with the same command above

## Different users
1. `go run ./cmd/api-proxy/main.go` 
2. Connect with userA
3. Kill with Ctrl+C
4. `go run ./cmd/api-proxy/main.go` 
5. Refresh page which will prompt reconnection
6. Connect with userB 