@echo off
setlocal enabledelayedexpansion

:: Read version from version.txt and remove any spaces
set /p version=<version.txt
set version=%version: =%

:: Extract major and minor version numbers
for /f "tokens=1,2 delims=." %%a in ("%version%") do (
    set major_version=%%a
    set minor_version=%%b
)

:: Increment the minor version
set /a new_minor_version=minor_version + 1
set new_version=%major_version%.%new_minor_version%

:: Write the new version to version.txt
echo %new_version% > version.txt

:: Build the Wails app with the versioned filename (ensure no spaces)
set output_file=cefp_scifi_v%version%.exe
wails build -o "%output_file%"

echo Build complete: %output_file%
