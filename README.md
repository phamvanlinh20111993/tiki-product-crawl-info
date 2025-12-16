
# When face error
WARNING: undefined behavior - version of Delve is too old for Go version go1.25.4 (maximum supported version 1.22)
=> fix: open cmd: go install github.com/go-delve/delve/cmd/dlv@latest

# How to build + run go program
    1. go build main.go
    2. main.exe

# go update all modules: https://stackoverflow.com/questions/67201708/go-update-all-modules

# go project structure: https://go.dev/doc/code#ImportingLocal

# Command with go in cmd: 
    - go get <dependencies>
    - go mod tidy
    - go get -u
    - go list -m -u all
    - go test
    - go get -u && go mod tidy 
    - go clean -modcache
    - ...