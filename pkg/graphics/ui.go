package Graphics

import (
	"AI30_-_BlackFriday/pkg/constants"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

func createUI(game *Game) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	slider := widget.NewSlider(
		widget.SliderOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(1, 200),
		),
		widget.SliderOpts.Orientation(widget.DirectionVertical),
		widget.SliderOpts.MinMax(1, constants.MAX_TIC_DURATION),
		widget.SliderOpts.InitialCurrent(constants.TIC_DURATION),
		widget.SliderOpts.Images(
			createSliderTrackImage(),
			createSliderButtonImage(),
		),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			game.Simulation.SetTicDuration(args.Current)
		}),
	)

	rootContainer.AddChild(slider)

	return &ebitenui.UI{
		Container: rootContainer,
	}
}

func createSliderTrackImage() *widget.SliderTrackImage {
	trackIdle := image.NewNineSliceColor(color.RGBA{100, 100, 100, 255})
	trackHover := image.NewNineSliceColor(color.RGBA{120, 120, 120, 255})

	return &widget.SliderTrackImage{
		Idle:     trackIdle,
		Hover:    trackHover,
		Disabled: trackIdle,
	}
}

func createSliderButtonImage() *widget.ButtonImage {
	buttonIdle := image.NewNineSliceColor(color.RGBA{200, 200, 200, 255})
	buttonHover := image.NewNineSliceColor(color.RGBA{220, 220, 220, 255})
	buttonPressed := image.NewNineSliceColor(color.RGBA{180, 180, 180, 255})

	return &widget.ButtonImage{
		Idle:     buttonIdle,
		Hover:    buttonHover,
		Pressed:  buttonPressed,
		Disabled: buttonIdle,
	}
}
