# djifpvvideoout

This is a port of the dji fpv videooutscript (from [here](https://github.com/fpv-wtf/voc-poc)) to **golang**

The intersting thing is that it **starts working in low power mode** as well

Also this script should support unplugging and replugging USB, disconnecting / reconnecting the drone and should not care if it is started while the goggles are off

It will also work with **multiple googles at once**, so if you connect new ones while it is running it will open new ffplay instances for you

FFMpeg (or gstreamer) needs to be available in your path, it will be started by the go binary automatically

libusb-1.0-dev needs to be available if you want to build this from source

## Running

Make sure you have go and ffmpeg installed, then run like this:

`go run .`

Per default this will use ffplay, if you want to use gstreamer though use it like this `go run . --output gstreamer`

The sync option is usefull when running on raspberrypi and directly using the framebuffer (when not using the desktop environment) It seems to work better

Other options are **fifo** and **ffplay** (default)

Crosscompilation does not work in an easy way due to the dependency on libusb
If you find any new intersting things about this let me know on the discordserver @s00500 or here in the issue section

# Using with OBS

The best way to work with OBS (at least on mac) I found was by creating a fifo, this is basically a special file where I pipe in the data from the googles. Then the obs-gstreamer plugin can be used
to read and decode the stream from this file and directly open it in OBS.

Install th eOBS gstreamer plugin from here https://github.com/fzwoch/obs-gstreamer

On macOS you can use the prebuilt binary, but carefull: this requires gstreamer to be installed via macports, not brew.

If you installed it with brew you will need to copy over the binaries to where the plugin expects them (TODO: Document paths)

Then create a new gstreamer source in obs and enter this command:

filesrc location=/Path/To/your/stream.fifo ! decodebin3  ! videoconvert n-threads=8 ! video.

if at some point your OBS does not start anymore make sure to remove any fifo files left, that should fix it
## Test status

For now I tested this with macOS bigsur, no additional drivers installd, although I might have a few on my computer anyway

Used Caddx Vista and GooglesV2

Next thing todo is testing this on a Raspberry Pi

Also I found that this works with a regular USB 2 cable (no need for the usb 3 one)
## More usefull Links

- https://github.com/fpv-wtf/voc-poc
- https://gist.github.com/fichek/c69326dba7e5a9dfb6ecc2c9e4e93224
- https://github.com/district-michael/fpv_live


Greetings,
Lukas

My Website: [lbsfilm.at](lbsfilm.at)

[Buy me a coffee ☕️](https://www.paypal.com/paypalme/lukasbachschwell/3)
