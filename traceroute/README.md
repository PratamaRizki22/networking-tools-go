go build -o traceroute main.go

sudo setcap cap_net_raw+ep ./traceroute
