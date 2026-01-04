# Poppins Font Files

This directory should contain the Poppins font files for the application.

## Required Font Files

Download the following font files from [Google Fonts](https://fonts.google.com/specimen/Poppins):

- `Poppins-Regular.ttf` (weight: 400)
- `Poppins-Medium.ttf` (weight: 500)
- `Poppins-SemiBold.ttf` (weight: 600)
- `Poppins-Bold.ttf` (weight: 700)

## Alternative: Using google_fonts Package

The app is configured to use the `google_fonts` package which can load Poppins dynamically from the internet. This means the app will work even without bundled font files.

However, for production apps, it's recommended to bundle the fonts locally for:
- Faster initial load times
- Offline support
- Consistent typography

## How to Download

1. Visit https://fonts.google.com/specimen/Poppins
2. Click "Download family"
3. Extract the ZIP file
4. Copy the required .ttf files to this directory

## License

Poppins is licensed under the [Open Font License](https://scripts.sil.org/OFL).
