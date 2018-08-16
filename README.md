# go-mvn-dl
go-mvn-dl is Go library that allows you to easily download maven artifacts from Go code. Supports custom maven repos and authentication.

# Install
```
go get github.com/navdeepsekhon/go-mvn-dl
```

# Import
```
import "github.com/navdeepsekhon/go-mvn-dl/download
```

# Download Function
`download.Download()` takes the following six arguments. First one is required rest can be empty strings.

 | **Argument** | **Type** | **Description**                                                                                         |
|--------------|--------|------------------------------------------------------------------------------------------------------------|
| artifactName | string | Maven artifact name, eg: `org.apache.commons:commons-lang3:3.4`                                            |
| destinaionDir| string | Directory where to write the downloaded file, if blank - defaults to directory where the code is running   |
| repositoryUrl| string | Maven repository url, if blank - defaults to `https://repo1.maven.org/maven2/`                             |
| fileName     | string | Filename to be used for the downloaded file. If blank - derived from artifact name eg commons-lang3-3.4.jar|
| username     | string | If authentication is required for repository, provide the username and password                            |
| password     | string | If authentication is required for repostiry, provide username and password                                 |

# Example
Following call will download `commons-lang3-3.4.jar` as `download.jar` from a custom maven repo `http://internal.repo.com/`

```
download.Download("org.apache.commons:commons-lang3:3.4", ".", "http://internal.repo.com/", "download.jar", "user", "password")
```

# Other Functions

The following helper functions are also available. See the listed test files for example useage

| **Function** | **Examples in** | **Description** |
|--------------|-----------------|-----------------|
| ParseName  | `test/parseName_test.go` | Converts the maven artifact name string to `Artifact` |
| ArtifactUrl  | `test/artifactUrl_test.go` | Returns the url to download the artifact from. |
| ArtifactUrl  | `/test/filename_test.go` | Generates a filename from artifact |

# Contributions

Please feel free to send pull requests.

# Credits

Inspired by npm based [mvn-dl](https://github.com/laat/mvn-dl/graphs/contributors)
