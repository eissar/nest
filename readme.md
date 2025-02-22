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
using links like `localhost:0000/M787F6GA16D3D` we can retrieve image data from eagle.
accepted links for this feature are as follows:
1. `/<eagleItemId>`
2. `/http://localhost:41595/item?<eagleItemId>`
where eagleItemId is the eagle item id. you can get this from copying the eagle link
(TODO: create link)  <how to get eagle link from eagle>

the second supported link scheme is the output you get from the context menu option `Copy Link`
in eagle (e.g, http://localhost:41595/item?id=M787F6GA16D3D).

### Tray icon
 - [X] quit the program <br>
 - [ ] open config files from tray icon <br>
manage the server (close)
Planned:
update configuration

### (WIP) Wrapper around the eagle api
You can call endpoints in the eagle api using the same syntax as the default eagle api.
/api/item/list?...

### (PLANNED) Re-Implementation/ Extensions to the eagle API
 - [ ] autogen docs <br>
 - [ ] improved error messages <br>
 - [ ] Search <br>
     - [ ] glob fts search <br>
     - [ ] Filter items by Tag Count <br>
 - [ ] api/item/reveal <br>
 - [ ] synchronous api/library/switch <br>

 - [ ] open config files from tray icon <br>
# TODO:
 - [X] generalized wrapper fallback for //api... <br>
 - [ ] method Query across libraries <br>
 - [ ] dynamic key in config/libraries.json <br>
 - [ ] use exe as server launcher (with cli flag) and also interface for CLI interaction with eagle.


### Configuration
(move to wiki) On starting the server, libraries are loaded from libraries.json.
 - [ ] dynamic library reading from eagle. <br>

> Disclaimer: <br>
> at the moment, this is experimental software. I will be making breaking changes and changing apis.

