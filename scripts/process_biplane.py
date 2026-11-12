from collections import deque
from pathlib import Path
from PIL import Image

ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / ".design-work" / "candidate-biplane.jpg"
OUTPUT = ROOT / "resources" / "sprites" / "gopher.png"
CROP = (32, 128, 912, 752)  # 880 × 624; excludes separate motion marks above the plane.

source = Image.open(SOURCE).convert("RGB")
gray = source.convert("L")
width, height = source.size
values = list(gray.getdata())
white = bytearray(value >= 232 for value in values)
seen = bytearray(width * height)
region = bytearray(width * height)  # 0=line/antialias, 1=outer background, 2=enclosed white.
tint = bytearray(width * height)    # 0=off-white, 1=Go blue, 2=warm Gopher.

for start in range(width * height):
    if not white[start] or seen[start]:
        continue
    queue = deque([start])
    seen[start] = 1
    component = []
    touches_border = False
    sum_x = sum_y = 0
    while queue:
        index = queue.popleft()
        component.append(index)
        y, x = divmod(index, width)
        sum_x += x
        sum_y += y
        if x == 0 or y == 0 or x == width - 1 or y == height - 1:
            touches_border = True
        neighbors = []
        if x > 0:
            neighbors.append(index - 1)
        if x + 1 < width:
            neighbors.append(index + 1)
        if y > 0:
            neighbors.append(index - width)
        if y + 1 < height:
            neighbors.append(index + width)
        for neighbor in neighbors:
            if white[neighbor] and not seen[neighbor]:
                seen[neighbor] = 1
                queue.append(neighbor)
    kind = 1 if touches_border else 2
    area = len(component)
    center_x = sum_x / area
    center_y = sum_y / area
    tint_kind = 0
    if kind == 2 and 340 <= center_x <= 700 and 160 <= center_y <= 380 and area >= 120:
        tint_kind = 2
    elif kind == 2 and area >= 500:
        tint_kind = 1
    for index in component:
        region[index] = kind
        tint[index] = tint_kind

result = Image.new("RGBA", source.size)
pixels = result.load()
for y in range(height):
    for x in range(width):
        index = y * width + x
        value = values[index]
        if region[index] == 1:
            pixels[x, y] = (0, 0, 0, 0)
        elif region[index] == 2:
            color = (245, 250, 248)
            if tint[index] == 1:
                color = (170, 231, 242)
            elif tint[index] == 2:
                color = (246, 207, 137)
            pixels[x, y] = (*color, 255)
        else:
            alpha = max(0, min(255, (235 - value) * 8))
            pixels[x, y] = (23, 57, 74, alpha)

# Remove the separate hand-drawn motion marks above the tail; the game draws its own trail.
for y in range(170):
    for x in range(700, width):
        pixels[x, y] = (0, 0, 0, 0)

result = result.crop(CROP)
if result.size != (880, 624):
    raise RuntimeError(f"unexpected crop size: {result.size}")
if result.getchannel("A").getbbox() is None:
    raise RuntimeError("background removal erased the complete image")
OUTPUT.parent.mkdir(parents=True, exist_ok=True)
result.save(OUTPUT, format="PNG", optimize=True)
print(f"{OUTPUT}: {result.size[0]}x{result.size[1]}, alpha={result.getchannel('A').getextrema()}")
