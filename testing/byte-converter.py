import io
from PIL import Image

def array_to_image(byte_array):
    # Create a byte stream from the array
    byte_stream = io.BytesIO(byte_array)

    # Open the byte stream as an image using PIL/Pillow
    image = Image.open(byte_stream)

    return image

# Example usage
byte_array = b'\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR...'
image = array_to_image(byte_array)
image.show()  # Display the image