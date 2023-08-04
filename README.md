# Product Development // Iconik IO Tool

This tool is a POC that allows you to:

- Input a given CSV into the metadata fields of a given asset in Iconik
- Output the metadata fields of a given asset in Iconik to a CSV



## Installation
To install this tool, please run:

`go get -u github.com/base-media-cloud/pd-iconik-io-rd`



## General Usage

### Build

```go build -o pd-iconik-io-rd app/main.go```

### Usage

There are two modes of use for this command-line tool: `input`, and `output`.



#### Input Mode

Input mode takes a pre-populated CSV file conforming to the respective schema constraints and writes the contained metadata field values contained therein to the corresponding Metadata View for each Asset in iconik based on the supplied Collection ID.

##### Example

```bash
$ pd-iconik-io-rd -input input.csv -app-id <AppID> \
-auth-token <AuthToken> -collection-id <CollectionID> -iconik-url \ 
<IconikURL> -metadata-view-id <ViewID>
```

##### Schema constraints

- First row MUST be a header row.
- R1C1 MUST be `id`.
- R1C2 MUST be `original_name`.
- R1C3 MUST be `title`.
- R1C4 -> R1Cn are the name attributes of the metadata fields in the view you want to manipulate.
- First column MUST always be the UUID of the asset.
- Second column MUST always be the original filename of the asset.
- Third column MUST always be the title of the asset.
- Columns 4->n are the values of the metadata fields in R1.
- If a field can have multiple values (e.g., Tags), they must be comma separated in the appropriate cell.
- If a field is boolean, it must be either true or false.


##### Example CSV

| id     | original_name | title               | field1_name   | field2_name                    | bool_field_name |  
|--------|---------------|---------------------|---------------|--------------------------------|-----------------|  
| `UUID` | filename1.mp4 | My asset title      | Field 1 Value | Field 2 Value1, Field 2 Value2 | `true`          |  
| `UUID` | filename2.mp4 | Another asset title | Field 1 Value | Field 2 Value1, Field 2 Value2 | `false`         |


The command line arguments required by input mode are listed in Table 2 – Input command-line arguments.

Table 2 – Input command-line arguments

| Flag                       | Required                            | Description                                                |
|----------------------------|-------------------------------------|------------------------------------------------------------|
| `-input <FILE_PATH>`       | no, provided output is used instead | Path to properly formatted CSV file                        |
| `iconik-url <URL>`         | no                                  | iconik URL (default "https://app.iconik.io")               |
| `-metadata-view-id <UUID>` | YES                                 | UUID of metadata view containing fields you want to update |
| `-collection-id <UUID>`    | YES                                 | UUID of collection containing assets you want to update    |
| `app-id <UUID>`            | YES                                 | App ID (provided by iconik)                                |
| `auth-token <JWT>`         | YES                                 | Auth token (provided by iconik)                            |





#### Output Mode

Output mode creates a CSV file conforming to the respective schema constraints, and which contains the metadata field values of the corresponding Metadata View for each Asset in iconik which resides in the provided collection.

##### Example

```bash
$ pd-iconik-io-rd -output ~/Desktop -app-id <AppID> \
-auth-token <AuthToken> -collection-id <CollectionID> -iconik-url \ 
<IconikURL> -metadata-view-id <ViewID>
```

The command line arguments required by output mode are listed in Table 3 – Output command-line arguments.

Table 3 – Output command-line arguments

| Flag                       | Required                           | Description                                                        |
|----------------------------|------------------------------------|--------------------------------------------------------------------|
| `-output <DIR_PATH>`       | no, provided input is used instead | Path to directory where you want to save your CSV                  |
| `iconik-url <URL>`         | no                                 | iconik URL (default "https://app.iconik.io")                       |
| `-metadata-view-id <UUID>` | YES                                | UUID of metadata view containing fields you want to include in CSV |
| `-collection-id <UUID>`    | YES                                | UUID of collection containing assets you want to include in CSV    |
| `app-id <UUID>`            | YES                                | App ID (provided by iconik)                                        |
| `auth-token <JWT>`         | YES                                | Auth token (provided by iconik)                                    |



#### Iconik io Command Reference

##### Name

```pd-iconik-io-rd #iconik metadata CSV reader/writer```

##### Synopsis

```pd-iconik-io-rd [-h][-version]```

```pd-iconik-io-rd [-output <csv-filename>][-app-id <AppID>][-auth-token <AuthToken>][-collection-id <CollectionID>][-iconik-url <IconikURL>][-metadata-view-id <ViewID>]```

##### Options

```bash
-output #toggles the tool to output mode ready to write a CSV file based on the supplied flag values.
-input #toggles the tool to input mode ready to read a CSV file based on the supplied flag values.
-iconik-url #expects a target URL for the iconik instance conforming the https URL schema. Default is https://app.iconik.io.
-app-id #the application key id corresponding to the JWT bearer Token generated in the iconik UI.
-auth-token #the JWT bearer Token generated in the iconik UI.
-collection-id #the ID of the collection in iconik where the assets reside.
-metadata-view-id #the ID of the Metadata View of interest.
```

##### Notes

If neither `input` or `output` mode is selected, the tool will display the version, and then exit.

The tool will also recursively traverse collections in iconik. Therefore, if you provide the top-level collection ID, it will search through all the collections nested within it for assets.




## License
See [License](LICENSE.txt)
