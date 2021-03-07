# pprof

1. `go get github.com/pkg/profile`
2. embed `defer profile.Start(profile.ProfilePath(".")).Stop()` into `main.Run`
3. `make build`
4. `build/darwin-amd64/worldwide XXX.gb`
5. `go tool pprof build/darwin-amd64/worldwide cpu.pprof`
6. `go tool pprof -png build/darwin-amd64/worldwide cpu.pprof > pprof.png`