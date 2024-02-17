go install
cd ../../example
go build -o example
cd ../cmd/goforever
goforever --http
