=============================================================================================================
For Windows System Use...
=============================================================================================================

A. To build  a Smart BI Editor standalone installer... 

run the command in terminal 'npm run dist'

With "win": { "target": "nsis" }, electron-builder typically produces:

dist/
├── SmartBI Editor Setup 1.0.0.exe    ← the installer (this is what you distribute)
├── win-unpacked/                      ← unpacked app (for testing only)
│   └── SmartBI Editor.exe             ← runs directly, no install needed
├── latest.yml                         ← auto-update metadata (ignore for now)

Two ways to Run smarbi-editor...

1. Quick test — run it unpacked (no install)
Go into dist/win-unpacked/ and double-click SmartBI Editor.exe directly. This launches the app exactly as it'll behave once installed, without going through the installer. Good for a fast sanity check.

2. Proper install — run the Setup installer
Double-click SmartBI Editor Setup 1.0.0.exe. This:

Installs the app (usually to %LOCALAPPDATA%\Programs\SmartBI Editor by default with NSIS)
Creates a Start Menu shortcut and optionally a desktop shortcut
Creates an uninstaller entry in "Add or Remove Programs"

After installing, the user launches it from the Start Menu like any normal Windows app — no terminal, no npm start, nothing dev-related needed.

Test this one first, not the unpacked version, since the installer path is what real users will actually experience — you want to catch any path-resolution issues (like your node_modules/grapesjs relative paths) after a real install, not just from the dev folder.

How to distribute it

Just share the single file: SmartBI Editor Setup 1.0.0.exe. Options depending on your audience:

Internal team / few people: upload it to a shared drive, Slack, email, or internal file server — anyone can just download and double-click it.
Wider distribution: host it somewhere with a stable link — GitHub Releases (if the code's on GitHub, even privately), or in company's internal downloads page.
Unsigned binary warning: since you likely haven't code-signed the app, Windows SmartScreen will show a "Windows protected your PC" warning on first run. Users need to click "More info" → "Run anyway." This is normal for unsigned installers, not a bug. If this matters for wider/public distribution, code-signing requires purchasing a certificate.

Check in Environment Variables the PATH variable is having the installed folder path. if not add and restart system.

B. To build  a Smart BI Editor portable single self-contained .exe with no install step... 

execute the command in terminal 'npm run dist'

With "win": { "target": ["portable"] }, electron-builder typically produces:
a single file like dist/SmartBI Editor 1.0.0.exe
copy it to any folder (e.g. your common path), double-click to run, no separate resources folder needed alongside it since portable mode bundles everything into that single file (it self-extracts to a temp location at runtime).

Check in Environment Variables the PATH variable is having the installed folder path. if not add and restart system.

Recommended Option is B

=============================================================================================================
For Linux System Use...
=============================================================================================================

the equivalent of Windows' "portable single file" is AppImage.
1. Add to the existing "build" json key in package.json the following.
 "win": {
      "target": "portable"
    },
    "linux": {
      "target": "AppImage",
      "category": "Utility"
    }

Important: electron-builder builds for the OS it's running on by default. 
You can't cross-build a Linux AppImage from Windows without extra setup (Docker or WSL2).
Execute the command 'npm run dist' from an actual Linux machine or WSL2.
It shall produce..
dist/SmartBI Editor-1.0.0.AppImage

Make dist/SmartBI Editor-1.0.0.AppImage as executable and run the command
chmod +x "SmartBI Editor-1.0.0.AppImage"
./"SmartBI Editor-1.0.0.AppImage"

Move it to a common path and run from anywhere
mkdir -p ~/Applications
mv "SmartBI Editor-1.0.0.AppImage" ~/Applications/

To execute from anywhere...
sudo ln -s ~/Applications/"SmartBI Editor-1.0.0.AppImage" /usr/local/bin/smartbi


End of Set up



