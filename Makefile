all: kacg

kacg: 
	go build

clean:
	- rm -f kacg
