package colorcandy

// TODO(Alexander Yunin): make color type for SRGB
// add func to convert 16bit colors to 8bit
// use new types for something else (for map with pix-count-percent ?)

import (
	"image/color"
	"math"
	"math/rand"
	"sort"
	"strconv"

	"github.com/improvemedia/colorcandy.git/candy"
)

type Config struct {
	BaseColorsStr       []string          `json:"base_colors"`
	BaseColors          map[string]Color  `json:"-"`
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

	config.BaseColors = map[string]Color{}
	for i, c := range config.BaseColorsStr {
		config.BaseColors[strconv.Itoa(i)] = ColorFromString(c)
	}

	config.ClusterColors = map[Color]Color{}
	for k, v := range config.ClusterColorsStr {
		config.ClusterColors[ColorFromString(k)] = ColorFromString(v)
	}

	return &ColorCandy{config}
}

func (colorCandy *ColorCandy) ExtractColors(path string) (map[string]*candy.ColorMeta, map[string]*candy.ColorCount, error) {
	histogram := CompactToCommonColors(ImageHistogram(path))

	colors := map[string]*candy.ColorMeta{}
	colorsHex := map[string]*candy.ColorCount{}

	for color, count := range histogram {
		cluster, delta := colorCandy.closestColorTo(color)
		hexColor := color.Hex()
		var id string
		for i, v := range colorCandy.BaseColors { // FIXME:(Alexander Yunin): if defined?(Rails) { SearchColor.find_or_create_by(color: color).id }
			if v == cluster {
				id = i
				break
			}
		}

		colorCount := &candy.ColorCount{
			Total:      count.Total,
			Percentage: count.Percentage,
		}
		colorsHex["#"+hexColor] = colorCount

		if oldMeta, found := colors[id]; found {
			oldMeta.OriginalColor["#"+hexColor] = colorCount
			oldMeta.SearchFactor += count.Percentage
			colors[id] = oldMeta
		} else {
			colors[id] = &candy.ColorMeta{
				Id:            id,
				SearchFactor:  count.Percentage,
				Distance:      delta,
				Hex:           hexColor,
				OriginalColor: map[string]*candy.ColorCount{"#" + hexColor: colorCount},
				HexOfBase:     colorCandy.BaseColors[id].Hex(),
			}
		}
	}

	return colors, colorsHex, nil
}

func (colorCandy *ColorCandy) CreatePalette(colors map[string]*candy.ColorCount) map[string]*candy.ColorCount {

	for len(colors) > colorCandy.PaletteColorsMaxNum {
		colorsArr := []*ColorCount{}
		for k, v := range colors {
			colorsArr = append(colorsArr, &ColorCount{
				color:      ColorFromString(k),
				Total:      v.Total,
				Percentage: v.Percentage,
			})
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
		var min float64 = 100001
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
		colors[add.color.Hex()] = &candy.ColorCount{
			Total:      add.Total,
			Percentage: add.Percentage,
		}
		delete(colors, remove.color.Hex())
	}
	if len(colors) == colorCandy.PaletteColorsMaxNum {
		return colors
	}
	for len(colors) < colorCandy.PaletteColorsMaxNum {
		colorsArr := []*ColorCount{}
		for k, v := range colors {
			colorsArr = append(colorsArr, &ColorCount{
				color:      ColorFromString(k),
				Total:      v.Total,
				Percentage: v.Percentage,
			})
		}
		count := colorsArr[rand.Intn(len(colorsArr))]
		r, g, b, _ := count.color.RGBA()
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
			newColor := Color{r + uint32(shift), g + uint32(shift), b + uint32(shift)}
			colors[newColor.Hex()] = &candy.ColorCount{1, 2}
		}
	}

	return colors
}

func (colorCandy *ColorCandy) closestColorTo(c color.Color) (color.Color, float64) {
	var closestColor color.Color
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
