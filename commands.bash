# List of commands
go get gopkg.in/yaml.v3

# ------------------
# Workspace
# ------------------
# https://go.dev/doc/tutorial/workspaces

go work init ./hello
# Add module to workspace
go work use ./example/hello
# syncs dependencies from the workspace’s build list into each of the workspace modules.
go work sync

go work use ./workspace/io/filedb
go work use ./workspace/dir1/dir2/dir3/mymodule

# https://go.dev/ref/mod#environment-variables
GOWORK=off go run main.go

# ------------------
# Modules
# ------------------
go mod init github.com/koubae/go-example

mkdir hello; cd hello;
go mod init github.com/koubae/go-example/hello

mkdir cmd pkg internal examples scripts

cd examples
mkdir http; cd http;
mkdir server-1 server-2

cd server-1
go mod init github.com/koubae/go-example/examples/http/server-1

cd server-2
go mod init github.com/koubae/go-example/examples/http/server-2