# Wistia 1080p Downloader
This tool enables downloading 1080p MP4 videos from Wistia pages. It extracts the video ID, retrieves JSON data, detects the 1080p video URL, and downloads the video locally.
## Features
- : Extracts the video ID via the `wvideo` parameter.
- : Retrieves JSON data from the iframe and completes it if necessary.
- : Finds the 1080p video URL using data from the JSON.

# Wistia 1080p Downloader

This application is designed to download 1080p resolution MP4 files from Wistia video pages. It extracts the video ID from a given Wistia course or video URL, retrieves and parses the JSON data from the corresponding iframe page, detects the URL of the 1080p video file, and downloads it locally.

## Features

- **Video ID Extraction**  
  Extracts the Wistia video ID from the URL by searching for the `wvideo` parameter.

- **JSON Data Extraction**  
  Retrieves the JSON data from the Wistia iframe page using regex. If the JSON string is missing closing braces, the `balanceBraces` function completes it.

- **1080p Video URL Detection**  
  Parses the JSON to find the video asset with a height of 1080 pixels. If the asset URL ends with `.bin`, it converts it to an `.mp4` URL.

- **File Download**  
  Downloads the video file from the detected URL and saves it locally with a specified file name.

## Requirements

- **Go Programming Language**  
  The project requires Go (version 1.15 or higher is recommended).

- **Internet Connection**  
  Needed to access Wistia servers and download video files.

## Installation and Running

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   
2. **Build the Application**
   ```bash
    go build -o wistia_downloader
    
3. **Alternatively, you can run the application directly with**
    ```bash
   go run main.go