# otherside

Web application to generate a Spotify playlist based on today's events happening in particular US cities. Built with VueJS, JavaScript/TypeScript, Golang, Redis, Spotify API, and SeatGeek API.

Full website: branch331.github.io/otherside

## Frontend

The frontend is built in VueJS and TypeScript, allowing users to first select cities and genres that they're interested in:

![City selection](/gh-images/selectedcitiesdt.png)

*City selection - desktop*

![City selection](/gh-images/selectedcitiesmb.png)

*City seleciton - mobile*

Each time a city or genre selection is changed, a request is made to the backend for local events for said cities and genres. 

Once the desired cities and genres are selected, the user can Log In using Spotify's authentication...

![Spotify login](/gh-images/spotifylogin.png)

...after which, they will see a rotating preview of the artists for which they can generate a playlist for. Logging in will generate a Spotify client that can be used for subsequent Spotify API requests.

![Built playlist](/gh-images/builtplaylistdt.png)

*Built playlist - desktop*

![Built playlist mb](/gh-images/builtplaylistmb.png)

*Built playlist - mobile*


## Backend

The backend is built in Golang, and uses Redis to cache requested data from the Spotify and SeatGeek APIs. 

Notable HTTP endpoints include:
* _authenticate_ - This endpoint makes a request to the Spotify API for a new authentication URL and returns it to the VueJS frontend. In this case, the Golang backend acts as a middleman for the frontend - the Spotify client will actually reside in the backend, since our goal is to make requests from the backend and store them directly in Redis for caching. The frontend will still have to pass along the auth token in the Header for valid Spotify requests.
* _localevents_ - Breaks down 'cities' and 'genres' query parameter arrays and makes requests to the SeatGeek API for events in the aforementioned cities. These events are then filtered by genre and returned to the requester.
* _artistids_ - Takes an array of local events and returns an array of corresponding Spotify IDs for those performing at the event.
* _toptracks_ - Takes an array of Spotify artist IDs and returns an array of the top track for each of those artists.
* _buildplaylist_ - Simply builds a Spotify playlist based on top tracks.