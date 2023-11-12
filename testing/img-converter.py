import base64

with open("1.jpg", "rb") as imageFile:
    str = base64.b64encode(imageFile.read())
    print(str)
