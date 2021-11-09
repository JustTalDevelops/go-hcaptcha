package hcaptcha

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"github.com/go-redis/redis/v8"
	"github.com/justtaldevelops/go-hcaptcha/utils"
	"github.com/sirupsen/logrus"
	"github.com/wimspaargaren/yolov3"
	"gocv.io/x/gocv"
)

// Solver is an interface to solve hCaptcha tasks.
type Solver interface {
	// Solve solves the hCaptcha tasks using the category, question, and the task. If it was successful,
	// it returns true, and in all other cases, it returns false.
	Solve(category, question string, tasks []Task) []Task
}

// GuessSolver solves hCaptcha tasks by guessing the solution.
type GuessSolver struct{}

// Solve ...
func (s *GuessSolver) Solve(_, _ string, tasks []Task) (answers []Task) {
	for _, task := range tasks {
		if utils.Chance(0.5) {
			answers = append(answers, task)
		}
	}
	return answers
}

// yolo is the YOLO v3 network.
var yolo yolov3.Net

// init initializes the YOLO v3 network.
func init() {
	yolo, _ = yolov3.NewNet("yolo/yolov3.weights", "yolo/yolov3.cfg", "yolo/coco.names")
}

// YOLOSolver uses the "You Only Look Once" (YOLO) algorithm to solve hCaptcha tasks.
type YOLOSolver struct {
	// Redis is the redis client used to store the hCaptcha task data.
	Redis *redis.Client
	// Log is the logger for the solver.
	Log *logrus.Logger
}

// Solve ...
func (s *YOLOSolver) Solve(category, object string, tasks []Task) []Task {
	// Make sure the YOLO network is initialized.
	if yolo == nil {
		panic("yolov3 data is not in expected folders")
	}

	// Make sure we can solve the challenge.
	if category != "image_label_binary" {
		s.Log.Debugf("cannot solve challenge with category %s", category)
		return []Task{}
	}

	// Answer the challenge.
	var answers []Task
	for _, task := range tasks {
		h := sha1.New()
		h.Write(task.Image)
		hash := hex.EncodeToString(h.Sum(nil))
		customId := object + "|" + hash
		score := s.tileScore(customId)
		if score < 0 {
			// Impossible.
			continue
		}

		// Possible!
		if score >= 1 {
			s.Log.Debugf("Detected %v in provided image (from cache!)", object)

			answers = append(answers, task)
			_ = s.increaseTileScore(customId, 1)
			continue
		}

		// Decode and detect the object.
		frame, err := gocv.IMDecode(task.Image, gocv.IMReadColor)
		if err != nil {
			continue
		}

		detections, err := yolo.GetDetections(frame)
		if err != nil {
			continue
		}

		var detected bool
		for _, detection := range detections {
			if detection.ClassName == object && detection.Confidence > 0.6 {
				s.Log.Debugf("Detected %v in provided image", object)

				answers = append(answers, task)
				_ = s.increaseTileScore(customId, 1)
				detected = true
				break
			}
		}

		if !detected {
			_ = s.decreaseScore(customId, 1)
		}
	}

	return answers
}

// decreaseScore sets the impossible flag for the given task.
func (s *YOLOSolver) decreaseScore(id string, delta int) error {
	if s.Redis == nil {
		return nil
	}
	return s.Redis.DecrBy(context.Background(), id, int64(delta)).Err()
}

// increaseTileScore increases the score of the tile.
func (s *YOLOSolver) increaseTileScore(id string, delta int) error {
	if s.Redis == nil {
		return nil
	}
	return s.Redis.IncrBy(context.Background(), id, int64(delta)).Err()
}

// tileScore is the score of a tile.
func (s *YOLOSolver) tileScore(id string) int {
	if s.Redis == nil {
		return 0
	}

	score, err := s.Redis.Get(context.Background(), id).Int()
	if err != nil {
		return 0
	}
	return score
}
