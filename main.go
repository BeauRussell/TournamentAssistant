package main

import (
	"image/color"
	"log"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	tourneyApp "github.com/BeauRussell/TournamentAssistant/app"
	"github.com/BeauRussell/TournamentAssistant/db"
)

type captureResult struct {
	tournamentUrl string
	key           string
}

func main() {
	logFile, err := os.OpenFile("output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	resultChan := make(chan captureResult)
	go func() {
		captureWindow := new(app.Window)
		captureWindow.Option(app.Title("Tournament Assistant"))
		captureWindow.Option(app.Size(unit.Dp(800), unit.Dp(300)))
		var result captureResult
		result.key, result.tournamentUrl = db.CheckCaptureCache()
		err := runCapture(captureWindow, &result.tournamentUrl, &result.key)
		if err != nil {
			log.Println("Window closed with error:", err)
		}
		resultChan <- result
		close(resultChan)
	}()

	result := <-resultChan
	db.WriteCapture(result.key, result.tournamentUrl)

	tourneyApp.TournamentAssistant(result.tournamentUrl, result.key)
	app.Main()
}

func runCapture(window *app.Window, tournamentUrl *string, key *string) error {
	theme := material.NewTheme()

	var submitButton widget.Clickable
	var urlInput widget.Editor
	urlInput.SetText(*tournamentUrl)
	var keyInput widget.Editor
	keyInput.SetText(*key)

	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)
			if submitButton.Clicked(gtx) {
				*tournamentUrl = strings.TrimSpace(urlInput.Text())
				*key = strings.TrimSpace(keyInput.Text())
				window.Perform(system.ActionClose)
				continue
			}

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

						title := material.Label(theme, unit.Sp(25), "Start.gg URL")
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
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						urlInput.SingleLine = true
						urlInput.Alignment = text.Middle
						inputBox := material.Editor(theme, &urlInput, "")
						inputBox.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
						inputBox.TextSize = unit.Sp(16)

						margins := layout.Inset{
							Top:    unit.Dp(0),
							Right:  unit.Dp(25),
							Bottom: unit.Dp(10),
							Left:   unit.Dp(25),
						}

						border := widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(3),
							Width:        unit.Dp(2),
						}

						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return border.Layout(gtx, inputBox.Layout)
							},
						)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						margins := layout.Inset{
							Top:    25,
							Left:   25,
							Bottom: 5,
							Right:  0,
						}

						title := material.Label(theme, unit.Sp(25), "API Key")
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
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						keyInput.SingleLine = true
						keyInput.Alignment = text.Middle
						inputBox := material.Editor(theme, &keyInput, "")
						inputBox.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
						inputBox.TextSize = unit.Sp(16)

						margins := layout.Inset{
							Top:    unit.Dp(0),
							Right:  unit.Dp(25),
							Bottom: unit.Dp(10),
							Left:   unit.Dp(25),
						}

						border := widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(3),
							Width:        unit.Dp(2),
						}

						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return border.Layout(gtx, inputBox.Layout)
							},
						)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						margins := layout.Inset{
							Top:    unit.Dp(0),
							Bottom: unit.Dp(25),
							Left:   unit.Dp(25),
							Right:  unit.Dp(25),
						}

						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {

								submit := material.Button(theme, &submitButton, "Submit")
								return submit.Layout(gtx)
							},
						)
					},
				),

				layout.Rigid(
					layout.Spacer{Height: unit.Dp(25)}.Layout,
				),
			)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
