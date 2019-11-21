<template>
  <div class="callback">

    <h2>CALLBACK</h2>

    <h5>Selected Cities</h5>
    <ul v-if="selectedCities">
      <li v-for="city in selectedCities" v-bind:key="city" v-on:click="transferArrayValue(selectedCities, availableCities, city)">
        {{city}}
      </li>
    </ul>

    <h5>Available Cities</h5>
    <ul v-if="availableCities">
      <li v-for="city in availableCities" v-bind:key="city" v-on:click="transferArrayValue(availableCities, selectedCities, city)">
        {{city}}
      </li>
    </ul>

    <h5>Selected Genres</h5>
    <ul v-if="selectedGenres">
      <li v-for="genre in selectedGenres" v-bind:key="genre" v-on:click="transferArrayValue(selectedGenres, availableGenres, genre)">
        {{genre}}
      </li>
    </ul>

    <h5>Available Genres</h5>
    <ul v-if="availableGenres">
      <li v-for="genre in availableGenres" v-bind:key="genre" v-on:click="transferArrayValue(availableGenres, selectedGenres, genre)">
        {{genre}}
      </li>
    </ul>

    <h2>Events Data:</h2>
    <p>THIS component's spotify state: {{ spotifyStateString }}</p>

    <div v-if="localEvents">
      {{localEvents}}
    </div>

    <div v-if="topTracks">
      {{topTracks.length}}
    </div>

    <div v-if="topTracks">
      <button v-on:click="buildPlaylist('Spooky Title', 'Spooooky Description!')">Build Playlist</button>
    </div>

    <h2>Playlist Status:</h2>
    <h2 v-if="playlistStatus">{{playlistStatus}}</h2>

  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'callback',
  data () {
    return {
      localEvents: null,
      playlistStatus: null,
      spotifyStateString: null,
      selectedCities: null,
      availableCities: null,
      selectedGenres: null,
      availableGenres: null,
      topTracks: null
    }
  },
  mounted () {
    this.getAvailableCities();
    this.getAvailableGenres();
    //this.getLocalEvents('[Austin TX,Washington DC,Nashville TN]', '[rock,electronic,hip-hop]');
  },
  methods: {
    transferArrayValue: function(sourceArray, targetArray, value) {
      var index = sourceArray.indexOf(value);
      if (index > -1) {
        if (!(targetArray == this.selectedCities && targetArray.length == 6) && 
        !(targetArray == this.selectedGenres && targetArray.length == 10)) {
          sourceArray.splice(index, 1);
          targetArray.push(value);
        }
      }

      if (this.selectedCities.length != 0 && this.selectedGenres.length != 0) {
        this.getLocalEvents(this.selectedCities, this.selectedGenres);
      } else {
        this.localEvents = null;
      }
    },
    getAvailableCities: function() {
      var citiesURL = "http://localhost:8081/cities";

      axios.get(citiesURL)
        .then((response => {
          this.availableCities = response.data;
          this.selectedCities = new Array();
        }));
    },
    getAvailableGenres: function() {
      var genresURL = "http://localhost:8081/genres";

      axios.get(genresURL)
        .then((response => {
          this.availableGenres = response.data;
          this.selectedGenres = new Array();
        }));
    },
    getLocalEvents: function (cities, genres) {
      var cityString = this.arrayToQueryString(cities);
      var genreString = this.arrayToQueryString(genres);

      var localEventsURL = "http://localhost:8081/localevents?cities=" +
      cityString +
      "&genres=" +
      genreString;

      axios.get(localEventsURL)
        .then((response => {
          this.localEvents = response.data;
          this.getTopTracks(this.localEvents);
        }));
    },
    arrayToQueryString: function (array) {
      var queryString = '[';

      for (var i = 0; i < array.length; i++) {
        queryString = queryString.concat(array[i] + ',');
      }

      queryString = queryString.slice(0, -1).concat(']');
      return queryString;
    },
    getTopTracks: function(events) {
      var topTracksURL = "http://localhost:8081/toptracks";

      axios.post(topTracksURL, JSON.stringify(events))
        .then((response => {
          this.topTracks = response.data;
        }));
    },
    buildPlaylist: function (name, desc) {
      var buildPlaylistURL = "http://localhost:8081/buildplaylist?name=" +
        name +
        "&desc=" +
        desc;

        axios.post(buildPlaylistURL, JSON.stringify(this.topTracks))
          .then((response => {
            this.playlistStatus = response.status;
        }));
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
