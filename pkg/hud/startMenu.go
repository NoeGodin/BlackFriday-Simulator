package Hud

import (
	"AI30_-_BlackFriday/pkg/constants"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Properties struct {
	Filename            string
	ClientNumber        int
	GuardNumber         int
	AgentAggressiveness float64
}
type ListEntry struct {
	name string
}
type StartMenu struct {
	ui         *ebitenui.UI
	properties *Properties
}
type OnStartFunc func(Properties)

func NewStartMenu(onStart OnStartFunc) *StartMenu {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{20, 20, 40, 255})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, false, false}),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.GridLayoutOpts.Spacing(20, 20),
		)),
	)
	game := StartMenu{
		ui: &ebitenui.UI{
			Container: rootContainer,
		},
		properties: &Properties{},
	}
	slider1Container := createSliderWithLabel("Nombre de clients", 0, 2000, constants.NUMBER_OF_CLIENTS, "", func(args *widget.SliderChangedEventArgs) {
		game.properties.ClientNumber = args.Current
	})
	rootContainer.AddChild(slider1Container)

	sliderGuardContainer := createSliderWithLabel("Nombre d'agents de sécurité", 0, 100, constants.NUMBER_OF_GUARDS, "", func(args *widget.SliderChangedEventArgs) {
		game.properties.GuardNumber = args.Current
	})
	rootContainer.AddChild(sliderGuardContainer)

	slider2Container := createSliderWithLabel("Taux d'aggressivité moyen", 0, 100, 0, "%", func(args *widget.SliderChangedEventArgs) {
		game.properties.AgentAggressiveness = float64(args.Current) / float64(args.Slider.Max)
	})
	rootContainer.AddChild(slider2Container)

	entries, err := os.ReadDir(constants.MAPS_PATH)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, v := range entries {
		if v.IsDir() {
			continue
		}
		if filepath.Ext(v.Name()) == ".txt" {
			files = append(files, v.Name())
		}
	}
	boxEntries := make([]any, 0, len(files))
	for _, file := range files {
		boxEntries = append(boxEntries, ListEntry{file})
	}
	rootContainer.AddChild(createComboBox(boxEntries, func(args *widget.ListComboButtonEntrySelectedEventArgs) {
		game.properties.Filename = args.Button.Label()
	}))

	// Simple button
	buttonImage, _ := loadButtonImage()
	fontFace, _ := loadFont(20)
	button := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Démarrer", &fontFace, &widget.ButtonTextColor{
			Idle: color.White,
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Left:   30,
			Right:  30,
			Top:    10,
			Bottom: 10,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			onStart(*game.properties)
		}),
	)
	rootContainer.AddChild(button)

	return &game

}

func (g *StartMenu) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *StartMenu) Update() error {
	g.ui.Update()

	if list, ok := g.ui.GetFocusedWidget().(*widget.ListComboButton); ok {
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			list.FocusPrevious()
		} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			list.FocusNext()
		} else if inpututil.IsKeyJustPressed(ebiten.KeyB) {
			list.SelectFocused()
		}
	}
	return nil
}

func (g *StartMenu) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)

}
