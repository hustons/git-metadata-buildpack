# GIT Metadata Buildpack

[![CircleCI](https://img.shields.io/circleci/project/github/bstick12/git-metadata-buildpack.svg)](https://circleci.com/gh/bstick12/git-metadata-buildpack) 
[![Download](https://api.bintray.com/packages/bstick12/buildpacks/git-metadata-buildpack/images/download.svg?version=0.3.0) ](https://bintray.com/bstick12/buildpacks/git-metadata-buildpack/0.3.0/link)
[![codecov](https://codecov.io/gh/bstick12/git-metadata-buildpack/branch/master/graph/badge.svg)](https://codecov.io/gh/bstick12/git-metadata-buildpack)

This is a [Cloud Native Buildpack V3](https://buildpacks.io/) that adds GIT metadata as a layer to the built container.

This buildpack is designed to work in collaboration with other buildpacks.

## Usage

```
pack build <image-name> --builder cloudfoundry/cnb:cflinuxfs3 --buildpack https://bintray.com/bstick12/buildpacks/download_file?file_path=git-metadata-buildpack-0.3.0.tgz
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
