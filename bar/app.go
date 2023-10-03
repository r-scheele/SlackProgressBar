package bar

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/nlopes/slack"
)

type SlackProgress struct {
	Token              string
	Channel            string
	Suffix             string
	ProgressBarChar    string
	Precision          int
	RateLimit          time.Duration
	CompletionCallback func()
}

type ProgressBar struct {
	Sp        *SlackProgress
	Total     int
	Done      int
	Pos       int
	MsgTs     string
	ChannelID string
	MsgLog    []string
	StartTime time.Time
	Paused    bool
	ResumeCh  chan bool
	mu        sync.Mutex
}

func (sp *SlackProgress) New(total int) (*ProgressBar, error) {
	api := slack.New(sp.Token)
	channelID, timestamp, err := api.PostMessage(sp.Channel, slack.MsgOptionText("Starting...", false))
	if err != nil {
		return nil, fmt.Errorf("unable to post initial message: %w", err)
	}
	return &ProgressBar{
		Sp:        sp,
		Total:     total,
		MsgTs:     timestamp,
		ChannelID: channelID,
		StartTime: time.Now(),
		ResumeCh:  make(chan bool),
	}, nil // Return nil error
}

func (pb *ProgressBar) Update() {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if pb.Paused {
		<-pb.ResumeCh
	}

	ticker := time.NewTicker(pb.Sp.RateLimit)
	defer ticker.Stop()

	<-ticker.C

	// Update the position and percentage completion
	pb.Pos = int(float64(pb.Done) / float64(pb.Total) * 100)

	// Call makeBar to get the progress bar string with memory usage
	bar := pb.makeBar()

	_, _, _, err := slack.New(pb.Sp.Token).UpdateMessage(pb.ChannelID, pb.MsgTs, slack.MsgOptionText(bar, false))
	if err != nil {
		pb.Log(fmt.Sprintf("Error updating progress bar: %v", err))
	}
}

func (pb *ProgressBar) Log(msg string) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	timestamp := time.Now().Format("15:04:05")
	pb.MsgLog = append(pb.MsgLog, fmt.Sprintf("*%s* - [%s]", timestamp, msg))
	content := fmt.Sprintf("%s\n%s", pb.makeBar(), strings.Join(pb.MsgLog, "\n"))

	_, _, _, err := slack.New(pb.Sp.Token).UpdateMessage(pb.ChannelID, pb.MsgTs, slack.MsgOptionText(content, false))
	if err != nil {
		fmt.Printf("Error logging message: %v\n", err)
	}
}

func (pb *ProgressBar) makeBar() string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	memUsage := fmt.Sprintf("Memory Usage: %v KB", mem.Alloc/1024)

	elapsed := time.Since(pb.StartTime).Seconds()
	remaining := (elapsed / float64(pb.Done+1)) * float64(pb.Total-pb.Done-1)
	eta := time.Duration(remaining) * time.Second

	perc := fmt.Sprintf("%.*f", pb.Sp.Precision, float64(pb.Pos))
	bar := fmt.Sprintf("%s %s%s ETA: %v", strings.Repeat(pb.Sp.ProgressBarChar, pb.Pos), perc, pb.Sp.Suffix, eta)
	return fmt.Sprintf("%s | %s", bar, memUsage)
}

func (pb *ProgressBar) Pause() {
	pb.Paused = true
}

func (pb *ProgressBar) Resume() {
	pb.Paused = false
	pb.ResumeCh <- true
}

func (sp *SlackProgress) Iter(items []string, processFunc func(item string)) {
	pb, err := sp.New(len(items))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for idx, item := range items {
		processFunc(item)
		pb.Done = idx + 1 // Increment idx by 1 to reflect the completed item
		pb.Update()
	}
	if sp.CompletionCallback != nil {
		sp.CompletionCallback()
	}
}
