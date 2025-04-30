# RSS Radio Generator

This project is a simple web application that allows users to generate an RSS feed URL for internet radio stations. Users can input multiple stations along with their titles, URLs, and descriptions, and the application will generate a valid RSS feed URL based on the provided information.

## Project Structure

```
rss-radio-generator
├── src
│   ├── main.go          # Entry point of the application
│   ├── static
│   │   ├── index.html   # HTML form for generating the RSS URL
│   │   └── script.js     # JavaScript for managing form functionality
├── README.md            # Documentation for the project
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd rss-radio-generator
   ```

2. **Navigate to the `src` directory:**
   ```
   cd src
   ```

3. **Run the application:**
   ```
   go run main.go
   ```

4. **Access the application:**
   Open your web browser and go to `http://localhost:8080`.

## Usage

1. Fill in the RSS title, description, and language in the provided fields.
2. Add multiple stations by clicking the "+" button. Each station requires a title, URL, and description.
3. Once all fields are filled, the generated RSS URL will be displayed.
4. Use the "Copy URL" button to copy the generated URL to your clipboard for easy sharing.

## Contributing

Feel free to submit issues or pull requests if you have suggestions or improvements for the project.