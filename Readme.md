# djifpvvideoout

This is a straightforward port of the dji fpv videooutscript (from [here](https://github.com/fpv-wtf/voc-poc)) to **golang**

The intersting thing is that it seems to **start working in low power mode** as well

## Running

I usually run this like so:

go run . | ffplay -i - -fast -flags2 fast -fflags nobuffer -flags low_delay -strict experimental -vf "setpts=N/60/TB" -framedrop -sync ext -probesize 32 -analyzeduration 0


Crosscompilation does not work in an easy way due to the dependency on libusb

If you find any new intersting things about this let me know on the discordserver @s00500 or here in the issue section


## More usefull Links

- https://github.com/fpv-wtf/voc-poc
- https://gist.github.com/fichek/c69326dba7e5a9dfb6ecc2c9e4e93224
- https://github.com/district-michael/fpv_live


Greetings,
Lukas

My Website: [lbsfilm.at](lbsfilm.at)

[Buy me a coffee ☕️](https://www.paypal.com/paypalme/lukasbachschwell/3)
