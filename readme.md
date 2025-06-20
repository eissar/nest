<!--toc:start-->
- [Overview of features](#overview-of-features)
    - [Command line interface](#command-line-interface)
      - [Add eagle item](#add-eagle-item)
      - [List eagle item](#list-eagle-item)
      - [Reveal item id in explorer](#reveal-item-id-in-explorer)
      - [Switch library (synchronously)](#switch-library-synchronously)
    - [Host your eagle images](#host-your-eagle-images)
      - [Render your eagle images in obsidian](#render-your-eagle-images-in-obsidian)
    - [Tray icon](#tray-icon)
    - [extensions to the eagle API](#extensions-to-the-eagle-api)
    - [Configuration](#configuration)
- [TODO:](#todo)
<!--toc:end-->

## Usage
in the command line, run:
```bash
nest start
```
a tray icon should appear in the bottom right.

## Overview of features

### Command line interface

#### Add eagle item
```bash
nest add ./image.png
```

#### List eagle item
```bash
nest list 200
```

```bash
nest list 200 | jq
```

#### Reveal item id in explorer
```bash
nest reveal <itemId>
```


(wip pipe support)
```ps1
nest list | select -First 1 -ExpandProperty id | nest reveal
```

#### Switch library (synchronously)
this will switch libraries synchronously. by default, the eagle api
endpoint will return a status code as soon as possible, without waiting for
the switching process to complete. This is inconvenient for scripting and automation
as you cannot immediately perform next actions or retrieve the current library information
without the risk of retrieving stale data. this endpoint will wait until the current library
is in fact switched and ready for subsequent requests.

```ps1
nest switch inspo && nest list
```

### Host your eagle images
render your images from disk:
using links like `localhost:port/M787F6GA16D3D` we can retrieve image data from eagle.
You could use this to render your eagle images anywhere that takes an image url.
#### Render your eagle images in obsidian
By default, dragging an item from eagle into obsidian creates an inline preview. However, this copies the image into the obsidian vault, which
creates unnecessary data duplication. instead, using this server, you can render obsidian images using a link like:
`![][http://localhost:1323/M787F6GA16D3D]`.

This also allows for more advanced application, like the example video below:
where I am querying the eagle API for items filtered by a tag and rendering the images directly in my
obsidian vault using https://blacksmithgu.github.io/obsidian-dataview/ This could additionally be extended with
html / css.

[![demo video](http://img.youtube.com/vi/UfN2Ad-iLoE/0.jpg)](http://www.youtube.com/watch?v=UfN2Ad-iLoE "Obsidian dataview demo")
js from video:
```js
let data = await fetch("http://localhost:41595/api/item/list?orderBy=CREATEDATE&limit=10&tags=eagle-demo")
	.then(r => r.json())
let imgs = data.data.map(d => d.id)
let innerHtml = imgs.map(id => {
	return `<img src="http://localhost:1323/${id}"></img>`
})

dv.paragraph(`<div style ="height:300px;">${innerHtml}</div>`)
```

accepted links for this feature are as follows:
1. `/<eagleItemId>`
2. (WIP) `/http://localhost:41595/item?<eagleItemId>`
where eagleItemId is the eagle item id.

to get an eagle item id, you can use the context menu option `Copy Link`
in eagle (i.e, http://localhost:41595/item?id=M787F6GA16D3D).

### Tray icon
- [X] quit the program
- [ ] open config files from tray icon
manage the server (close)
Planned:
update configuration

### extensions to the eagle API
- [X] autogen docs
    - [ ] document each endpoint
- [X] get library filepath
- [ ] improved error messages
- [ ] Search
    - [ ] glob fts search
    - [ ] Filter items by Tag Count
- [X] api/item/reveal
- [X] synchronous api/library/switch
- [ ] Config
    - [X] open config files from tray icon

### Configuration
config is set in %USERPROFILE%/.config/nest directory.
most of the things in there are not used right now, but you can change the port to something else.
Just make sure it isn't being used by any other service on your computer.

# TODO:
- [X] generalized wrapper fallback for //api...
- [X] create watcher for mtime.json
    - [X] basic impl. nest.WatchMtime

- [ ] method Query across libraries
- [ ] dynamic key in config/libraries.json
- [ ] use exe as server launcher (with cli flag) and also interface for CLI interaction with eagle.
- [X] On starting the server, libraries are loaded from libraries.json.
- [ ] create events channel for dynamic library reading from eagle.

* Command Line Interface
- [X] Start server -start
- [X] stop server -stop
* Subcommands
- [X] nest add
- [X] nest list
- [X] nest reveal
- [X] nest info

- [ ] query eagle via tags
- [ ] eaglepack generator (bulk url import solution?)
    opt: include date as tag?

- [ ] embed c# object message converter/ convertFrom-json in output?

https://github.com/eissar/nest/browser-query -> has to define event loop. this should be opt-in.
how to opt in via: config flag? cli flag?

> "I've learned a few things from the contributions of Wox, one in particular is not to put resources in a non-Github location, as they will become invalid and outdated over time as they migrate. All future Wox resources will be hosted on Github, including the plugin store, theme store, documentation, discussions and more!" <https://github.com/Wox-launcher/Wox/discussions/3937>


> Disclaimer: <br>
> at the moment, this is experimental software. I will be making breaking changes and changing apis.

