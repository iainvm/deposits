# Deposits

This repo is to house an interview coding test for Touch

## Project Dependencies

Nix (specifically nix flakes) handles the majority of project dependencies ensuring they're versioned
If you already have these dependencies in your system there is no need to use them


## Project Layout

```
/cmd/grpc/main.go   # Creates the binary for hosting the server
/gen/...            # Any generated files through protoc
/protos/...         # The proto specs are stored here
/internal/...       # Contains all the domain business logic

```
