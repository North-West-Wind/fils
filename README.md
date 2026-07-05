# FILS (F*ck It, Link Shortener)

The simplest link shortener you can ask for.

This is me after 3 hours of frustration trying to selfhost a "feature-rich" link shortener.  
Turns out it's better to make my own.

## Install
I don't know if I'm bothered enough to publish the binary as a release.
If there is, download from there.

## Building
This link shortener is written in Go. You must have Go installed on your computer to compile it.

```sh
git clone https://codeberg.org/NorthWestWind/fils # download the repository
cd fils                                           # change directory
bash compile.sh                                   # compile
```

The `compile.sh` script handles minifying client files.
It uses `npx`, which is a part of `npm`.
If `npm` is not found, it will default to just copying the files.
Remove it if you don't want to use `npm`.

## Usage
Run the binary file:
```sh
./fils # or "fils" if you put it in $PATH
```

Configuration is done via environment variables.
Either use a `.env` file or pass them directly to the program.

`.env` example:
```sh
PORT=3456              # (default: 3000) port of the web server
PASSPHRASE=very-secret # (default: empty) passphrase to use when shortening links, omit to allow everyone to use
STORE=store.txt        # (default: store.txt) file where links will be stored
```

Passing directly to program example:
```sh
PORT=3456 fils
```

There's an extra option `APP_ENV`. If you set `APP_ENV=dev`, it will use the disk `client/` directory instead of the embedded ones.
Useful if you want to change the client code.