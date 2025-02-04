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


[ ] - recent eagle items
[ ] - recent




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
