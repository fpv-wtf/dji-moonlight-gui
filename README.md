# dji-moonlight-gui

Stream games via Moonlight and [fpv.wtf](https://github.com/fpv-wtf) to your DJI
FPV Goggles. This is the PC side of the whole thing. See also:

- [dji-moonlight-shim](https://github.com/fpv-wtf/dji-moonlight-shim)
- [dji-moonlight-embedded](https://github.com/fpv-wtf/dji-moonlight-embedded)

![splash](media/screenshot.png)

## Usage

### Setup

1. Install and configure host streaming software.

   _NVIDIA GPU_: Install [GeForce
   Experience](https://www.nvidia.com/en-us/geforce/geforce-experience/) and
   enable _GAMESTREAM_ under _Settings_ > _SHIELD_.

   _Otherwise_: Install and configure
   [Sunshine](https://github.com/LizardByte/Sunshine/).

1. Install [dji-moonlight-shim](https://github.com/fpv-wtf/dji-moonlight-shim)
   on your goggles.
2. Download the [latest
   release](https://github.com/fpv-wtf/dji-moonlight-gui/releases/latest) and
   extract it.

### Running

The GUI is just a wrapper around
[dji-moonlight-embedded](https://github.com/fpv-wtf/dji-moonlight-embedded).
Keep an eye on the console output for anything weird.

1. Run `dji-moonlight-gui.exe`.
2. Pair Moonlight with your streaming software by pressing _Pair_. The PIN will
   be displayed in the console. This only needs to be done once.

   _GeForce Experience_: A popup will appear on your PC asking you to confirm
   the pin.

   _Sunshine_: You will need to go to the _Settings_ page in your browser and
    enter the pin.
3. Press _Refresh_ to list your games.
4. Select a game and press _Stream_.

_Protip_: You can stream your desktop by adding `C:\Windows\System32\mstsc.exe`
to the list of games in GeForce Experience. You need to press _Quit_ manually to
stop streaming, though.
