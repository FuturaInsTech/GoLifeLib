@echo off

REM Get the folder where the BAT file resides
set "BASEDIR=%~dp0"

REM Add Tesseract DLL folder to PATH
set "PATH=%BASEDIR%tesslib;%PATH%"

REM Tell Tesseract where trained data files are located
set "TESSDATA_PREFIX=%BASEDIR%tessdata"

REM Execute program and pass all arguments
"%BASEDIR%pngtohtml.exe" %*