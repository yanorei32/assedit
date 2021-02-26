package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/yanorei32/assedit/assutil"
	"github.com/yanorei32/assedit/avsparser"
)

func frame2duration(frame uint64, framerate float64) time.Duration {
	frameDurationNs := float64(time.Second.Nanoseconds()) / framerate

	totalDurationNs := frameDurationNs * float64(frame)

	totalDurationNsStr := fmt.Sprintf("%dns", uint64(totalDurationNs))

	d, _ := time.ParseDuration(totalDurationNsStr)

	return d
}

func main() {
	assP := flag.String("x", "", "[required] Input SubStation Alpha file (.ass)")
	avsP := flag.String("t", "", "[required] AviSynthScript (.avs)")
	outputAssP := flag.String("o", "", "[optional] Output file (default: .trimmed.ass)")
	framerate := flag.Float64("f", 29.97, "framerate (example: 29.97, 59.94, 119.88)")

	flag.Parse()

	if *avsP == "" || *assP == "" {
		flag.PrintDefaults()
		return
	}

	if !strings.HasSuffix(*avsP, ".avs") {
		log.Fatal("AVS file name is not '*.avs'")
	}

	avsF, err := os.Open(*avsP)
	if err != nil {
		log.Fatal("AVS (open): " + err.Error())
	}

	avs, err := avsparser.Parse(bufio.NewReader(avsF))
	if err != nil {
		log.Fatal("AVS (parse): " + err.Error())
	}

	if !strings.HasSuffix(*assP, ".ass") {
		log.Fatal("ASS file name is not '*.ass'")
	}

	assF, err := os.Open(*assP)
	if err != nil {
		log.Fatal("ASS (open): " + err.Error())
	}

	ass, err := assutil.Parse(bufio.NewReader(assF))
	if err != nil {
		log.Fatal("ASS (parse): " + err.Error())
	}

	outputAssP_ := *outputAssP

	if outputAssP_ == "" {
		outputAssP_ = (*assP)[:len(*assP)-3] + "trimmed.ass"
	}

	outputAssF, err := os.Create(outputAssP_)

	if err != nil {
		log.Fatal("Output ASS (create): " + err.Error())
	}

	outputAss := assutil.Ass{
		OtherSectionsHead: ass.OtherSectionsHead,
		OtherSectionsTail: ass.OtherSectionsTail,
	}

	lastFrame := uint64(0)
	trimmedFrames := uint64(0)

	for _, trim := range avs.Trims {
		trimmedFrames += trim.StartFrame - lastFrame
		lastFrame = trim.EndFrame

		trimmedDuration := frame2duration(trimmedFrames, *framerate)
		sectionHead := frame2duration(trim.StartFrame, *framerate)
		sectionEnd := frame2duration(trim.EndFrame, *framerate)

		for _, event := range ass.EventSection {
			if event.EndTime < sectionHead {
				continue
			}

			if sectionEnd < event.StartTime {
				continue
			}

			if event.StartTime < sectionHead {
				event.StartTime = sectionHead
			}

			if sectionEnd < event.EndTime {
				event.EndTime = sectionEnd
			}

			event.StartTime -= trimmedDuration
			event.EndTime -= trimmedDuration

			outputAss.EventSection = append(outputAss.EventSection, event)
		}
	}

	outputAssF.WriteString(outputAss.Format())
}
