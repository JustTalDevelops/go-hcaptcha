package hcaptcha

import (
	"github.com/justtaldevelops/go-hcaptcha/utils"
	"github.com/sirupsen/logrus"
	"github.com/wimspaargaren/yolov3"
	"gocv.io/x/gocv"
)

// Solver is an interface to solve an hCaptcha task.
type Solver interface {
	// Solve solves the hCaptcha task using the category, question, and the task. If it was successful,
	// it returns true, and in all other cases, it returns false.
	Solve(category, question string, task Task) bool
}

// GuessSolver solves hCaptcha tasks by guessing the solution.
type GuessSolver struct{}

// Solve ...
func (s *GuessSolver) Solve(_, _ string, _ Task) bool {
	return utils.Chance(0.5)
}

// yolo is the YOLO v3 network.
var yolo yolov3.Net

// init initializes the YOLO v3 network.
func init() {
	yolo, _ = yolov3.NewNet("yolo/yolov3.weights", "yolo/yolov3.cfg", "yolo/coco.names")
}

// YOLOSolver uses the "You Only Look Once" (YOLO) algorithm to solve hCaptcha tasks.
type YOLOSolver struct {
	// Log is the logger for the solver.
	Log *logrus.Logger
}

// Solve ...
func (s *YOLOSolver) Solve(category, object string, task Task) bool {
	// Make sure the YOLO network is initialized.
	if yolo == nil {
		panic("yolov3 data is not in expected folders")
	}

	// Make sure we can solve the challenge.
	if category != "image_label_binary" {
		s.Log.Debugf("cannot solve challenge with category %s", category)
		return false
	}

	// Decode and detect the object.
	frame, err := gocv.IMDecode(task.Image, gocv.IMReadColor)
	if err != nil {
		return false
	}

	detections, err := yolo.GetDetections(frame)
	if err != nil {
		return false
	}

	for _, detection := range detections {
		if detection.ClassName == object && detection.Confidence > 0.6 {
			s.Log.Debugf("Detected %v in provided image", detection.ClassName)
			return true
		}
	}
	return false
}

// Compile-time check that GuessSolver implements Solver.
var _ Solver = (*GuessSolver)(nil)
