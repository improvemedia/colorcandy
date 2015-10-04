require 'rmagick'
require 'json'

module Magick
  class Pixel
    def to_s
      r = self.red
      b = self.blue
      g = self.green
      r /= 257 if r / 255 > 0
      b /= 257 if b / 255 > 0
      g /= 257 if g / 255 > 0
      "#" + [r, g, b].pack("C*").unpack("H*")[0]
    end
  end
end

image = ::Magick::ImageList.new(ARGV[0])
image_quantized = image.quantize(60, Magick::YIQColorspace)
palette = image_quantized.color_histogram
image_quantized.destroy!
image.destroy!

sum = palette.values.inject(:+)

result = {}
palette.each do |k, v|
  result[k.to_s] = [v, v / (sum / 100.0)]
end

puts JSON.pretty_generate(result)
