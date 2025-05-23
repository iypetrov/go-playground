package main

import (
	"fmt"
	"time"
)

const (
	minDuration = time.Duration(1)
)

type Span struct {
	SpanID string
	Start  time.Time
	End    time.Time
}

func (span *Span) isStartBeforeEnd() bool {
	return span.Start.Before(span.End)
}

func (span *Span) compareStartAndEnd() int {
	return span.Start.Compare(span.End)
}

func (span *Span) sanitizeDuration() bool {
	if span.Start.Compare(span.End) == -1 {
		return false
	}

	span.End = span.Start.Add(minDuration)
	return true
}

func printSpanInfo(label string, span *Span) {
	fmt.Printf("before %s.Start = %v\n", label, span.Start)
	fmt.Printf("before %s.End = %v\n", label, span.End)
	fmt.Printf("duration = %v\n", span.End.Sub(span.Start))
	fmt.Printf("isStartBeforeEnd(%s) = %v\n", label, span.isStartBeforeEnd())
	fmt.Printf("compareStartAndEnd(%s) = %d\n", label, span.compareStartAndEnd())
	fmt.Printf("sanitizeDuration(%s) = %v\n", label, span.sanitizeDuration())
	fmt.Printf("after %s.Start = %v\n", label, span.Start)
	fmt.Printf("after %s.End = %v\n", label, span.End)
}

func main() {
	span1 := Span{
		SpanID: "1",
		Start:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	span2 := Span{
		SpanID: "2",
		Start:  time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
		End:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	span3 := Span{
		SpanID: "3",
		Start:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	printSpanInfo("span1", &span1)
	fmt.Printf("---------------------------\n")
	printSpanInfo("span2", &span2)
	fmt.Printf("---------------------------\n")
	printSpanInfo("span3", &span3)
}
