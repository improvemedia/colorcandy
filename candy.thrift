struct ColorCount {
  1: i64 total,
  2: double percentage
}

struct ColorMeta {
  1: string id,
  2: double search_factor,
  3: double distance,
  4: string hex,
  5: map<string, ColorCount> original_color,
  6: string hex_of_base
}

service Candy {
  map<string,ColorMeta> candify(1: string url)
}
