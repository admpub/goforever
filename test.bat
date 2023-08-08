cd C:\Users\PC\go\src\github.com\admpub\goforever
go test -v -count=1 -run "TestProcessStartByUser" --user=hank-minipc\test
# go test -v -count=1 -run "TestWindowsToken"
go test -v -count=1 -run "TestGetTokenByPid"
pause