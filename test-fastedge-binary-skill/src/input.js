async function eventHandler(event) { return new Response("Modified\!") } addEventListener("fetch", (event) => { event.respondWith(eventHandler(event)) })
