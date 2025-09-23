go-example
==========


_Go Example Workspace with multiple GoLang Recipes from simple one, http servers and more_

* [Programming-CookBook (GO)](https://github.com/Koubae/Programming-CookBook/tree/master/Programming%20Languages/GO)
* [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
* [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests)

This is meant to represent your local Go workspace. 
So you would not commit this to source but rather content inside [workspace](./workspace) would be separate
Git repositories/sources; this project is just to show what a Go workspace may look like.

Also, of course, contains actual simple Go / Started examples. 

Check also the official: **[Go example projects (github.com/golang)](https://github.com/golang/example)**


### Quickstart


```bash
make init
```

Workspace
---------

As write above, the [go.work](go.work) is **intentionally** committed to source; this is an example project on how
a Go Workspace may be managed.

Add a new module to the workspace:

```bash
make workspace-add MODULE=path/to/module
```

You should commit `./workspace/` part.


### Note about workspace

I stabled upon a "weird" and "unexpected behavior" while working on a workspace, well is unexpected till you know
 it, then it becomes... "expected" I guess.

Once you create a workspace and have a [go.work](go.work) in some root directory... then everything from the same
level of "root" and all its "children/subdirectory" are "part of Go Workspace"...
You **can't escape it**.

If you have a Go Workspace in `~/workspace_go/go.work` and create a Go Module in `~/workspace/dir1/dir2/dir3/mymodule/go.mod`
go will know that module is part of a workspace. So if you attempt to create a package in it and import, **it will not work**
unless you specifically add that into your `go.work`.

Example:

* See Example: [mymodule](./workspace/dir1/dir2/dir3/mymodule)

You create this folder `~/workspace/dir1/dir2/dir3/mymodule` and `cd` in it.
Create a Go Module called `mymodule` and a file `main.go`.
Then you create a package called `utils` and have `utils.go` in it

```bash
~/workspace_go/dir1/dir2/dir3/mymodule
    |
      - go.mod 
      - main.go # <-- entry point.. all good
      /utils
          |
            - utils.go
```

So you have a package named `package utils` and you try to import it in `main.go` it will fail...


```bash
go run main.go
main.go:4:2: package mymodule/utils is not in std (/usr/local/go/src/mymodule/utils)
```

### Solution

So you can either:

1. Add Module `my-module` in `go.work` by adding it with this command `go work use ./workspace/dir1/dir2/dir3/mymodule`
2. Or you turn off `GOWORK`:

```bash
GOWORK=off go run main.go
```


The "official" documentation and how I understood this behavior was expected I to follow those links/doc in order:

* [go build not working when used with workspaces](https://stackoverflow.com/a/76180815/13903942)
  * [Minimal version selection (MVS)](https://go.dev/ref/mod#minimal-version-selection)
  * [Workspace](https://go.dev/ref/mod#workspaces)
  * [Environment variables](https://go.dev/ref/mod#environment-variables) 


### Honest Opinion

While I find the workspace neat and pretty cool, this behavior is... "weird"... I may see one day if I can create a
Issue ticket to Go... or even create a pool request and change this or at least, have the workspace for a specific 
module **turned off** inside the `go.mod`... and not just an env variable (`GOWORK`)


Cool Stuff
=========

### Go Hot Reloading

* [air-verse/air](https://github.com/air-verse/air)

```bash
go install github.com/air-verse/air@latest

# Make sure that GOPATH and GOROOT is in your PATH
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin

# initialize
air init
# run
air
air server --port 8080
# Will run ./tmp/main -h
air -- -h

# Will run air with custom config and pass -h argument to the built binary
air -c .air.toml -- -h
```


* [Go Examples](https://github.dev/gin-gonic/examples)


Courses & Learning 
=======

### Coursera

* [Programming with Google Go Specialization](https://www.coursera.org/specializations/google-golang)
* [C, Go, and C++: A Comprehensive Introduction to Programming Specialization](https://www.coursera.org/specializations/c-go-c-plus-plus)
* [Developing a REST API with Go and Cloud Run](https://www.coursera.org/projects/googlecloud-developing-a-rest-api-with-go-and-cloud-run-tdech)
* [Golang for Beginners: Data Types, Functions, and Packages](https://www.coursera.org/projects/golang-beginners-data-types-functions-packages)


### Youtube

* [Complete Backend Engineering Course in Go](https://www.youtube.com/watch?v=h3fqD6IprIA&t=591s)
    * [sikozonpc/GopherSocial](https://github.com/sikozonpc/GopherSocial)
    * [Self Made Engineer](https://www.selfmadeengineer.com/)
