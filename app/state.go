package app

import (
	"fmt"
	"sync"
	"time"
)

type State struct {
	mu       sync.RWMutex
	videos   []*Video
	selected int
	onChange func()
}

func NewState() *State {
	return &State{
		videos:   make([]*Video, 0),
		selected: -1,
	}
}

func (s *State) SetOnChange(fn func()) {
	s.onChange = fn
}

func (s *State) notifyChange() {
	if s.onChange != nil {
		s.onChange()
	}
}

func (s *State) AddVideo(path string) error {
	video, err := NewVideo(path)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.videos = append(s.videos, video)
	s.mu.Unlock()

	s.notifyChange()
	return nil
}

func (s *State) RemoveSelected() {
	s.mu.Lock()
	if s.selected < 0 || s.selected >= len(s.videos) {
		s.mu.Unlock()
		return
	}

	s.videos = append(s.videos[:s.selected], s.videos[s.selected+1:]...)

	if s.selected >= len(s.videos) {
		s.selected = len(s.videos) - 1
	}
	s.mu.Unlock()

	s.notifyChange()
}

func (s *State) MoveUp() {
	s.mu.Lock()
	if s.selected <= 0 || s.selected >= len(s.videos) {
		s.mu.Unlock()
		return
	}

	s.videos[s.selected], s.videos[s.selected-1] = s.videos[s.selected-1], s.videos[s.selected]
	s.selected--
	s.mu.Unlock()

	s.notifyChange()
}

func (s *State) MoveDown() {
	s.mu.Lock()
	if s.selected < 0 || s.selected >= len(s.videos)-1 {
		s.mu.Unlock()
		return
	}

	s.videos[s.selected], s.videos[s.selected+1] = s.videos[s.selected+1], s.videos[s.selected]
	s.selected++
	s.mu.Unlock()

	s.notifyChange()
}

func (s *State) Clear() {
	s.mu.Lock()
	s.videos = make([]*Video, 0)
	s.selected = -1
	s.mu.Unlock()

	s.notifyChange()
}

func (s *State) SetSelected(index int) {
	s.mu.Lock()
	s.selected = index
	s.mu.Unlock()

	s.notifyChange()
}

func (s *State) GetSelected() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.selected
}

func (s *State) GetVideos() []*Video {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Video, len(s.videos))
	copy(result, s.videos)
	return result
}

func (s *State) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.videos)
}

func (s *State) Move(from, to int) {
	s.mu.Lock()
	if from < 0 || from >= len(s.videos) || to < 0 || to >= len(s.videos) || from == to {
		s.mu.Unlock()
		return
	}

	video := s.videos[from]
	s.videos = append(s.videos[:from], s.videos[from+1:]...)

	newVideos := make([]*Video, 0, len(s.videos)+1)
	newVideos = append(newVideos, s.videos[:to]...)
	newVideos = append(newVideos, video)
	newVideos = append(newVideos, s.videos[to:]...)
	s.videos = newVideos

	s.selected = to
	s.mu.Unlock()

	s.notifyChange()
}

func (s *State) TotalDuration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var total time.Duration
	for _, video := range s.videos {
		total += video.Duration
	}
	return total
}

func (s *State) TotalDurationString() string {
	total := s.TotalDuration()
	if total == 0 {
		return ""
	}

	totalSec := int(total.Seconds())
	hours := totalSec / 3600
	minutes := (totalSec % 3600) / 60
	seconds := totalSec % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
