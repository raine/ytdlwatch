build:
  go build -o ytdlwatch .

build-arm64:
  GOOS=linux GOARCH=arm64 go build -o ytdlwatch_arm64 .

run-w:
  fd .go | entr -r go run . 

test +FLAGS='./...':
  richgo test {{FLAGS}}

test-w +FLAGS='./...':
  fd .go | entr richgo test {{FLAGS}}
