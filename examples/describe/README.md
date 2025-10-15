# describe

Describes an image using a Vision Language Model. Works with either remote or local files.

## Install

```
go install ./examples/describe
```

## Running

```shell
describe https://www.publicdomainpictures.net/pictures/220000/t2/llama-1496506141lWt.jpg
```

Set any flags before the image you want to describe like this:

```shell
describe -models="/home/ron/models/" -v https://www.publicdomainpictures.net/pictures/220000/t2/llama-1496506141lWt.jpg
```

