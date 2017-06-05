package decor

import (
	"fmt"
	"math"
	"time"
	"unicode/utf8"
)

const (
	// DidentRight specifies identation direction.
	// |   foo|     b| Without DidentRight
	// |foo   |b     | With DidentRight
	DidentRight = 1 << iota

	// DwidthSync will auto sync max width
	DwidthSync

	// DextraSpace adds extra space, makes sence with DwidthSync only.
	// When DidentRight bit set, the space will be added to the right,
	// otherwise to the left.
	DextraSpace
)

// Statistics represents statistics of the progress bar.
// Cantains: Total, Current, TimeElapsed, TimePerItemEstimate
type Statistics struct {
	ID                  int
	Completed           bool
	Aborted             bool
	Total               int
	Current             int
	StartTime           time.Time
	TimeElapsed         time.Duration
	TimePerItemEstimate time.Duration
}

// Eta returns exponential-weighted-moving-average ETA estimator
func (s *Statistics) Eta() time.Duration {
	return time.Duration(s.Total-s.Current) * s.TimePerItemEstimate
}

// DecoratorFunc is a function that can be prepended and appended to the progress bar
type DecoratorFunc func(s *Statistics, myWidth chan<- int, maxWidth <-chan int) string

func Name(name string, minWidth int, conf byte) DecoratorFunc {
	format := "%%"
	if (conf & DidentRight) != 0 {
		format += "-"
	}
	format += "%ds"
	return func(s *Statistics, myWidth chan<- int, maxWidth <-chan int) string {
		if (conf & DwidthSync) != 0 {
			myWidth <- utf8.RuneCountInString(name)
			max := <-maxWidth
			if (conf & DextraSpace) != 0 {
				max++
			}
			return fmt.Sprintf(fmt.Sprintf(format, max), name)
		}
		return fmt.Sprintf(fmt.Sprintf(format, minWidth), name)
	}
}

func Counters(pairFormat string, unit Units, minWidth int, conf byte) DecoratorFunc {
	format := "%%"
	if (conf & DidentRight) != 0 {
		format += "-"
	}
	format += "%ds"
	return func(s *Statistics, myWidth chan<- int, maxWidth <-chan int) string {
		current := Format(s.Current).To(unit)
		total := Format(s.Total).To(unit)
		str := fmt.Sprintf(pairFormat, current, total)
		if (conf & DwidthSync) != 0 {
			myWidth <- utf8.RuneCountInString(str)
			max := <-maxWidth
			if (conf & DextraSpace) != 0 {
				max++
			}
			return fmt.Sprintf(fmt.Sprintf(format, max), str)
		}
		return fmt.Sprintf(fmt.Sprintf(format, minWidth), str)
	}
}

func ETA(minWidth int, conf byte) DecoratorFunc {
	format := "%%"
	if (conf & DidentRight) != 0 {
		format += "-"
	}
	format += "%ds"
	return func(s *Statistics, myWidth chan<- int, maxWidth <-chan int) string {
		str := fmt.Sprint(time.Duration(s.Eta().Seconds()) * time.Second)
		if (conf & DwidthSync) != 0 {
			myWidth <- utf8.RuneCountInString(str)
			max := <-maxWidth
			if (conf & DextraSpace) != 0 {
				max++
			}
			return fmt.Sprintf(fmt.Sprintf(format, max), str)
		}
		return fmt.Sprintf(fmt.Sprintf(format, minWidth), str)
	}
}

func Elapsed(minWidth int, conf byte) DecoratorFunc {
	format := "%%"
	if (conf & DidentRight) != 0 {
		format += "-"
	}
	format += "%ds"
	return func(s *Statistics, myWidth chan<- int, maxWidth <-chan int) string {
		str := fmt.Sprint(time.Duration(s.TimeElapsed.Seconds()) * time.Second)
		if (conf & DwidthSync) != 0 {
			myWidth <- utf8.RuneCountInString(str)
			max := <-maxWidth
			if (conf & DextraSpace) != 0 {
				max++
			}
			return fmt.Sprintf(fmt.Sprintf(format, max), str)
		}
		return fmt.Sprintf(fmt.Sprintf(format, minWidth), str)
	}
}

func Percentage(minWidth int, conf byte) DecoratorFunc {
	format := "%%"
	if (conf & DidentRight) != 0 {
		format += "-"
	}
	format += "%ds"
	return func(s *Statistics, myWidth chan<- int, maxWidth <-chan int) string {
		str := fmt.Sprintf("%d %%", CalcPercentage(s.Total, s.Current, 100))
		if (conf & DwidthSync) != 0 {
			myWidth <- utf8.RuneCountInString(str)
			max := <-maxWidth
			if (conf & DextraSpace) != 0 {
				max++
			}
			return fmt.Sprintf(fmt.Sprintf(format, max), str)
		}
		return fmt.Sprintf(fmt.Sprintf(format, minWidth), str)
	}
}

func CalcPercentage(total, current, width int) int {
	if total == 0 || current > total {
		return 0
	}
	num := float64(width) * float64(current) / float64(total)
	ceil := math.Ceil(num)
	diff := ceil - num
	// num = 2.34 will return 2
	// num = 2.44 will return 3
	if math.Max(diff, 0.6) == diff {
		return int(num)
	}
	return int(ceil)
}
