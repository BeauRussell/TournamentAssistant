package main

import (
	"database/sql"
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

	_ "github.com/mattn/go-sqlite3"

	tourneyApp "github.com/BeauRussell/TournamentAssistant/app"
)

type captureResult struct {
	tournamentSlug string
	key            string
}

func main() {
	resultChan := make(chan captureResult)
	go func() {
		captureWindow := new(app.Window)
		captureWindow.Option(app.Title("Tournament Assistant"))
		captureWindow.Option(app.Size(unit.Dp(800), unit.Dp(300)))
		var result captureResult
		result.key = checkKey()
		err := runCapture(captureWindow, &result.tournamentSlug, &result.key)
		if err != nil {
			log.Println("Window closed with error:", err)
		}
		resultChan <- result
		close(resultChan)
	}()

	result := <-resultChan
	writeKey(result.key)

	tourneyApp.TournamentAssistant(result.tournamentSlug, result.key)
	app.Main()
}

func runCapture(window *app.Window, tournamentSlug *string, key *string) error {
	theme := material.NewTheme()

	var submitButton widget.Clickable
	var slugInput widget.Editor
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
				*tournamentSlug = strings.TrimSpace(slugInput.Text())
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

func checkKey() string {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Println("Failed to Open DB file")
		panic(err)
	}
	defer db.Close()

	createCmd := `
		create table if not exists smash (id integer not null primary key, api);
	`
	_, err = db.Exec(createCmd)
	if err != nil {
		log.Printf("%q: %s\n", err, createCmd)
		panic(err)
	}

	checkCmd := `
		select api from smash where id = 1;
	`

	var key string
	err = db.QueryRow(checkCmd).Scan(&key)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		log.Printf("%q: %s\n", err, checkCmd)
		panic(err)
	}
	log.Println(key)
	return key
}

func writeKey(key string) {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Println("Failed to Open DB file")
		panic(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to begin insert transaction")
		panic(err)
	}
	stmt, err := tx.Prepare("insert or replace into smash (id, api) values(?,?)")
	if err != nil {
		log.Println("Failed to prepare insert transaction")
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(1, key)
	if err != nil {
		log.Println("Failed to insert key")
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Failed to commit insert transaction")
		panic(err)
	}
}
