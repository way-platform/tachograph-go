# Mage API Reference

## Package `mg` (`github.com/magefile/mage/mg`)

Contains the core logic for dependencies and running targets.

```go
const AnsiColorReset = "\033[0m"
const CacheEnv = "MAGEFILE_CACHE"
const DebugEnv = "MAGEFILE_DEBUG"
const EnableColorEnv = "MAGEFILE_ENABLE_COLOR"
const GoCmdEnv = "MAGEFILE_GOCMD"
const HashFastEnv = "MAGEFILE_HASHFAST"
const IgnoreDefaultEnv = "MAGEFILE_IGNOREDEFAULT"
const TargetColorEnv = "MAGEFILE_TARGET_COLOR"
const VerboseEnv = "MAGEFILE_VERBOSE"

func CacheDir() string
func CtxDeps(ctx context.Context, fns ...interface{})
func Debug() bool
func Deps(fns ...interface{})
func EnableColor() bool
func ExitStatus(err error) int
func Fatal(code int, args ...interface{}) error
func Fatalf(code int, format string, args ...interface{}) error
func GoCmd() string
func HashFast() bool
func IgnoreDefault() bool
func SerialCtxDeps(ctx context.Context, fns ...interface{})
func SerialDeps(fns ...interface{})
func TargetColor() string
func Verbose() bool

type Fn interface{ ... }
    func F(target interface{}, args ...interface{}) Fn
type Namespace struct{}
```

## Package `sh` (`github.com/magefile/mage/sh`)

Helper library for running shell commands.

```go
func CmdRan(err error) bool
func Copy(dst string, src string) error
func Exec(env map[string]string, stdout, stderr io.Writer, cmd string, args ...string) (ran bool, err error)
func ExitStatus(err error) int
func OutCmd(cmd string, args ...string) func(args ...string) (string, error)
func Output(cmd string, args ...string) (string, error)
func OutputWith(env map[string]string, cmd string, args ...string) (string, error)
func Rm(path string) error
func Run(cmd string, args ...string) error
func RunCmd(cmd string, args ...string) func(args ...string) error
func RunV(cmd string, args ...string) error
func RunWith(env map[string]string, cmd string, args ...string) error
func RunWithV(env map[string]string, cmd string, args ...string) error
```

## Package `target` (`github.com/magefile/mage/target`)

Helper library for determining if targets are out of date.

```go
func Dir(dst string, sources ...string) (bool, error)
func DirNewer(target time.Time, sources ...string) (bool, error)
func Glob(dst string, globs ...string) (bool, error)
func GlobNewer(target time.Time, sources ...string) (bool, error)
func NewestModTime(targets ...string) (time.Time, error)
func OldestModTime(targets ...string) (time.Time, error)
func Path(dst string, sources ...string) (bool, error)
func PathNewer(target time.Time, sources ...string) (bool, error)
```
