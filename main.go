package main

import (
    "time"
    "snake/terminal"
    "snake/inputreader"
    "snake/game"
    "snake/renderer"
)

func atExit(t *terminal.Terminal, r *renderer.Renderer, ticker *time.Ticker) {
    t.Restore()
    r.Restore()
    ticker.Stop()
}

func main() {
    t, err := terminal.New()
    if err != nil {
        panic(err)
    }

    ir := inputreader.New()
    g := game.New(t.NRows, t.NCols)
    r := renderer.New()

    tickRate := time.Second / 10
    ticker := time.NewTicker(tickRate)

    events := make(chan byte)
    go ir.Read(events)
    var input byte

    defer atExit(t, r, ticker)
    for {
        select {
        case <-ticker.C:
            r.Render(g)
            g.Update(input)
            if g.GameOver() {
                return
            }
        case input = <-events:
            if input == 'q' {
                return
            }
        }
    }
}
