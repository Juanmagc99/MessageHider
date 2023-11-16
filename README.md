# MessageHider
This project allows you to hide messages in images. However, be aware that when using the JPG format to save images, 
there is a potential for information loss due to compression.

## How it Works

1. **Sending Messages:**
   - You need to send a POST request to the server with the content you want to hide.
   - Use the form-data format to attach the image and provide the message to hide.

2. **Request Format:**
   - Key "image": Attach the image you want to process.
   - Key "message": Provide the message you want to hide in the image (must be a string).

3. **Example Request (curl):**
   ```bash
   curl -X POST -H "Content-Type: multipart/form-data" \
        -F "image=@path/to/your/image.jpg" \
        -F "message=Hello, this is my secret message" \
        http://your-server.com/process-image

