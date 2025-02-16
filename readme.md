/*
	Support following color codes:
	Foreground
	Black
	Bright black
	Red
	Bright red
	Green
	Bright green
	Yellow
	Bright yellow
	Blue
	Bright blue
	Purple
	Bright purple
	Cyan
	Bright cyan
	White
	Bright white
*/


# Features
[ ] - Filter by Tag Count [1]
[ ] - Query across libraries
[ ] - recent eagle items

1. (make-a-wish-readonly) <https://discord.com/channels/1169553626860101672/1178804692655034539/1178804692655034539>






/* __MACRO__
```lua
!wt.exe -d "$env:CLOUD_DIR\Code\go\web-dashboard" pwsh -c ./build.ps1
```
:so ./build.lua

Explanation:
open new terminal instance with the server running.
*/

/*
// TODO:
			[X] - move types to types.go
			[X] - add HTMX
			[X] - add middleware to serve url without file suffix.
			[X] - Recent notes
			[X] - move time.now calls to middleware (custom)
			q5s: enumerate-Windows.ps1
			[ ] - Try reflection for template functions?
			[ ] - add action parameter to recent notes
			[X] - Create build.lua
			[X] - add sse listener to eagle-plugin
				[X] - for eagle.tabs.query({})
				on BroadcastSSETargeted({event:"getTabs"...})
					-> eagle.tabs.query({})
					-> post("api/uploadTabs")
					-> fmt.PrintLn tabs
*/

[ ] - recent eagle items
[X] - browser tabs
    [X] - filtering
    [ ] - 'Raw' flag or special char at start of query to indicate not to automatically double up on asterisks.

[ ] - Integrate into a GUI search engine.
[ ] - Powertoys run integration (WIP)
[ ] - Everything integration
[ ] - custom search engine in go extensible with Lua?
[ ] - ? clone <https://github.com/srwi/EverythingToolbar>
[ ] - tauri (entirely custom) how much overhead? need immediate


[ ] - frecency?

[ ] - event based architecture?
    [ ] - request id -> resolved via channels, impl. timeout
[ ] - websockets, or SSE?

[ ] - Trim pending requests V2 (eagleModule)
[ ] - rename eagleModule to browserModule

[X] - does htmx support SSE? A: YES:
<https://htmx.org/extensions/sse/>
<https://htmx.org/extensions/ws/>

## security considerations
[ ] - Possible require api key for only non-local connections?
[ ] - How to tell local vs not?


eagle item link format: eagle://item/LVN5ZJY1XNVSX
this allows you to use eagle as an image server.

this works:
<img src="http://localhost:3000/eagleApp/images/LVN5ZJY1XNVSX" alt="GFG image">
<img src="http://localhost:3000/item/M43KSCG4AHK3T" width = '80%'>
<img src="http://localhost:3000/item/LVN5ZJY1XNVSX">

also just leave a link like this:
<eagle://item/LVN5ZJY1XNVSX>
(clicking on this opens the item in eagle.)

ideas :
server naming:
<http://eagleServ/item/LVN5ZJY1XNVSX>


Userscript Logic (Youtube Music)
<https://gist.github.com/eissar/ce251e8d49afc20888e4d6398d1ee7bd>
["X:/Dropbox/Code/javascript/Userscripts/yt-music-broadcast.userscript.js"]
```js
    const targetNodes = document.querySelectorAll("yt-formatted-string.ytmusic-player-bar:not(.complex-string)");

    function postMessage(msg) {
        if (!msg) {
            return
        };
        if (typeof msg !== "string"|| msg.length === 0){
            return
        };
        try {
            fetch(`http://localhost:1323/api/broadcast/yt-music?elem=<a id="message-container">${msg}</a>`)
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`HTTP error! status: ${response.status}`); // Throw error for non-2xx responses
                    }
                    return response.json(); // Or response.text() if the response is not JSON
                })
                .then(data => {
                    // Process the successful response data
                    console.log("Broadcast successful:", data);
                })
                .catch(error => {
                    // Handle errors during the fetch or processing
                    console.error("Error during broadcast:", error);
                });
        } catch (e) {
            console.error("An unexpected error has occurred:", e); //Catch any unexpected errors.
        }
    }

    const handleTextChange = (mutationsList, observer) => {
      for (const mutation of mutationsList) {
        if (mutation.type === 'characterData' || (mutation.type === 'childList' && mutation.addedNodes.length > 0 && mutation.addedNodes[0].nodeType === 3)) { // Check for text changes
          //console.log('changed', mutation.target.textContent);
          postMessage(mutation.target.textContent)
        }

        if (mutation.type === 'childList' && mutation.addedNodes.length > 0) {
            mutation.addedNodes.forEach(node => {
                if (node.nodeType === 1 && node.matches("yt-formatted-string.ytmusic-player-bar:not(.complex-string)")) {
                    observer.observe(node, config);
                }
            })
        }
      }
    };
    // Create an observer instance linked to the callback function.
    const observer = new MutationObserver(handleTextChange);

    // observer cfg
    const config = { characterData: true, childList: true, subtree: true }; // Observe characterData changes and childList changes
    // Start observing each target node.
    targetNodes.forEach(targetNode => {
      observer.observe(targetNode, config);
    });
