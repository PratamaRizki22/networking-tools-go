go build -o pingtool main.go

sudo setcap cap_net_raw+ep ./pingtool
