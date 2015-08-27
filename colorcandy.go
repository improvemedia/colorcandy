package colorcandy

// TODO(Alexander Yunin): make color type for SRGB
// add func to convert 16bit colors to 8bit
// use new types for something else (for map with pix-count-percent ?)

import (
	"image/color"
	"math"

	"github.com/improvemedia/colorcandy.git/candy"
)

type Config struct {
	BaseColorsStr       []string `json:"base_colors"`
	BaseColors          map[string]Color
	ClusterColorsStr    map[string]string `json:"cluster_colors"`
	ClusterColors       map[Color]Color
	ColorsCount         uint    `json:"colors_count"`
	PaletteColorsMaxNum int     `json:"palette_colors_max_num"`
	WhiteThreshold      int     `json:"white_threshold"`
	BlackThreshold      int     `json:"black_threshold"`
	Delta               float64 `json:"delta"`
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
		config.BaseColors[string(i)] = ColorFromString(c)
	}

	config.ClusterColors = map[Color]Color{}
	for k, v := range config.ClusterColorsStr {
		config.ClusterColors[ColorFromString(k)] = ColorFromString(v)
	}

	return &ColorCandy{config}
}

func (colorCandy *ColorCandy) Candify(path string) (map[string]*candy.ColorMeta, error) {
	const delta float64 = 2.5

	delta_count := 0
	mapped_palette := ImageHistogram(path)

	cp_mapped_palette := mapped_palette
	new_palette := map[Color]*ColorCount{}

	// cant iterate through map values inside a cycle after assigning new values
	// I have to launch new cycle or will get mapped_palette[k] -> nil %)
	for color1, v := range mapped_palette {
		common_colors := []*ColorCount{v}
		lab := Lab.Convert(color1)

		for color2, v2 := range cp_mapped_palette {
			lab2 := Lab.Convert(color2)

			if delta > DeltaE(lab, lab2) {
				delta_count += 1 //debug
				if color1 != color2 {
					common_colors = append(common_colors, v2)
					common_colors[0].Percentage += v2.Percentage
				}
			}
		}

		for _, count := range common_colors[1:] {
			delete(cp_mapped_palette, count.color)
		}

		new_palette[common_colors[0].color] = common_colors[0]
	}

	colors := map[string]*candy.ColorMeta{}
	//colors := make(map[int]color_meta)
	for color, count := range new_palette {
		cluster, delta := colorCandy.ClosestColorTo(color)
		hexColor := color.Hex()
		var id string
		for i, v := range colorCandy.BaseColors { // FIXME:(Alexander Yunin): if defined?(Rails) { SearchColor.find_or_create_by(color: color).id }
			if v == cluster {
				id = string(i)
				break
			}
		}

		colorCount := &candy.ColorCount{
			Total:      int64(count.Total),
			Percentage: count.Percentage,
		}
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

	return colors, nil
}

func (colorCandy *ColorCandy) ClosestColorTo(c color.Color) (color.Color, float64) {
	var closest_color color.Color
	var cluster Color
	min_delta := math.MaxFloat64 // pls no buf overflow

	lab := Lab.Convert(c)
	for k, _ := range colorCandy.ClusterColors {
		delta := DeltaE(Lab.Convert(k), lab)
		if delta < min_delta {
			min_delta = delta
			closest_color = k
			cluster = colorCandy.ClusterColors[k]
		}
	}
	return cluster, DeltaE(Lab.Convert(cluster), Lab.Convert(closest_color))
}
