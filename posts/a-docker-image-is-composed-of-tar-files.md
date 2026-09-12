---
title: A Docker image is composed of .tar files
date: 2026-09-12
---

A Docker image is a `.tar` file that contains metadata and its respective layers. Each layer is also a `.tar` file containing the file system diff

## Example

We pull the `python` image, save it to a `.tar` file, and extract that file

```bash
docker pull python
docker save python -o python.tar
mkdir python-extracted && tar -xf python.tar -C python-extracted
cd python-extracted
```

The end result is a set of blobs. Some of these blobs are configuration files, while others are layers. The `manifest.json` file tells us which blobs are layers

```
.
├── blobs
│   └── sha256
│       ├── 375591c23c...
│       ├── 56a9d4a189...
│       ├── 92385f4b6d...
│       └── bb56896431...
│       └── ...
├── index.json
├── manifest.json
├── oci-layout
└── repositories
```

To extract all layer blobs and view the file system diffs, run the script below (inside `python-extracted`)

```bash
while IFS= read -r f; do
  if file --mime-type "$f" | grep -q "application/x-tar"; then
    dir=$(dirname "$f")
    base=$(basename "$f")
    name="${base%.tar}"
    tmp="$dir/.tmp_extract_$base"
    mv "$f" "$tmp"
    mkdir -p "$dir/$name"
    tar xf "$tmp" -C "$dir/$name" && rm "$tmp"
  fi
done < <(find ./blobs/sha256 -maxdepth 1 -type f)
```
