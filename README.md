# hcaptcha-solver-go

An HCaptcha solver for GoLang.
It uses a Playwright tab with the HSW script injected to generate HSW
on the fly, and then pipes that to other go routines which are made to purely
solve HCaptchas until it gets a password UUID.

# Reason

This project is purely made to abuse on MCPE Pocket Servers (lmao fuck you),
so a lot of code still in development will be centered around that site.

# Resources

This project is not fully complete. If you are interested in contributing,
I would recommend you check out JimmyLaurent's HCaptcha solver in JavaScript,
and aw1875's HCaptcha solver (which is more up to date) and uses Puppeteer.

# Disclaimer

This project is not complete however contributions are welcome.