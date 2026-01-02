# Annotator
Reading OpenEdge 4GL annotations and outputting them to a JSON file.

There are four types of annotations: 
- class
- method
- property
- free
An annotation is considered `free` when it cannot be tied to one of the other types.
For an example of what to expect from the output, see the example JSON at the bottom.

## Usage

```
annotator parse <directory> [flags]
```

### Flags
- `-o, --output <file>` - Output file (default: annotations.json)
- `--stdout` - Output to stdout instead of file
- `--compact` - Compact JSON output
- `-l, --loglevel <level>` - Log level: none, error, info, debug, trace (default: info)
- `--logtoconsole` - Log to console instead of annotations.log
- `-h, --help` - Display help
- `-v, --version` - Display version

### Example
```
annotator parse c:\myproject -o results.json
```

## Docker
The releases are put in an container image and can be found at:
`docker.io/devbfvio/openedge-annotator`
Two volume mappings are needed:
-  `/app/src` - the root directory which is traversed for .cls files.
-  `/app/output` - a directory which holds the resulting annotations.json

```
docker run -d -v ./4gl:/app/src -v ./result:/app/output annotator:latest 
```

Notes:
- The image is derived `FROM scratch`, meaning that it's just a wrapper for the command.
- logging is done to `stdout`
  
## Release Notes

| Version | Date | Description |
|---------|------|-------------|
| 0.10.0  | 2026-01-02 | added property annotations |
|         |            | added labels to container image |
|         |            | added documentation to README.md |
| 0.9.0   | 2026-01-02 | added docker image major/minor semver's |
| 0.0.4   | 2026-01-02 | added devbfvio/openedge-annotator container image creation |
|         |            | added elapsed time to log |
|         |            | removed version from archive file name |
| 0.0.3   | 2026-01-01 | fixed immutable release issue | 
| 0.0.2   | 2025-12-31 | fixed bug with false annotations within comments |
| 0.0.1   | 2025-12-31 | Initial release with basic annotation parsing for 4GL class files |

## Example JSON

The following is a fragment of `annotation.json` of a parse of the examples classes in the `4gl` directory.
```
{
  "annotations": {
    "auth": [
      {
        "name": "auth",
        "attributes": [
          {
            "name": "required",
            "value": "true"
          }
        ],
        "file": "UserService.cls",
        "classname": "UserService",
        "type": "method",
        "constructName": "GetUser",
        "annotationLine": 8,
        "constructLine": 9
      },
      {
        "name": "auth",
        "attributes": [
          {
            "name": "role",
            "value": "admin"
          }
        ],
        "file": "api/controller/ProductController.cls",
        "classname": "api.controller.ProductController",
        "type": "method",
        "constructName": "DeleteProduct",
        "annotationLine": 28,
        "constructLine": 29
      }
    ]
  }
}
```
