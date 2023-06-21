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

There are two primary modes of use for this script, input and output.  Input mode requires you feed the script a CSV file that is properly formatted.  The format of this CSV has some basic requirements.

	- First row MUST be a header row.  
	- R1C1 MUST be simply `id`
	- R1C2 MUST be simply `title`
	- R1C3 -> R1Cn are the `name` attributes of the metadata fields in the view you want to manipulte
	- First column is ALWAYS the UUID of the asset
	- Second column is ALWAYS the title of the asset
	- Columns 3->n are the values of the metadata fields in R1
<!-- 	- If a field can have multiple values, they must be comma separated in the appropriate cell.
	- If a field is a boolean, it must be either `TRUE` or `FALSE` -->

| id | title | field1_name | field2_name | bool_field_name |
| ------ | ------ | ------ | ------ | ------ |
| `UUID` | My asset title | Field 1 Value | Field 2 Value1, Field 2 Value2 | `TRUE` |
| `UUID` | Another asset title | Field 1 Value | Field 2 Value1, Field 2 Value2 | `FALSE` |


For input mode, there are a few required command line arguments

| Flag | Description |
| ------ | ------ |
|  `-input <FILE_PATH>` | Path to properly formatted CSV file |
|  `-metadata-view-id <UUID>` | UUID of metadata view containing fields you want to update |

For output mode, you are required to use a few more flags

| Flag | Description |
| ------ | ------ |
|  `-output <DIR_PATH>` | Path to directory where you want to save your CSV |
|  `-metadata-view-id <UUID>` | UUID of metadata view containing fields you want to put into CSV |


```bash
-app-id string #iconik Application ID
-auth-token string #iconik Authentication token
-collection-id string #iconik Collection ID
-iconik-url string #iconik URL (default "https://preview.iconik.cloud")
-metadata-view-id string #iconik Metadata View ID
-input string #Input mode - requires path to input CSV file
-output string #Output mode - requires path to save CSV file

# Eg.
#Input mode:
pd-iconik-io-rd -input ~/Desktop/input.csv -app-id <AppID> -auth-token <AuthToken> -collection-id <CollectionID> -iconik-url <IconikURL> -metadata-view-id <ViewID>

#Output mode:
pd-iconik-io-rd -output ~/Desktop -app-id <AppID> -auth-token <AuthToken> -collection-id <CollectionID> -iconik-url <IconikURL> -metadata-view-id <ViewID>
```



## License
See [License](LICENSE.txt)
