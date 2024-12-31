package main

import (
	"fmt"
	"image/color"
	"log"
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
)

func main() {
	tournamentSlugChan := make(chan string)
	go func() {
		captureWindow := new(app.Window)
		captureWindow.Option(app.Title("Tournament Assistant"))
		captureWindow.Option(app.Size(unit.Dp(800), unit.Dp(200)))
		var tournamentSlug string
		err := runCapture(captureWindow, &tournamentSlug)
		if err != nil {
			log.Println("Window closed with error:", err)
		}
		tournamentSlugChan <- tournamentSlug // Send the slug back to the main thread
		close(tournamentSlugChan)            // Close the channel
		app.Main()                           // Run the Gio event loop
	}()

	// Wait for the result from the goroutine
	tournamentSlug := <-tournamentSlugChan
	fmt.Println("Tournament Slug:", tournamentSlug)
}

func runCapture(window *app.Window, tournamentSlug *string) error {
	theme := material.NewTheme()

	var submitButton widget.Clickable
	var slugInput widget.Editor

	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)
			if submitButton.Clicked(gtx) {
				*tournamentSlug = strings.TrimSpace(slugInput.Text())
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
						slugInput.SingleLine = true
						slugInput.Alignment = text.Middle
						inputBox := material.Editor(theme, &slugInput, "")
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
