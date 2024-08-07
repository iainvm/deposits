# Deposits

This repo is to house an interview coding test for Touch

## Project Dependencies

Nix (specifically nix flakes) handles the majority of project dependencies ensuring they're versioned
If you already have these dependencies in your system there is no need to use them


## Project Layout

```
/application/grpc/protos    # The proto specs are stored here
/application/grpc/gen       # Any generated files through protoc
/application/grpc/main.go   # Creates the binary for hosting the gRPC server
/internal                   # Contains all the domain business logic
/common                     # Stores packages that are agnostic to this product
/infrastructure             # Stored what's needed to build and host the server locally
```

## Host

If you have Taskfile installed then a simple `task docker` will run the `docker-compose.yaml` and bring up the `api` and `postrgres` contianers

`docker compose -f infrastructure/docker-compose.yml up -d --build` ran in the root of the project will work as well
