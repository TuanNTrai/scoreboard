# Score tracker for discgolf and golf

- En käytä vuex koska toistaiseksi se tuntuu lisäävän liikaa

## Things to consider
- PWA (Offline support)
- Great UX
- HTML Router
- Random avatars
- User cookies to set id and createdAt
- Measure time spent (holes also)
- GPS Nearest course in Helsinki
- After action report with linegraph if total is saved for each basket
- Take care if course aint played in right order
- Mutex
- Secure cookie
- If cookie is found but no game, fix that injecting starttime or something
- Updating stats require hold button
- Sound effects
- Copy id to clipboard
- Multiple languages
- Delete based on updated at
- User specific friends

## UI
- Star icon to leader
- Round plus and inc buttons
- Avatars

## Create new course

### POST
```json
{
	"basketCount": 2,
	"players": ["Tiger King", "Ying Yang"]
}
```

### Server responses
```json
{
  "id": "2ty",
  "basketCount": 2,
  "active": 1,
  "Baskets": {
    "1": {
      "OrderNum": 1,
      "Par": 0,
      "Scores": {
        "Tiger King": {
          "Score": 0,
          "OB": 0
        },
        "Ying Yang": {
          "Score": 0,
          "OB": 0
        }
      }
    },
    "2": {
      "OrderNum": 2,
      "Par": 0,
      "Scores": {
        "Tiger King": {
          "Score": 0,
          "OB": 0
        },
        "Ying Yang": {
          "Score": 0,
          "OB": 0
        }
      }
    }
  }
}
```
