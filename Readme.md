# djifpvvideoout

This is a straightforward port of the dji fpv videooutscript (from [here](https://github.com/fpv-wtf/voc-poc)) to **golang**

The intersting thing is that it seems to **start working in low power mode** as well

Also this script should support unplugging and replugging USB, disconnecting / reconnecting the drone and should not care if it is started while the goggles are off

FFMpeg needs to be available in your path, it will be started by the go binary automatically

## Running

I usually run this like so:
go run .

Crosscompilation does not work in an easy way due to the dependency on libusb

If you find any new intersting things about this let me know on the discordserver @s00500 or here in the issue section


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
