package components

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/BeauRussell/TournamentAssistant/startgg"
)

type Option struct {
	ID   int
	Name string
}

type SelectBox struct {
	SelectedOption Option
	Options        []Option
	IsOpen         bool
	Button         widget.Clickable
	List           layout.List
	OptionClicks   []widget.Clickable
}

func NewSelectBox(options []Option) *SelectBox {
	sb := &SelectBox{
		Options:        options,
		SelectedOption: options[0],
		List:           layout.List{Axis: layout.Vertical},
	}

	sb.OptionClicks = make([]widget.Clickable, len(options))
	return sb
}

func (sb *SelectBox) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if sb.Button.Clicked(gtx) {
		sb.IsOpen = !sb.IsOpen
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &sb.Button, sb.SelectedOption.Name)
			btn.Background = color.NRGBA{R: 100, G: 149, B: 237, A: 255}
			btn.TextSize = unit.Sp(16)
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !sb.IsOpen {
				return layout.Dimensions{}
			}
			return sb.List.Layout(gtx, len(sb.Options), func(gtx layout.Context, i int) layout.Dimensions {
				if sb.OptionClicks[i].Clicked(gtx) {
					sb.SelectedOption = sb.Options[i]
					sb.IsOpen = false
				}
				btn := material.Button(th, &sb.OptionClicks[i], sb.Options[i].Name)
				btn.Background = color.NRGBA{R: 220, G: 220, B: 220, A: 255}
				btn.TextSize = unit.Sp(14)
				btn.Inset = layout.UniformInset(unit.Dp(5))
				return btn.Layout(gtx)
			})
		}),
	)
}

func ConvertEventsToOptions(events []startgg.Event) []Option {
	options := make([]Option, len(events))
	for i, event := range events {
		options[i] = Option{
			ID:   event.ID,
			Name: event.Name,
		}
	}
	return options
}
