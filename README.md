# hcaptcha-solver-go

An HCaptcha solver for GoLang.
It uses a Playwright tab with the HSW script injected to generate HSW
on the fly, and then pipes that to other go routines which are made to purely
solve HCaptchas until it gets a password UUID.

# Usage

## No proxies, one worker:

```go
s, err := hcaptcha.NewSolver("example.com", 1)
if err != nil {
  panic(err)
}
defer s.Close()
solution, err := s.Solve()
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
solution, err := s.Solve()
if err != nil {
  panic(err)
}
// F0_eyJ0eXAiOiJKV1Q...
fmt.Println(solution)
```

# Resources

This project is not fully complete. If you are interested in contributing,
I would recommend you check out JimmyLaurent's HCaptcha solver in JavaScript,
and aw1875's HCaptcha solver (which is more up to date) and uses Puppeteer.

# Disclaimer

This project is not complete however contributions are welcome.
