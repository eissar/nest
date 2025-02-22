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
1. /<eagleItemId>
2. /http://localhost:41595/item?<eagleItemId>
where eagleItemId is the eagle item id. you can get this from copying the eagle link
(TODO: create link)  <how to get eagle link from eagle>

the second supported link scheme is the output you get from the context menu option `Copy Link`
in eagle (e.g, http://localhost:41595/item?id=M787F6GA16D3D).

### Tray icon
[X] - quit the program
[ ] - open config files from tray icon
manage the server (close)
Planned:
update configuration

### (WIP) Wrapper around the eagle api
You can call endpoints in the eagle api using the same syntax as the default eagle api.
/api/item/list?...

### (PLANNED) Re-Implementation/ Extensions to the eagle API
[ ] - autogen docs
[ ] - improved error messages
[ ] - Search
    [ ] - glob fts search
    [ ] - Filter items by Tag Count
[ ] - api/item/reveal
[ ] - synchronous api/library/switch

[ ] - open config files from tray icon
# TODO:
[X] - generalized wrapper fallback for //api...
[ ] - method Query across libraries
[ ] - dynamic key in config/libraries.json


### Configuration
(move to wiki) On starting the server, libraries are loaded from libraries.json.
[ ] - dynamic library reading from eagle.

> Disclaimer: <br>
> at the moment, this is experimental software. I will be making breaking changes and changing apis.

