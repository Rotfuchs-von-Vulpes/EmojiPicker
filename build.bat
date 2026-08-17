rsrc -ico EmojiPicker.ico
go build -ldflags "-H windowsgui -linkmode external -extldflags '-static-libgcc -static-libstdc++'" .