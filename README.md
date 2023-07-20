# endgame

This repo holds the code for all the testing we do on a web asset after our vulnerability scanner.

```
Note: All the new scanning functionality added should be appended under `scannerlist.go` file
```

### Docker build command

```bash
docker build . -t registry.digitalocean.com/getastra/endgame:<version to build>
```

### How to use with config.json
Working with this service can be done in two ways: -
1. Create a file name `config.json` from `default_config.json` file and declare all the required data in it
   1. `export SCAN_ID="<SCAN_ID>" && go run ./*.go -local`

2. Export all the required data as environment variables
   1. `export SCAN_ID="<SCAN_ID>" && go run ./*.go`

