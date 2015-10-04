package colorcandy

// TODO(Alexander Yunin): make color type for SRGB
// add func to convert 16bit colors to 8bit
// use new types for something else (for map with pix-count-percent ?)

import (
	"math"
	"math/rand"
	"sort"

	"github.com/improvemedia/colorcandy.git/candy"
)

type Config struct {
	ClusterColorsStr    map[string]string `json:"cluster_colors"`
	ClusterColors       map[Color]Color   `json:"-"`
	ColorsCount         uint              `json:"colors_count"`
	PaletteColorsMaxNum int               `json:"palette_colors_max_num"`
	WhiteThreshold      int               `json:"white_threshold"`
	BlackThreshold      int               `json:"black_threshold"`
	Delta               float64           `json:"delta"`
}

type ColorCandy struct {
	Config
}

func New(config Config) *ColorCandy {
	// defaults
	if config.ColorsCount == 0 {
		config.ColorsCount = 60
	}
	if config.PaletteColorsMaxNum == 0 {
		config.PaletteColorsMaxNum = 5
	}
	if config.WhiteThreshold == 0 {
		config.WhiteThreshold = 55000
	}
	if config.BlackThreshold == 0 {
		config.BlackThreshold = 2000
	}
	if config.Delta == 0 {
		config.Delta = 2.5
	}

	config.ClusterColors = map[Color]Color{}
	for k, v := range config.ClusterColorsStr {
		config.ClusterColors[ColorFromString(k)] = ColorFromString(v)
	}

	return &ColorCandy{config}
}

func (colorCandy *ColorCandy) Candify(path string, searchColors []string) (*candy.Result, error) {
	colors, colorsCount, baseColorsCount := colorCandy.ExtractColors(path)

	var paletteColors map[string]*ColorCount
	if len(searchColors) == 0 {
		paletteColors = colorsCount
	} else {
		for _, searchColor := range searchColors {
			if baseColorCount, ok := baseColorsCount[searchColor]; ok {
				for color, count := range baseColorCount {
					paletteColors[color] = count
				}
			}
		}
	}

	palette := map[string]*candy.ColorCount{}
	for k, v := range colorCandy.CreatePalette(paletteColors) {
		palette[k] = &candy.ColorCount{
			Total:      v.Total,
			Percentage: v.Percentage,
		}
	}

	return &candy.Result{
		Colors:  colors,
		Palette: palette,
	}, nil
}

func (colorCandy *ColorCandy) ExtractColors(path string) (map[string]*candy.ColorMeta, map[string]*ColorCount, map[string]map[string]*ColorCount) {
	histogram := ImageHistogram(path)
	compacted := CompactToCommonColors(histogram, colorCandy.Delta)

	return colorCandy.extractColorsFromHistogram(compacted)
}

func (colorCandy *ColorCandy) extractColorsFromHistogram(histogram map[Color]*ColorCount) (map[string]*candy.ColorMeta, map[string]*ColorCount, map[string]map[string]*ColorCount) {
	colors := map[string]*candy.ColorMeta{}
	colorsCount := map[string]*ColorCount{}
	baseColorsCount := map[string]map[string]*ColorCount{}

	for color, count := range histogram {
		baseColor, delta := colorCandy.closestBaseColorTo(color)
		baseColorHex := baseColor.Hex()

		colorsCount[color.Hex()] = count
		if _, ok := baseColorsCount[baseColorHex]; ok {
			baseColorsCount[baseColorHex][color.Hex()] = count
		} else {
			baseColorsCount[baseColorHex] = map[string]*ColorCount{
				color.Hex(): count,
			}
		}

		if meta, found := colors[baseColor.Hex()]; found {
			meta.SearchFactor += count.Percentage
			meta.Colors = append(meta.Colors, color.Hex())
		} else {
			colors[baseColorHex] = &candy.ColorMeta{
				Colors:       []string{color.Hex()},
				BaseColor:    baseColor.Hex(),
				SearchFactor: count.Percentage,
				Distance:     delta,
			}
		}
	}

	return colors, colorsCount, baseColorsCount
}

