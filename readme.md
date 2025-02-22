<!--toc:start-->
- [Overview of features](#overview-of-features)
    - [Host your eagle images](#host-your-eagle-images)
    - [Tray icon](#tray-icon)
    - [(WIP) Wrapper around the eagle api](#wip-wrapper-around-the-eagle-api)
    - [(PLANNED) Re-Implementation/ Extensions to the eagle API](#planned-re-implementation-extensions-to-the-eagle-api)
- [TODO:](#todo)
    - [Configuration](#configuration)
<!--toc:end-->

# Overview of features

### Host your eagle images
render your images from disk:
using links like `localhost:port/M787F6GA16D3D` we can retrieve image data from eagle.
You could use this to render your eagle images anywhere that takes an image url.
#### Host your eagle images in obsidian
By default, dragging an item from eagle into obsidian creates an inline preview. However, this copies the image into the obsidian vault, which
creates unnecessary data duplication. instead, using this server, you can render obsidian images using a link like:
`![][http://localhost:1323/M787F6GA16D3D]`. 

This also allows for more advanced application, like the example video below:
where I am querying the eagle API for items filtered by a tag and rendering the images directly in my
obsidian vault using https://blacksmithgu.github.io/obsidian-dataview/ This could additionally be extended with more advanced
html / css.

[![demo video](http://img.youtube.com/vi/UfN2Ad-iLoE/0.jpg)](http://www.youtube.com/watch?v=UfN2Ad-iLoE "Obsidian dataview demo")

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

### (WIP) Wrapper around the eagle api
You can call endpoints in the eagle api using the same syntax as the default eagle api.
/api/item/list?...

### (PLANNED) Re-Implementation/ Extensions to the eagle API
- [ ] autogen docs
- [ ] improved error messages
- [ ] Search
    - [ ] glob fts search
    - [ ] Filter items by Tag Count 
- [ ] api/item/reveal <br>
- [ ] synchronous api/library/switch <br>

- [ ] open config files from tray icon <br>

# TODO:
- [X] generalized wrapper fallback for //api... <br>
- [ ] method Query across libraries <br>
- [ ] dynamic key in config/libraries.json <br>
- [ ] use exe as server launcher (with cli flag) and also interface for CLI interaction with eagle.

### Configuration
(ignore this) (move to wiki) On starting the server, libraries are loaded from libraries.json.
- [ ] dynamic library reading from eagle. <br>


> Disclaimer: <br>
> at the moment, this is experimental software. I will be making breaking changes and changing apis.

