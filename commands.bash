# List of commands


# ------------------
# Workspace
# ------------------
# https://go.dev/doc/tutorial/workspaces

go work init ./hello
# Add module to workspace
go work use ./example/hello
# syncs dependencies from the workspace’s build list into each of the workspace modules.
go work sync

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