struct ColorCount {
  1: i64 total,
  2: double percentage
}

struct ColorMeta {
  1: string color,
  2: string baseColor,
  3: double search_factor,
  4: double distance
}

struct Result {
  1: map<string, ColorMeta> colors,
  2: map<string, ColorCount> palette
}

service Candy {
  Result candify(1: string url, 2: list<string> searchColors)
}
