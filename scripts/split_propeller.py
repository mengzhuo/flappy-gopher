from pathlib import Path
from PIL import Image, ImageDraw, ImageOps

ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / "resources" / "sprites" / "gopher.png"
BODY = ROOT / "resources" / "sprites" / "gopher-body.png"
PROPELLER = ROOT / "resources" / "sprites" / "gopher-propeller.png"
SIZE = (880, 624)

# Polygons enclose the two complete hand-drawn blades and their motion arcs.
# The hub at (766, 390) and the engine behind it deliberately stay in the body.
BLADES_AND_ARCS = [
    [(766, 382), (758, 354), (775, 320), (797, 286), (823, 255),
     (842, 250), (856, 260), (860, 278), (853, 300), (835, 327),
     (810, 354), (784, 382)],
    [(759, 407), (744, 425), (725, 451), (705, 480), (688, 507),
     (686, 525), (697, 539), (714, 541), (731, 527), (747, 502),
     (762, 472), (776, 438), (784, 410), (779, 402)],
    [(820, 285), (879, 285), (879, 450), (820, 450)],
]

source = Image.open(SOURCE).convert("RGBA")
if source.size != SIZE:
    raise RuntimeError(f"unexpected source size: {source.size}")
if source.getchannel("A").getbbox() is None:
    raise RuntimeError("source sprite is empty")

mask = Image.new("L", SIZE, 0)
draw = ImageDraw.Draw(mask)
for polygon in BLADES_AND_ARCS:
    draw.polygon(polygon, fill=255)
draw.ellipse((718, 350, 800, 420), fill=0)

body = Image.new("RGBA", SIZE, (0, 0, 0, 0))
propeller = Image.new("RGBA", SIZE, (0, 0, 0, 0))
body.paste(source, (0, 0), ImageOps.invert(mask))
propeller.paste(source, (0, 0), mask)

for name, image in (("body", body), ("propeller", propeller)):
    if image.getchannel("A").getbbox() is None:
        raise RuntimeError(f"{name} output is empty")
if propeller.getpixel((766, 390))[3] != 0:
    raise RuntimeError("propeller mask covers the hub")
if body.getpixel((766, 390)) != source.getpixel((766, 390)):
    raise RuntimeError("body does not preserve the hub")

body.save(BODY, format="PNG", optimize=True)
propeller.save(PROPELLER, format="PNG", optimize=True)
print(f"{BODY}: {body.size}")
print(f"{PROPELLER}: {propeller.size}")