func (colorCandy *ColorCandy) CreatePalette(colors map[string]*ColorCount) (result map[string]*ColorCount) {
	for len(colors) > colorCandy.PaletteColorsMaxNum {
		colorsArr := []*ColorCount{}
		for _, v := range colors {
			colorsArr = append(colorsArr, v)
		}

		matrix := make([][]float64, len(colors))
		for i, row := range colorsArr {
			matrix[i] = make([]float64, len(colors))
			for j, col := range colorsArr {
				rgbColor1 := row.color
				rgbColor2 := col.color
				pixel1 := Lab.Convert(rgbColor1)
				pixel2 := Lab.Convert(rgbColor2)
				diff := DeltaE(pixel1, pixel2)
				if diff == 0 {
					diff = 100000
				}
				matrix[i][j] = diff
			}
		}

		pos1, pos2 := 0, 0
		var min float64 = math.MaxInt64
		for i := 0; i < len(colorsArr); i++ {
			for j := 0; j < len(colorsArr); j++ {
				v := matrix[i][j]
				if v < min {
					min = v
					pos1, pos2 = i, j
				}
			}
		}

		add, remove := LabMerge(colorsArr[pos1], colorsArr[pos2])
		colors[add.color.Hex()] = add
		delete(colors, remove.color.Hex())
	}
	if len(colors) == colorCandy.PaletteColorsMaxNum {
		return
	}
	for len(colors) < colorCandy.PaletteColorsMaxNum {
		colorsArr := []Color{}
		for k, _ := range colors {
			colorsArr = append(colorsArr, ColorFromString(k))
		}
		r, g, b := colorsArr[rand.Intn(len(colorsArr))].RGB()
		min := r
		max := r
		if g < min {
			min = g
		}
		if g > max {
			max = g
		}
		if b < min {
			min = b
		}
		if b > max {
			max = b
		}
		shifts := [][]int{}
		for i := -min; i < 255-max; i++ {
			if math.Abs(float64(i)) < 30 {
				shifts = append(shifts, []int{int(i), 2, -int(math.Abs(float64(i)))})
			} else if int(math.Abs(float64(i)))%30 != 0 {
				shifts = append(shifts, []int{int(i), 1, int(math.Abs(float64(i)))})
			} else {
				shifts = append(shifts, []int{int(i), 0, int(math.Abs(float64(i)))})
			}
		}

		sort.Sort(ShiftSorter{shifts})

		d := colorCandy.PaletteColorsMaxNum - len(colors)
		for _, e := range shifts[0:d] {
			shift := e[0]
			newColor := Color{r + int32(shift), g + int32(shift), b + int32(shift)}
			colors[newColor.Hex()] = &ColorCount{newColor, 1, 2}
		}
	}

	return
}

func (colorCandy *ColorCandy) closestBaseColorTo(c Color) (Color, float64) {
	var closestColor Color
	var cluster Color
	minDelta := math.MaxFloat64 // pls no buf overflow

	lab := Lab.Convert(c)
	for k, _ := range colorCandy.ClusterColors {
		delta := DeltaE(Lab.Convert(k), lab)
		if delta < minDelta {
			minDelta = delta
			closestColor = k
			cluster = colorCandy.ClusterColors[k]
		}
	}
	return cluster, DeltaE(Lab.Convert(cluster), Lab.Convert(closestColor))
}

type ShiftSorter struct {
	Shifts [][]int
}

func (s ShiftSorter) Len() int      { return len(s.Shifts) }
func (s ShiftSorter) Swap(i, j int) { s.Shifts[i], s.Shifts[j] = s.Shifts[j], s.Shifts[i] }
func (s ShiftSorter) Less(i, j int) bool {
	if s.Shifts[i][1] == s.Shifts[j][1] {
		return s.Shifts[i][2] < s.Shifts[j][2]
	}
	return s.Shifts[i][1] < s.Shifts[j][1]
}
