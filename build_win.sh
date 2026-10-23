CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o ./build/win/runner.exe -ldflags "-s -w" ./runner
CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o ./build/win/editor.exe -ldflags "-s -w" ./editor
cp -r ./editor/templates ./build/win
cp ./editor/chr_default.png ./build/win
cp ./editor/pal_default.png ./build/win
cp ./editor/keybinds.json ./build/win