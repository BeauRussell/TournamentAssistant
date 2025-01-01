package app

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/BeauRussell/TournamentAssistant/startgg"
)

func TournamentAssistant(tournamentUrl string, key string) {
	var tournamentSlug string
	parts := strings.Split(tournamentUrl, "/")
	if len(parts) > 4 {
		tournamentSlug = parts[4]
	} else {
		fmt.Println("Invalid URL format")
	}

	startClient := startgg.Start{}
	startClient.Setup(tournamentSlug, key)

	go func() {
		captureWindow := new(app.Window)
		captureWindow.Option(app.Title("Tournament Assistant"))
		captureWindow.Option(app.Size(unit.Dp(800), unit.Dp(200)))
		err := runApp(captureWindow, tournamentSlug)
		if err != nil {
			log.Println("Window closed with error:", err)
		}

		os.Exit(0)
	}()
	app.Main()
}

func runApp(window *app.Window, tournamentSlug string) error {
	theme := material.NewTheme()
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
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						margins := layout.Inset{
							Top:    25,
							Left:   25,
							Bottom: 5,
							Right:  0,
						}

						title := material.Label(theme, unit.Sp(25), tournamentSlug)
						white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
						title.Color = white
						title.Alignment = text.Middle
						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return title.Layout(gtx)
							},
						)
					},
				),
			)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
