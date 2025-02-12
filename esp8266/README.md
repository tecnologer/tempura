# Platformio

This guide will help you set up and build projects using [Platformio][1] by defining necessary environment variables in a `.env` file.

## Installation

You can install Platformio on Visual Studio Code or Atom using the links below:

- [Visual Studio Code][2]
- [Atom][3]

## Configuration

### Environment Variables

To run your projects, you need to set up an environment file (`.env`) that includes the Wi-Fi credentials and API URL. Create a `.env` file in your project's root directory and populate it with the following content:

```env
# Wi-Fi Configuration
WIFI_SSID=TheNameOfYourWifi
WIFI_PASSWORD=YourWifiPassword

# API Configuration
API_URL=http://localhost:88/v1/records

```

**Note:** Ensure that your Wi-Fi password does not contain special characters that might be misinterpreted by your system or the application.

[1]: https://platformio.org/
[2]: https://code.visualstudio.com/
[3]: https://atom-editor.cc/
