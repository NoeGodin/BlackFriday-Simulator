package Hud

import (
	"AI30_-_BlackFriday/pkg/constants"
	"bytes"
	"fmt"
	i "image"
	"image/color"
	"log"
	"os"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/utilities/constantutil"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})
	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})
	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadFont(size float64) (text.Face, error) {
	fontData, err := os.ReadFile(constants.FONT_PATH)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

func createComboBox(entries []any, f widget.ListComboButtonEntrySelectedHandlerFunc) *widget.Container {
	buttonImage, _ := loadButtonImage()
	face, _ := loadFont(20)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	comboBox := widget.NewListComboButton(
		widget.ListComboButtonOpts.Entries(entries),
		widget.ListComboButtonOpts.MaxContentHeight(150),
		widget.ListComboButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ListComboButtonOpts.ButtonParams(&widget.ButtonParams{
			Image:       buttonImage,
			TextPadding: widget.NewInsetsSimple(5),
			TextColor: &widget.ButtonTextColor{
				Idle:     color.White,
				Disabled: color.White,
			},
			TextFace: &face,
			MinSize:  &i.Point{200, 0},
		}),
		widget.ListComboButtonOpts.ListParams(&widget.ListParams{
			ScrollContainerImage: &widget.ScrollContainerImage{
				Idle:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Disabled: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Mask:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			Slider: &widget.SliderParams{
				TrackImage: &widget.SliderTrackImage{
					Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
					Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				},
				HandleImage:   buttonImage,
				MinHandleSize: constantutil.ConstantToPointer(5),
				TrackPadding:  widget.NewInsetsSimple(2),
			},
			EntryFace: &face,
			EntryColor: &widget.ListEntryColor{
				Selected:                   color.NRGBA{254, 255, 255, 255},             //Foreground color for the unfocused selected entry
				Unselected:                 color.NRGBA{254, 255, 255, 255},             //Foreground color for the unfocused unselected entry
				SelectedBackground:         color.NRGBA{R: 130, G: 130, B: 200, A: 255}, //Background color for the unfocused selected entry
				SelectedFocusedBackground:  color.NRGBA{R: 130, G: 130, B: 170, A: 255}, //Background color for the focused selected entry
				FocusedBackground:          color.NRGBA{R: 170, G: 170, B: 180, A: 255}, //Background color for the focused unselected entry
				DisabledUnselected:         color.NRGBA{100, 100, 100, 255},             //Foreground color for the disabled unselected entry
				DisabledSelected:           color.NRGBA{100, 100, 100, 255},             //Foreground color for the disabled selected entry
				DisabledSelectedBackground: color.NRGBA{100, 100, 100, 255},             //Background color for the disabled selected entry
			},
			EntryTextPadding: widget.NewInsetsSimple(5),
			MinSize:          &i.Point{200, 0},
		}),

		widget.ListComboButtonOpts.EntryLabelFunc(
			func(e any) string {
				return e.(ListEntry).name
			},
			func(e any) string {
				return e.(ListEntry).name
			}),
		widget.ListComboButtonOpts.EntrySelectedHandler(f),
	)

	comboBox.SetSelectedEntry(entries[0])
	rootContainer.AddChild(comboBox)
	return rootContainer
}

func createSliderWithLabel(label string, min, max, defaultValue int, suffix string, changedHandler widget.SliderChangedHandlerFunc) *widget.Container {
	fontFace, _ := loadFont(20)
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// Container for label and value
	labelContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// Label
	labelWidget := widget.NewText(
		widget.TextOpts.Text(label, &fontFace, color.White),
	)
	labelContainer.AddChild(labelWidget)

	// Value display
	initialValue := fmt.Sprintf("%d%s", defaultValue, suffix)
	valueWidget := widget.NewText(
		widget.TextOpts.Text(initialValue, &fontFace, color.RGBA{100, 200, 255, 255}),
	)
	labelContainer.AddChild(valueWidget)

	container.AddChild(labelContainer)

	// Slider
	orientation := widget.DirectionHorizontal
	slider := widget.NewSlider(
		widget.SliderOpts.Orientation(orientation),
		widget.SliderOpts.MinMax(min, max),
		widget.SliderOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.MinSize(400, 20),
		),
		widget.SliderOpts.Images(
			&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.RGBA{50, 50, 100, 255}),
				Hover: image.NewNineSliceColor(color.RGBA{60, 60, 120, 255}),
			},
			&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.RGBA{100, 100, 200, 255}),
				Hover:   image.NewNineSliceColor(color.RGBA{120, 120, 220, 255}),
				Pressed: image.NewNineSliceColor(color.RGBA{80, 80, 180, 255}),
			},
		),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			valueWidget.Label = fmt.Sprintf("%d%s", args.Current, suffix)
			changedHandler(args)
		}),
	)
	slider.Current = defaultValue
	container.AddChild(slider)

	return container
}
