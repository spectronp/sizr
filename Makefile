VERSION=`git describe --always`
HELP_MESSAGE=`cat help.txt`
GO_MODULE="github.com/spectronp/sizr"

build:
	go build -ldflags=" \
		-X 'main.VERSION=${VERSION}' \
		-X 'main.HELP_MESSAGE=${HELP_MESSAGE}' \
		-X '${GO_MODULE}/vars.BASEDIR=/usr/share/sizr' \
		-X '${GO_MODULE}/vars.DB_FILE=/var/lib/sizr/db.json'" .
	touch db.json
	echo '{}' > db.json

clean:
	rm sizr
	rm db.json

install:
	sudo install --owner=root --group=root -m755 sizr /usr/bin/sizr
	sudo install --owner=root --group=root -dm755 /usr/share/sizr
		sudo cp -r scripts /usr/share/sizr/scripts
		sudo chown -R root:root /usr/share/sizr
		sudo chmod -R 755 /usr/share/sizr
	sudo install --owner=root --group=root -Dm644 db.json /var/lib/sizr/db.json
	sudo setfacl -m g:wheel:rw /var/lib/sizr/db.json

uninstall:
	sudo rm /usr/bin/sizr
	sudo rm -rf /usr/share/sizr
	sudo rm -rf /var/lib/sizr

test:
	./tests/wrap.sh -count=1

