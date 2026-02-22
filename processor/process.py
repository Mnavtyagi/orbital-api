import rasterio
from rasterio.windows import from_bounds
import numpy as np
from PIL import Image

file_path = "../data/o41078a5.tif"

with rasterio.open(file_path) as dataset:
    print("Original Bounds:", dataset.bounds)

    # Define crop area (small box inside original bounds)
    left = dataset.bounds.left + 1000
    right = dataset.bounds.left + 3000
    bottom = dataset.bounds.bottom + 1000
    top = dataset.bounds.bottom + 3000

    window = from_bounds(left, bottom, right, top, dataset.transform)
    cropped = dataset.read(1, window=window)

    print("Cropping region...")

    # Normalize to 0-255
    cropped_normalized = (
        (cropped - cropped.min()) /
        (cropped.max() - cropped.min()) * 255
    ).astype(np.uint8)

    image = Image.fromarray(cropped_normalized)
    image.save("../data/cropped_output.png")

print("Saved cropped_output.png")