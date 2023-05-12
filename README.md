# go-cursor-paginator
For paginating results in golang. It doesn't rely on database



## Usage/Examples

see the [examples folder](/examples)


## Running Tests
To run all tests, run the following command

```bash
  make
```

## Running Tests (when makefile not available)
### OR Do this to run Test
```bash
    go test -v  ./...
```

### runs coverage tests and generates the coverage report
```bash
    go test ./... -v -coverpkg=./...
```

### runs integration tests
```bash
	go test ./... -tags=integration ./...
```

## License

[MIT](https://choosealicense.com/licenses/mit/)

