package app

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/BeauRussell/TournamentAssistant/components"
	"github.com/BeauRussell/TournamentAssistant/startgg"
)

func TournamentAssistant(tournamentUrl string, key string) {
	var tournamentSlug string
	parts := strings.Split(tournamentUrl, "/")
	if len(parts) > 4 {
		tournamentSlug = parts[4]
	} else {
		log.Println("Invalid URL format")
		panic("Invalid URL format")
	}
	if key == "" {
		log.Println("No API Key Supplied")
		panic("No API Key Supplied")
	}

	startClient := startgg.Start{}
	startClient.Setup(tournamentSlug, key)
	tournamentData := startClient.GetEventData()

	go func() {
		captureWindow := new(app.Window)
		captureWindow.Option(app.Title("Tournament Assistant"))
		captureWindow.Option(app.Size(unit.Dp(800), unit.Dp(800)))
		err := runApp(captureWindow, tournamentData, &startClient)
		if err != nil {
			log.Println("Window closed with error:", err)
		}

		os.Exit(0)
	}()
	app.Main()
}

func runApp(window *app.Window, tournamentData *startgg.Tournament, startClient *startgg.Start) error {
	theme := material.NewTheme()
	selectBox := components.NewSelectBox(components.ConvertEventsToOptions(tournamentData.Events))

	tickerDuration := 5 * time.Second
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	var progress float32 = 0

	startTime := time.Now()

	go func() {
		for range ticker.C {
			pullEventStandings(startClient, selectBox.SelectedOption)
			pullBracketData(startClient, selectBox.SelectedOption)
			startTime = time.Now()
			progress = 0
		}
	}()

	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)
			paint.Fill(gtx.Ops, color.NRGBA{R: 10, G: 10, B: 10, A: 255})

			elapsed := time.Since(startTime)
			progress = float32(elapsed.Seconds() / tickerDuration.Seconds())

			layout.UniformInset(unit.Dp(25)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return selectBox.Layout(gtx, theme)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						margins := layout.Inset{
							Top:    25,
							Left:   25,
							Bottom: 5,
							Right:  0,
						}

						milliseconds := float32(tickerDuration) / float32(time.Millisecond)
						displayTime := float64((1-progress)*milliseconds) / float64(1000)
						if displayTime < 0 {
							displayTime = 0
						}
						timeLeft := fmt.Sprintf("%.1f sec", displayTime)

						title := material.Label(theme, unit.Sp(25), "Time until refresh: "+timeLeft)
						white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
						title.Color = white
						title.Alignment = text.Middle
						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return title.Layout(gtx)
							},
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.ProgressBar(theme, progress).Layout(gtx)
					}),
				)
			})

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)

			// Force re-rendering for the progress bar
			window.Invalidate()
		}
	}
}
