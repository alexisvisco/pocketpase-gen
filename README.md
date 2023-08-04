# pocketpase-gen

This is a poc for code gen of models of pocketpase instance. 

## How to use

Install the binary with `go install github.com/pocketbase/pocketbase-gen/cmd/pb-gen@latest`

```bash
pb-gen models --help
Generate models from a pocketbase instance sqlite file

Usage:
  pb-gen models [flags]

Flags:
  -h, --help                    help for models
      --package-folder string   path to package folder (default "modelspb")
      --package-name string     name of the package (default "modelspb")

Global Flags:
      --db-path string   path to pocketbase instance sqlite file (default "pb_data/data.db")
      --verbose          enable verbose mode
```

Runnin `pb-gen models` will create a folder `modelspb` with all the models of the pocketbase instance.
By default the db-path is `pb_data/data.db`.

## Currently supported things

- Getter
- Setter

## What can be imagined

- Better documentation for each fields (validations description, if required ...)
- Better documentation for model (rules ...)
- Expand relations generation
- Enums generations for select fields (field name as suffix ?)
- Hook models generation (beforeSave, afterSave, beforeDelete, afterDelete but with the model as parameter)
- Add/Remove for arrayable fields
- Add automated tests


