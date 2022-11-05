Permission:

* principal
* resource
* action

# Version 1

## Dush

| REST        | resource             | action      |
|-------------|----------------------|-------------|
| list-albums | album:owner/dush:*   | list-albums |
| list-medias | album:owner/dush:*   | list-medias |
| get-media   | media:owner/dush:*:* | get-media   |

## Vero

| REST        | resource                                                                    | action      |
|-------------|-----------------------------------------------------------------------------|-------------|
| list-albums | album:dush:2022-aout, album:dush:2022-feb                                   | list-albums |
| list-medias | album:dush:2022-aout, album:dush:2022-feb                                   | list-medias |
| get-media   | media:dush:2022-aout:*, media:dush:2022-feb:*, media:dush::mediaId:filename | get-media   |

## Claire

| REST        | resource                                                                    | action      |
|-------------|-----------------------------------------------------------------------------|-------------|
| list-albums | album:owner/claire, album:family/duchatelle-magnier                         | list-albums |
| list-medias | album:dush:2022-aout, album:dush:2022-feb                                   | list-medias |
| get-media   | media:dush:2022-aout:*, media:dush:2022-feb:*, media:dush::mediaId:filename | get-media   |

# Version 2

| REST                  | resource             | action      |
|-----------------------|----------------------|-------------|
| list-albums (album:*) | *:dush:*,            | *           |
| list-medias           | album:owner/dush:*   | list-medias |
| get-media             | media:owner/dush:*:* | get-media   |

# Version 3

```
Medias <- Album <- Family <- user
       ^---- Owner  <|
```

# Version 4 -- associations

| Attached to | Type   | resource id   | role                 | resource name |
|-------------|--------|---------------|----------------------|---------------|
| User        |        |               |                      |               |
| -           | owner  | <owner email> | manage               | -             |
| -           | family | <family name> | manage / contributor | -             |
| -           | media  | <media id>    | view                 | <filename>    |
| -           | api    | viewer / ...  | consumer             |               |

