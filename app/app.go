package app

import (
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
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
		captureWindow.Option(app.Size(unit.Dp(800), unit.Dp(200)))
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

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			pullEventStandings(startClient, selectBox.SelectedOption)
			pullBracketData(startClient, selectBox.SelectedOption)
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

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return selectBox.Layout(gtx, theme)
				}),
			)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
