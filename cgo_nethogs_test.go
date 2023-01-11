package main

import (
    "testing"
    "time"
)

func TestCallNethogs(t *testing.T) {
    CallNethogs("udp portrange 4800-4900")
    time.Sleep(time.Minute)
}
