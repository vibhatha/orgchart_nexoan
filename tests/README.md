To run the tests run:

```bash
go test ./tests
```

For developers:

If you change the api code (ie Nexoan) but not the tests code. Run the following to execute the tests without caching the previous test results:

```bash
go test ./tests -count=1
```
