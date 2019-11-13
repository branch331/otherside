<template>
  <div class="home">

    <h2>HOME</h2>

    <h2>Events Data:</h2>
    <p>THIS component's spotify state: {{ spotifyStateString }}</p>

    <ul v-if="cities">
      <li v-for="city in cities" v-bind:key="city">
        {{city}}
      </li>
    </ul>

    <ul v-if="genres">
      <li v-for="genre in genres" v-bind:key="genre">
        {{genre}}
      </li>
    </ul>

    <div v-if="localEvents">

      <button v-on:click="redirectToURL">Log In</button>

      <p>{{localEvents}}</p>

    </div>

    <h2>Playlist Status:</h2>
    <h2 v-if="playlistStatus">{{playlistStatus}}</h2>

  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'home',
  data () {
    return {
      localEvents: null,
      playlistStatus: null,
      spotifyAuthenticationUrl: null,
      spotifyStateString: null,
      cities: null,
      genres: null
    }
  },
  mounted () {
    this.getAvailableCities();
    this.getAvailableGenres();
    this.getLocalEvents('Austin TX,Washington DC,Nashville TN', '[rock,electronic,hip-hop]');
    this.setNewSpotifyAuthenticationUrl();
  },
  methods: {
    getAvailableCities: function() {
      var citiesURL = "http://localhost:8081/cities";

      axios.get(citiesURL)
        .then((response => {
          this.cities = response.data;
        }));
    },
    getAvailableGenres: function() {
      var genresURL = "http://localhost:8081/genres";

      axios.get(genresURL)
        .then((response => {
          this.genres = response.data;
        }));
    },
    getLocalEvents: function (cities, genres) {
      var localEventsURL = "http://localhost:8081/localevents?cities=" +
      cities +
      "&genres=" +
      genres;

      axios.get(localEventsURL)
        .then((response => {
          this.localEvents = response.data;
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
