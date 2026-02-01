# Mage Examples

## Basic Magefile

```go
//go:build mage

package main

import (
	"fmt"
	"os"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = Build

// Build runs the build steps
func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	return sh.Run("go", "build", "-o", "myapp", "./...")
}

// InstallDeps installs dependencies
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	return sh.Run("go", "mod", "download")
}

// Clean removes build artifacts
func Clean() {
	fmt.Println("Cleaning...")
	os.Remove("myapp")
}
```

## Using Context and Timeouts

```go
//go:build mage

package main

import (
	"context"
	"time"
	"github.com/magefile/mage/mg"
)

// Deploy simulates a deployment with a timeout
func Deploy(ctx context.Context) error {
	// Pass context to dependencies
	mg.CtxDeps(ctx, Build)
	
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(2 * time.Second):
		println("Deployed!")
		return nil
	}
}

func Build() {
    println("Building...")
}
```

## Namespaces

```go
//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
)

type Docker mg.Namespace

// Build builds the docker image
// usage: mage docker:build
func (Docker) Build() error {
	return nil
}

// Push pushes the docker image
// usage: mage docker:push
func (Docker) Push() error {
	return nil
}
```
