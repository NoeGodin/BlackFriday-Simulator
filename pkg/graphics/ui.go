package Graphics

import (
	"AI30_-_BlackFriday/pkg/constants"
	Hud "AI30_-_BlackFriday/pkg/hud"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

func createUI(game *Game) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10), // espace entre boutons
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				Padding:            &widget.Insets{Top: 10, Left: 10},
			}),
		),
	)
	fontFace, _ := Hud.LoadFont(20)
	pauseButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
			}),
			widget.WidgetOpts.MinSize(80, 30),
		),
		widget.ButtonOpts.Image(createButtonImage()),
		widget.ButtonOpts.Text("Pause", &fontFace, &widget.ButtonTextColor{
			Idle: color.RGBA{255, 255, 255, 255},
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.TogglePause()
		}),
	)

	guardButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
			}),
			widget.WidgetOpts.MinSize(80, 30),
		),
		widget.ButtonOpts.Image(createButtonImage()),
		widget.ButtonOpts.Text("afficher FOV gardes", &fontFace, &widget.ButtonTextColor{
			Idle: color.RGBA{255, 255, 255, 255},
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.guardsFOV = !game.guardsFOV
		}),
	)
	menuButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
			}),
			widget.WidgetOpts.MinSize(80, 30),
		),
		widget.ButtonOpts.Image(createButtonImage()),
		widget.ButtonOpts.Text("Menu", &fontFace, &widget.ButtonTextColor{
			Idle: color.RGBA{255, 255, 255, 255},
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.inMenu = true
			go game.Simulation.Stop()
			game.Simulation = nil
		}),
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
		widget.SliderOpts.MinMax(constants.MIN_TIC_DURATION, constants.MAX_TIC_DURATION),
		widget.SliderOpts.InitialCurrent(constants.TIC_DURATION),
		widget.SliderOpts.Images(
			createSliderTrackImage(),
			createSliderButtonImage(),
		),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			game.SetTicDuration(args.Current)
		}),
	)

	buttonsContainer.AddChild(guardButton)
	buttonsContainer.AddChild(pauseButton)
	buttonsContainer.AddChild(menuButton)
	rootContainer.AddChild(buttonsContainer)
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

func createButtonImage() *widget.ButtonImage {
	buttonIdle := image.NewNineSliceColor(color.RGBA{70, 130, 180, 255})
	buttonHover := image.NewNineSliceColor(color.RGBA{100, 149, 237, 255})
	buttonPressed := image.NewNineSliceColor(color.RGBA{50, 100, 150, 255})

	return &widget.ButtonImage{
		Idle:     buttonIdle,
		Hover:    buttonHover,
		Pressed:  buttonPressed,
		Disabled: buttonIdle,
	}
}
