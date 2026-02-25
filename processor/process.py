import sys
import rasterio
from rasterio.windows import from_bounds
from rasterio.warp import transform
import numpy as np
from PIL import Image
import uuid
import os

BASE_DIR = os.path.dirname(os.path.abspath(__file__))
file_path = os.path.join(BASE_DIR, "..", "data", "o41078a5.tif")

# Read lat/lng from command line
lat = float(sys.argv[1])
lng = float(sys.argv[2])

with rasterio.open(file_path) as dataset:
    xs, ys = transform(
        "EPSG:4326",
        dataset.crs,
        [lng],
        [lat]
    )

    x = xs[0]
    y = ys[0]

    buffer = 500
    left = x - buffer
    right = x + buffer
    bottom = y - buffer
    top = y + buffer

    window = from_bounds(left, bottom, right, top, dataset.transform)
    cropped = dataset.read(1, window=window)

    if cropped.size == 0:
        print("ERROR: Outside bounds")
        sys.exit(1)

    cropped_normalized = (
        (cropped - cropped.min()) /
        (cropped.max() - cropped.min()) * 255
    ).astype(np.uint8)

    unique_name = f"{uuid.uuid4()}.png"
    output_path = os.path.join(BASE_DIR, "..", "data", unique_name)

    image = Image.fromarray(cropped_normalized)
    image.save(output_path)

    print(unique_name)