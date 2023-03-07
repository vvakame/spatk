# spatk

Spanner ToolKit

## sidx

Spanner InDeX helper generator.

```shell
$ go run github.com/vvakame/spatk/cmd/sidx -output model_spanner_index.go ./db/schema.sql
```

## sig

Spanner Information code Generator.

```shell
$ go run github.com/vvakame/spatk/cmd/sig -private -output model_spanner_info.go .
```

## scur

Spanner CURsor constructor.

## sqb

Spanner Query Builder.

## How to test

Do this once.

```shell
$ export CLOUDSDK_ACTIVE_CONFIG_NAME=spatk-config
$ gcloud config configurations create --no-activate $CLOUDSDK_ACTIVE_CONFIG_NAME
$ gcloud config set auth/disable_credentials true
$ gcloud config set project spatk-test
$ gcloud config set api_endpoint_overrides/spanner http://localhost:9020/
```

Do this before runs test.

```shell
$ docker compose up -d --build --force-recreate
$ export CLOUDSDK_ACTIVE_CONFIG_NAME=spatk-config
$ export SPANNER_EMULATOR_INSTANCE_NAME=spatk-test-instance
$ gcloud spanner instances create "${SPANNER_EMULATOR_INSTANCE_NAME}" --config=emulator-config --nodes=1 --description "for testing"
$ export SPANNER_EMULATOR_HOST=localhost:9010
$ export SPANNER_EMULATOR_PROJECT_ID=spatk-test
$ export SPANNER_EMULATOR_INSTANCE_NAME=spatk-test-instance
$ export SPANNER_EMULATOR_DB_NAME=spatk-test-db
$ go test -v ./...
$ docker compose down
```
