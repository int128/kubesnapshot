all: kubesnapshot

kubesnapshot: *.go */*.go
	go build -o $@
