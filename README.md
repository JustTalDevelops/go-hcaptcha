# go-hcaptcha

A Go library for solving hCaptchas with YOLOv3 and options for other image recognition systems within twenty seconds.

![Example of the solver in action with YOLOv3.](images/example.png)

## Installation
In order to use the main YOLOv3 solver, you'll need to install [gocv](https://github.com/hybridgroup/gocv), 
and the YOLOv3 config, weights, and names, which can be downloaded from the 
[Go YOLOv3 repository](https://github.com/wimspaargaren/yolov3).

Save all three files in a `yolo` directory in the main directory of your project. `go-hcaptcha` will use this directory
to load the YOLOv3 model.

## Basic Usage
In order to solve, you need the site URL (not the domain!), and the site key, which can be found 
in the HTML of the website with the hCaptcha challenge.

Below is a basic example of how to use the solver with the two using YOLOv3.
```go
c, err := NewChallenge(siteUrl, siteKey)
if err != nil {
    panic(err)
}
err = c.Solve(&YOLOSolver{Log: c.log})
if err != nil {
    c.log.Panic(err)
}
c.log.Info(c.Token())
```

## Credits

### 2.0.0
The motion data capturing required with hCaptcha would not be possible without the work of 
[@h0nde](https://github.com/h0nde) and his [py-hcaptcha](https://github.com/h0nde/py-hcaptcha) solver in Python.

### 1.0.2:
There were quite a lot of changes with the hCaptcha API, so the solver was updated to reflect these changes, with
the generous help of [@aw1875](https://github.com/aw1875) and his [puppeteer-hcaptcha](https://github.com/aw1875/puppeteer-hcaptcha) 
solver in JavaScript.

### 1.0.0
This project was inspired by the work of [@JimmyLaurent](https://github.com/JimmyLaurent) and his [hcaptcha-solver](https://github.com/JimmyLaurent/hcaptcha-solver)
also in JavaScript. I'd like to thank him for his work, and for being a motivation to create this library.