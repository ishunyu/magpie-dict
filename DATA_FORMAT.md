# Data Format
This is the expected format for how the data file should be organized.

## Folder Structure
Each show should get its own folder, with a manifest file and data folder as its contents.

```
/------
 |- show_a
  |- manifest.json
  |- data
   |- 01.csv
   |- 02.csv
   |- ...
 |- show_b
  |- manifest.json
  |- data
   |- 01.csv
   |- 02.csv
   |- ...
```

### manifest.json
JSON file containing information about the show.
```
{
    "title": "Show A"
}
```

### data
Directory containing subtitle data.

### 01.csv
Subtitle data in the following format.
```
<subtitle_A_timestamp_start>,<subtitle_A_timestamp_end>,<subtitle_A_subtitle>,<subtitle_B_timestamp_start>,<subtitle_B_timestamp_end>,<subtitle_B_subtitle>
```
