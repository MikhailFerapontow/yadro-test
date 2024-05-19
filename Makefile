image_name = "ferapontov/yadro-test"

make: build

test:
	docker run --rm $(image_name) go test ./...

build:
	docker build -t $(image_name) .

run:
	@docker run --rm \
	--mount type=bind,source=$(realpath $(INPUT)),target=/app/input.txt \
	-w /app \
	$(image_name) ./app input.txt

clean:
	docker rmi $(image_name)