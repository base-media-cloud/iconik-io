There are two modes of use for this command-line tool: `input`, and `output`.



#### Input Mode

Input mode takes a pre-populated CSV file conforming to the respective schema constraints and writes the contained metadata field values contained therein to the corresponding Metadata View for each Asset in iconik based on the supplied Collection ID.

###### Example

```bash
$ pd-iconik-io-rd -input input.csv -app-id <AppID> \
-auth-token <AuthToken> -collection-id <CollectionID> -iconik-url \ 
<IconikURL> -metadata-view-id <ViewID>
```

###### Schema constraints

- First row MUST be a header row.
- R1C1 MUST be `id`.
- R1C2 MUST be `original_name`.
- R1C3 MUST be `size`.
- R1C4 MUST be `title`.
- R1C5 -> R1Cn are the name attributes of the metadata fields in the view you want to manipulate.
- First column MUST always be the UUID of the asset.
- Second column MUST always be the original filename of the asset.
- Third column can include the filesize of the asset (in bytes), but if not including filesize MUST be left blank.
- Fourth column MUST always be the title of the asset.
- Columns 5->n are the values of the metadata fields in R1.
- If a field can have multiple values (e.g., Tags), they must be comma separated in the appropriate cell.
- If a field is boolean, it must be either true or false.


###### Example CSV

| id     | original_name | size   | title               | field1_name   | field2_name                    | bool_field_name |  
|--------|---------------|--------|---------------------|---------------|--------------------------------|-----------------|  
| `UUID` | filename1.mp4 | 176985 | My asset title      | Field 1 Value | Field 2 Value1, Field 2 Value2 | `true`          |  
| `UUID` | filename2.mp4 | 176985 | Another asset title | Field 1 Value | Field 2 Value1, Field 2 Value2 | `false`         |


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

###### Example

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
