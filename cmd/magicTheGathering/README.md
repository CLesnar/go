# WeatherDate Service
This service utilizes Magic: The Gathering API to create a web based game. 

## Magic: The Gathering API Details
The MTG APIs used relate to card data and set data. The official website for MTG API can be found here: https://magicthegathering.io/#. The official docs found in linked from above can be found here: https://docs.magicthegathering.io/. And the go SDK can be found here: https://github.com/MagicTheGathering/mtg-sdk-go. 

The go client is used:
    `go get github.com/MagicTheGathering/mtg-sdk-go`

## Example of magicTheGathering (this project) API Call
** TODO **
http://localhost:8111/v1/magic

## Magic: The Gathering How to Play
https://magic.wizards.com/en/how-to-play#ltp-video

## Run weatherWebApp
Checkout the git repo clesnar/go and rely on the .vscode config - run debug. Or build and run executable.
The default `host:port` is `localhost:8111`.
Once weatherWebApp is running, then in a browser, navigate to http://localhost:8111/v1/magic.

## debug build
go build -buildvcs=false -o /workspaces/clesnar-go/cmd/weatherWebApp/__debug_bin1800304272 -gcflags all=-N

## Flowcharts
Mermaid Diagrams; for details see https://mermaid.js.org/syntax/examples.html#a-commit-flow-diagram or https://mermaid.live/edit.
### Building the Deck
```mermaid
flowchart TD
    A[Service] --> |Request|M[MTG Public API]
    M --> |GET|C[Cards]
    M --> |GET|S[Sets]
    S --> |Validate|P(Card Pool)
    C --> |Validate|P
    UC[Client] <--> |GET|P
    UC --> |POST|PD[Player Deck]
```
### Main Thread Operations Sequence Diagram
```mermaid
sequenceDiagram
    participant db as Storage
    participant gm as Game Manager
    participant p1 as Player 1
    participant p2 as Player 2

    Note over db,p2: Setup Player Decks
    Note over db,p2: Game Manager allows asynchronous player input in building Player Decks

    p1 ->> gm: Build Deck
    gm ->> db: validate & store p1 update
    gm ->> p1: Game State
    p2 ->> gm: Build Deck
    gm ->> db: validate & store p1 update
    gm ->> p2: Game State

    Note over db,p2: When Player Decks are valid Game can commence
    Note over db,p2: Game Manager waits for player input in turn based game

    gm ->> p1: Game State
    gm ->> p2: Game State

    p1 ->> gm: action
    gm ->> db: is player's turn, validate action
    gm ->> p1: Game State 
    gm ->> p2: Game State 
    
    p2 ->> gm: action
    gm ->> db: is player's turn, validate action
    gm ->> p2: Game State 
    gm ->> p1: Game State 
    
    Note over db,p2: Iterate until Game is complete

```