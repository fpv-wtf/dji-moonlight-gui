# dji-moonlight-gui

Stream games via Moonlight and [fpv.wtf](https://github.com/fpv-wtf) to your DJI
FPV Goggles!

![splash](media/logo.png)

The DJI Moonlight project is made up of three parts:

- **[dji-moonlight-shim](https://github.com/fpv-wtf/dji-moonlight-shim)**: a
  goggle-side app that displays a video stream coming in over USB.
- **[dji-moonlight-gui](https://github.com/fpv-wtf/dji-moonlight-gui)**: a
  Windows app that streams games to the shim via Moonlight and friends. _You are
  here._
- [dji-moonlight-embedded](https://github.com/fpv-wtf/dji-moonlight-embedded): a
  fork of Moonlight Embedded that can stream to the shim. The GUI app uses this
  internally.

Latency is good, in the 7-14ms range at 120Hz (w/ 5900X + 3080Ti via GeForce
Experience).

![latency](media/latency.gif)

---

![splash](media/screenshot.png)

## Usage

### Setup

1. Install [dji-moonlight-shim](https://github.com/fpv-wtf/dji-moonlight-shim)
   on your goggles.
2. Download the [latest
   release](https://github.com/fpv-wtf/dji-moonlight-gui/releases/latest) and
   extract it.
3. Install and configure the host streaming software.

   _NVIDIA GPU?_ See [Setup: NVIDIA GPU](#setup-nvidia-gpu).

   _AMD/Intel GPU?_ See [Setup: AMD/Intel GPU](#setup-amdintel-gpu).

#### Setup: NVIDIA GPU

1. Install [GeForce
   Experience](https://www.nvidia.com/en-us/geforce/geforce-experience/). You
   almost definitely already have this installed if you've ever updated your
   drivers.
2. Open GeForce Experience and go to _Settings_ > _SHIELD_.
3. Enable _GameStream_.

   ![nvidia_1](media/nvidia_1.png)

4. GameStream will automatically detect a bunch of your games (including Steam),
   though they won't be shown in this tab. You can be more precise by adding
   games (i.e., your sim of choice) here manually.

   Want to stream everything on your desktop directly? Add
   `C:\Windows\System32\mstsc.exe` as a "game" while you're here. Selecting
   `mstsc` in the GUI will then do just that.

#### Setup: AMD/Intel GPU

See [Sunshine
documentation](https://docs.lizardbyte.dev/projects/sunshine/en/latest/) for
more guidance.

1. Download and install
   [Sunshine](https://github.com/LizardByte/Sunshine/releases/tag/latest).
2. Sunshine runs as a background service automatically and uses a locally-hosted
   web UI for settings.

   ![nvidia_1](media/sunshine_1.png)

   Go to [https://localhost:47990/](https://localhost:47990/) and set a username
   and password for future settings fanangling.

### Pairing

Before you can start streaming, you need to pair the GUI app with the host
streaming software. This only needs to be done once.

1. Run `dji-moonlight-gui.exe`.
2. Press _Pair_.
3. The PIN will be displayed in the console. This is what you'll need to enter
   in the host streaming software.

   _GeForce Experience_: A popup will appear on your PC asking you to enter the
   pin.

   _Sunshine_: Go to the [PIN tab](https://localhost:47990/pin) on the web UI
   and enter the pin.

### Streaming

1. Run `dji-moonlight-gui.exe`.
2. Configure the settings to your liking.
3. Press _Refresh_ to fetch the list of your games.
4. Select a game and press _Stream_.

## Modes: BULK vs. RNDIS

There are two streaming modes to choose from: BULK and RNDIS. What does that
mean?

**RNDIS**

This is the default mode. When the goggles are connected to your PC via USB,
they appear as a plain old network interface. If you've ever used your phone as
a hotspot over USB, then you've used this exact same mechanism before. Since
it's just a network interface, all regular networking conventions apply and we
can send data to the goggles like any other device.

_The main downside is that it's slow!_ Due to reasons yet unknown, the maximum
bitrate we're able to achieve is around 30Mbps before packet loss starts to
creep in.

**BULK**

This mode is more experimental. Rather than using this indirect network route,
we can instead send data _directly_ to the goggles via the USB interface. With
this, we can easily achieve a bitrate of 100Mbps (as long as your PC can keep
up).

_The only downside is: driver shenanigans!_

### Using BULK mode

1. Use [the fpv.wtf Driver
   Installer](https://github.com/fpv-wtf/driver-installer) to install the
   correct drivers for the goggles. You may have used this already when you
   rooted your goggles.
2. **Close any fpv.wtf configurator browser tabs**.
3. Run `dji-moonlight-gui.exe`.
4. Select _BULK_ mode.
5. Crank that bitrate.
6. ...carry on as normal.
