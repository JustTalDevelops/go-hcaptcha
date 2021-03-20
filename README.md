# hcaptcha-solver-go

An HCaptcha solver for GoLang.
It uses Playwright for generating HSW which is put in an HSW pool, which can then be accessed by a solver to solve captchas.

## Feature: Vision API
Feature: Vision API uses the Google Vision API to solve captchas.
It uses the object that we need to find, along with every image in the captcha's
task list to find the ones that match.

This would replace the random guessing that we do, and make captcha solving
much faster. Of course, there would still be an option to use the random guessing
if you aren't willing to use Vision API.

![An example of the Vision API with one of HCaptcha's images.](./images/bus_example.png)

# Usage
Below are some usage examples on how you would use the solver.

## No proxies, one worker:

```go
s, err := hcaptcha.NewSolver("example.com")
if err != nil {
  panic(err)
}
defer s.Close()
// We provide a deadline that the solver must have the solution done by.
// If the deadline is not reached, an error is sent instead of the solution.
solution, err := s.Solve(time.Now().Add(50 * time.Second))
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
// We provide a deadline that the solver must have the solution done by.
// If the deadline is not reached, an error is sent instead of the solution.
solution, err := s.Solve(time.Now().Add(50 * time.Second))
if err != nil {
  panic(err)
}
// F0_eyJ0eXAiOiJKV1Q...
fmt.Println(solution)
```

# Resources

## [JimmyLaurent's HCaptcha Solver in JS](https://github.com/JimmyLaurent/hcaptcha-solver)
JimmyLaurent's helped a lot with the core structure of HCaptcha's API.
If you are interested in building your own solver, I would check out his repository,
although it is a bit outdated.

## [aw1875's HCaptcha Solver using Puppeteer](https://github.com/JimmyLaurent/hcaptcha-solver)
aw1875 was a big help with the issues I was encountering with my own implementation.
His implementation is much more up to date then Jimmy Laurent's, so if you're 
a JS developer, I would recommend using his work.
