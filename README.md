# hcaptcha-solver-go

An HCaptcha solver for GoLang.
It uses a Playwright tab with the HSW script injected to generate HSW
on the fly, and then pipes that to other go routines which are made to purely
solve HCaptchas until it gets a password UUID.

# Usage
Below are some usage examples on how you would use the solver.

## No proxies, one worker:

```go
s, err := hcaptcha.NewSolver("example.com")
if err != nil {
  panic(err)
}
defer s.Close()
solution, err := s.Solve(time.Now().Add(50 * time.Second)) // We must provide a deadline
if err != nil {
panic(err)
}
// F0_eyJ0eXAiOiJKV1Q...
fmt.Println(solution)
```

## Proxied, two workers:

```go
s, err := hcaptcha.NewSolverWithProxies("example.com", 2, proxies)
if err != nil {
  panic(err)
}
defer s.Close()
solution, err := s.Solve(time.Now().Add(50 * time.Second)) // We must provide a deadline
if err != nil {
  panic(err)
}
// F0_eyJ0eXAiOiJKV1Q...
fmt.Println(solution)
```

# Resources

## [JimmyLaurent's HCaptcha Solver in JS](https://github.com/JimmyLaurent/hcaptcha-solver)
JimmyLaurent's solver was a big help with the core structure of HCaptcha's API.
If you are interested in building your own solver, I would check out his repository,
although it is a bit outdated.

## [aw1875's HCaptcha Solver using Puppeteer](https://github.com/JimmyLaurent/hcaptcha-solver)
aw1875 was a big help with the issues I was encountering with my own implementation.
His implementation is much more up to date then Jimmy Laurent's, so if you're 
a JS developer, I would recommend using his work.