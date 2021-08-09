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
