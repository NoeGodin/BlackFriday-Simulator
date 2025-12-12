package mapgenerator

type MapLayout struct {
	width, height 	int
	mapLayout 		[][]string
}

func NewMapLayout(width, height int) (MapLayout) {
	mapLayout := make([][]string, height)
	for y := range height {
		mapLayout[y] = make([]string, width)
	}

	return MapLayout{
		mapLayout: mapLayout,
		width: width,
		height: height,
	}
}

func (m *MapLayout) ToString() string {
	mapStr := ""

	for y := range m.height {
		for x := range m.width {
			if(m.mapLayout[y][x] == "") {
				mapStr += " "
			} 
			mapStr += m.mapLayout[y][x]
		}
		mapStr += "\n"
	}

	return mapStr
}

func (m *MapLayout) CleanMapLayout() {
	for y := range m.height {
		for x := range m.width {
			m.mapLayout[y][x] = ""
		}
	}
}