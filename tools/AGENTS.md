# Development Tools

This document provides guidance on how the build scripts and build targets work.

## Developing

- The project uses a [tools](./) directory with a separate Go module for build, lint, and generation tools.
- Build tasks are defined in [magefile.go](./magefile.go) using Mage.
- Run tests: `./tools/mage test`
- Lint code: `./tools/mage lint`
- Re-generate code: `./tools/mage generate`
