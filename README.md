# GIT Metadata Buildpack

This is a [Cloud Native Buildpack V3](https://buildpacks.io/) that adds GIT metadata as a layer to the built container.

This buildpack is designed to work in collaboration with other buildpacks.

## Usage

```
pack build <image-name> --builder cloudfoundry/cnb:cflinuxfs3 --buildpack /path/to/git-metadata-buildpack 
```

The following layer will be added to your container

```
/layers/io.bstick12.buildpacks.git-metadata/git-metadata/
```

The `git-metadata.toml` file will contain the following elements

```
sha = "<sha>"
branch = "<remote>/<branch>"
remote = "<remote url>"
```

## Development

`scripts/unit.sh` - Runs unit tests for the buildpack
`scripts/build.sh` - Builds the buildpack
`scripts/package.sh` - Package the buildpack
