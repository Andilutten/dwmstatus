# DWM Statusbar

Dynamic and configurable statusbar for dwm written in Go.
It uses rpc to setup a server that updates the root window name.
See [Remotely using it](#Remotely using it) how to use rpc to 
update values on the statusbar manually.

## Getting it

```sh
go get github.com/Andilutten/dwmstatus
```

## Running it

The application should be run as a daemon. Simply
call the following in your .xinitrc or something like it.

```sh
dwmstatus &
```

## Configuring it

Copy the config.yaml into .config/dwmstatus/ to get
a base config. The config file does not really need any
saying. Take a look at it and you will understand whats
going on.

## Remotly updating it

Values can be remotely updated using the -update flag.
Example:
```sh
dwmstatus -update [NAME]
```

NAME should be any of the names mapped in config.yaml.
This command triggers the given statusbar item to force reload.

