#!/bin/zsh

# sed -i.bak '/- deadcode/d' .golangci.yml
# sed -i '' '/- unused/d' .golangci.yml
# sed -i '' '/- structcheck/d' .golangci.yml

echo "Lint"
golangci-lint run ./...
cd ..

# mv .golangci.yml.bak .golangci.yml