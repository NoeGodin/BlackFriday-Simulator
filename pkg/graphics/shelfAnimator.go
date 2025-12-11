package Graphics

import (
	Map "AI30_-_BlackFriday/pkg/map"

	"github.com/hajimehoshi/ebiten/v2"
)

type ShelfAnimator struct {
	shelvesQuantity map[[2]float64]int
}

const (
	FULL_TRESHOLD         = 1.0
	ALMOST_FULL_TRESHOLD  = 0.66
	HALF_EMPTY_TRESHOLD   = 0.33
	ALMOST_EMPTY_TRESHOLD = 0.0
)

func (animator *ShelfAnimator) AnimationFrame(pos [2]float64, shelf *Map.Shelf) *ebiten.Image {
	firstQuantity, ok := animator.shelvesQuantity[pos]
	var currentQuantity int
	for _, item := range shelf.Items {
		currentQuantity += item.Quantity
	}
	if !ok {
		animator.shelvesQuantity[pos] = currentQuantity
		return itemImg
	}
	ratio := float64(currentQuantity) / float64(firstQuantity)
	if ratio == FULL_TRESHOLD {
		return itemImg
	}
	if ratio > ALMOST_FULL_TRESHOLD {
		return itemAlmostFullImg
	}
	if ratio > HALF_EMPTY_TRESHOLD {
		return itemHalfEmptyImg
	}
	if ratio > ALMOST_EMPTY_TRESHOLD {
		return itemAlmostEmptyImg
	}
	return itemEmptyImg
}

func NewShelfAnimator() *ShelfAnimator {
	return &ShelfAnimator{shelvesQuantity: make(map[[2]float64]int)}
}
