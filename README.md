# Arche Serde

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/arche-serde/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/arche-serde/actions/workflows/tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/mlange-42/arche-serde/badge.svg?branch=main)](https://coveralls.io/github/mlange-42/arche-serde?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/arche-serde)](https://goreportcard.com/report/github.com/mlange-42/arche-serde)
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/arche-serde.svg)](https://pkg.go.dev/github.com/mlange-42/arche-serde)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/arche-serde)
[![MIT license](https://img.shields.io/github/license/mlange-42/arche-serde)](https://github.com/mlange-42/arche-serde/blob/main/LICENSE)

*Arche Serde* provides JSON serialization and deserialization for the [Arche](https://github.com/mlange-42/arche) Entity Component System (ECS).

<div align="center">

<a href="https://github.com/mlange-42/arche">
<img src="https://user-images.githubusercontent.com/44003176/236701164-28178d13-7e52-4449-baa4-41b764183cbd.png" alt="Arche (logo)" width="500px" />
</a>

</div>

## Features

* Serialize/deserialize an entire *Arche* world in one line.
* Proper serialization of entity relations, as well as of entities stored in components.

## Installation

```
go get github.com/mlange-42/arche-serde
```

## Usage

See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche-serde) for more details and examples.  
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/arche-serde.svg)](https://pkg.go.dev/github.com/mlange-42/arche-serde)

Serialize a world:

```go
jsonData, err := archeserde.Serialize(&world)
if err != nil {
    // handle the error
}
```

Deserialize a world:

```go
err = archeserde.Deserialize(jsonData, &world)
if err != nil {
    // handle the error
}
```

## License

This project is distributed under the [MIT licence](./LICENSE).
