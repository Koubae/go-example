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