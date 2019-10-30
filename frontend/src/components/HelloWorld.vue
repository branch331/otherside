<template>
  <div class="hello">
    <h1>{{ msg }}</h1>

    <button v-on:click="redirectToURL">Log In</button>
    
    <h2>Test Data:</h2>
    <h2 v-if="localEvents">{{localEvents}}</h2>

    <h2>Test Top Tracks Data:</h2>
    <h2 v-if="topTracks">{{topTracks}}</h2>

    <h2>Playlist Status:</h2>
    <h2 v-if="playlistStatus">{{playlistStatus}}</h2>

  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'HelloWorld',
  props: {
    msg: String,
  },
  data () {
    return {
      localEvents: null,
      topTracks: null,
      playlistStatus: null,
      spotifyAuthenticationUrl: null,
      spotifyStateString: null
    }
  },
  mounted () {
    /*
    axios.get("http://localhost:8081/localevents?postcode=78745&miles=50")
      .then((response => {
        this.localEvents = response.data;
        this.getTopTracks(response.data);
      }));
    */
    
    this.setNewSpotifyAuthenticationUrl();
  },
  methods: {
    getTopTracks: function (localEventData) {
      axios.post("http://localhost:8081/toptracks", JSON.stringify(localEventData))
        .then((response => {
          this.topTracks = response.data;
          this.buildPlaylist("Greetings from Axios", "A description for the ages.", response.data);
        }));
    },
    buildPlaylist: function (name, desc, topTracks) {
      var buildPlaylistURL = "http://localhost:8081/buildplaylist?name=" +
        name +
        "&desc=" +
        desc;

      axios.post(buildPlaylistURL, JSON.stringify(topTracks))
        .then((response => {
          this.playlistStatus = response.data;
        }));
    },
    setNewSpotifyAuthenticationUrl: function() {
      this.spotifyStateString = this.getRandomStateString()

      var getAuthenticationRequestUrl = "http://localhost:8081/authenticate?state="
      +  this.spotifyStateString;

      axios.get(getAuthenticationRequestUrl)
        .then(response => {
          this.spotifyAuthenticationUrl = response.data;
        })
    },
    getRandomStateString: function() {
      return Math.random().toString(36).substring(2,15) + Math.random().toString(36).substring(2,15)
    },
    redirectToURL: function() {
      window.location = this.spotifyAuthenticationUrl;
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