```
`http://localhost:1323/api/broadcast/yt-music?elem=<a id="message-container">test</a>`

```html
<div id="messages" hx-ext="ws" ws-connect="/ws">
    <div id="message-container"></div>
</div>
```




```js
const evtSource = new EventSource("/sse"); // Replace with your SSE endpoint URL

evtSource.onmessage = (event) => {
  console.log("Received event:", event.data);
  // Process the received data (e.g., update the UI)
  try {
      const data = JSON.parse(event.data); // Try to parse as JSON
        if (data.action && data.action === "getSong") {
            console.log("YESS!")
        }
      // ... handle JSON data
      console.log("Parsed JSON:", data)

  } catch (e) {
      // Handle plain text data or JSON parsing errors
      console.log("Plain text:", event.data)
  }
};
evtSource.onerror = (error) => {
  console.error("EventSource failed:", error);
  // Handle errors (e.g., reconnect, display an error message)
  if (error.target.readyState === EventSource.CLOSED) {
      console.log("Connection closed. Reconnecting...");
      setTimeout(() => {
          const newSource = new EventSource("//api.example.com/sse-demo.php"); // Replace with your SSE endpoint URL
          // copy over the event handlers
          newSource.onmessage = evtSource.onmessage;
          newSource.onerror = evtSource.onerror;
          evtSource = newSource;
      }, 5000); // Reconnect after 5 seconds
  }
};

// Optional: Close the connection when you're done
// evtSource.close();
```



//const Port = "41595"
//const Host = "127.0.0.1" // prefer ip over localhost

// var EagleConfig map[string]string{
// 	"addItemFromPath": "/api/item/addFromPath",
// }

//Test["baseUrl"] = fmt

// func LoadConfig() {
// 	config.New()
// }
/*
	[X] - config
		 [X] - Api key
	[X] - read config
	[X] - check if server open
	[ ] - GetEagleThumbnailFromID (eagle api wrapper)
		[ ] - make custom? (really not neccessary)
*/


	/*
		server.GET("/template/notes-struct", handlers.DynamicTemplateHandler("notes-struct.html", apiroutes.PopulateGetNotesDetail))
		server.GET("/template/windows", handlers.DynamicTemplateHandler("windows.html", apiroutes.PopulateEnumerateWindows))
		server.GET("/template/recent-notes", handlers.PwshTemplateHandler("recent-notes.html", pwsh.PwshScript, "./powershell-utils/recentNotes.ps1"))
		server.GET("/template/key-value", handlers.PwshTemplateHandler("key-value.templ", pwsh.PwshScript, "./powershell-utils/mock_nvim.ps1"))
		server.GET("/template/open-tabs-count", handlers.PwshTemplateHandler("open-tabs-count.templ", pwsh.PwshScript, "./powershell-utils/waterfoxTabs.ps1"))
		//server.GET("/template/open-tabs", dynamicTemplateHandler("open-tabs.templ", apiroutes.PopulateOpenTabs))
		server.GET("/template/recent-eagle-items", handlers.PwshTemplateHandler("recent-eagle-items.templ", pwsh.PwshScript, "./powershell-utils/recentEagleItems.ps1"))
		server.GET("/template/sse-browser-tabs", handlers.StaticTemplateHandler("sse-browser-tabs.templ"))
		server.GET("/template/browser-tabs", handlers.StaticTemplateHandler("browser-tabs.templ"))
		server.GET("/template/recent-notes_layout", handlers.StaticTemplateHandler("recent-notes.layout.html"))
		server.GET("/template/timeline_layout", handlers.StaticTemplateHandler("timeline.layout.html"))
		server.GET("/template/now-playing", handlers.StaticTemplateHandler("ws-now-playing.ytm.templ")) // ./templates/ws-now-playing.ytm.templ
		server.GET("/template/ping", handlers.StaticTemplateHandler("ping.templ"))
	*/
	/*
		g.GET("/broadcast/getTabs", func(c echo.Context) error {
			req := PendingRequests.New()
			bc := fmt.Sprintf(`{"id":"%s","command":"getTabs"}`, req.Id)
			websocketCfg.Broadcast(bc)
			// add Pending WSRequest with callback?
			return c.String(200, "OK")
		})
	*/
	/*
		activates a browser tab if it exists, creates a new tab if it does not.
		cache is created from browser history
		api.BrowserTabActivateOrOpen
		api.GetBrowserHistory
		about:profiles
		["X:\Dropbox\Code\Projects\render-image-blazingly\db"]
	*/

	//server.GET("/api/broadcast/sse", broadcastHandler("getSong"))


	// Module philosophy:
	// UNDER NO CIRCUMSTANCES
	// should html or css be tightly coupled with
	// or packaged in a module (e.g., eagle_module) TODO:

	// ??? access routes in a module like:
	// server.GET("/eagleApp/*", eaglemodule.HandleModuleRoutes)
	// OR



// populate funcs
// set default params with pathparams
// they are empty strings not nil null case.
/*
Every PopulateFunction returns data that will be consumed by a template.
using the context, we can extract parameters or default arguments we can pass to the API calls.

if there is an error, the .error member is populated. This is checked first in the template and if it exists, the template is populated in the error case.

some of the populate functions bubbled errors by just returning c.string(400,err) which is less flexible.
*/

