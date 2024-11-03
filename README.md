# ncbuild

Declarative build system based on containers.

## Usage

Write a recipe in `ncbuild.json`.

```json
{
  "name": "hello",
  "executable": "/bin/sh",
  "args": ["-c", "mkdir $out && echo hello world > $out/hello.txt"]
}
```

- `name`: name of the recipe
- `executable`: name or path to the executable to run the build
- `args`: arguments to pass to the executable

Then run in the same folder:

```bash
ncbuild build
```
