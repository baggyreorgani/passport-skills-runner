# Passport Skills Runner

`passport-skills-runner` is a small, deterministic evaluator for an A2A agent's declared skill passport. It validates the parts that a caller needs before invoking a paid skill: protocol, version, input declarations, uniqueness and execution timeout.

```sh
go test ./...
go run . -manifest examples/passport.json
go run . -manifest examples/passport.json -json
```

The runner is intentionally local and side-effect free. It does not invoke skills or sign payments; it produces a pass/fail gate that can run before discovery or in CI.
