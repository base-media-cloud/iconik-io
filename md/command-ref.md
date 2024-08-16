#### Name

```pd-iconik-io-rd #iconik metadata CSV reader/writer```

#### Synopsis

```pd-iconik-io-rd [-h][-version]```

```pd-iconik-io-rd [-output <csv-filename>][-app-id <AppID>][-auth-token <AuthToken>][-collection-id <CollectionID>][-iconik-url <IconikURL>][-metadata-view-id <ViewID>]```

#### Options

```bash
-output #toggles the tool to output mode ready to write a CSV file based on the supplied flag values.
-input #toggles the tool to input mode ready to read a CSV file based on the supplied flag values.
-iconik-url #expects a target URL for the iconik instance conforming the https URL schema. Default is https://app.iconik.io.
-app-id #the application key id corresponding to the JWT bearer Token generated in the iconik UI.
-auth-token #the JWT bearer Token generated in the iconik UI.
-collection-id #the ID of the collection in iconik where the assets reside.
-metadata-view-id #the ID of the Metadata View of interest.

```

#### Notes

If neither `input` or `output` mode is selected, the tool will display the version, and then exit.
The `size` value is returned in Bytes.