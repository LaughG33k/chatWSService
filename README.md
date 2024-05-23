ChatWSService

The service opens websocket connections and starts communication with clients.
It accpets incoming messages on save, delete, edit and sends these to messageDBService.
It's also redirects message from sender to recipient by redis's pub/sub