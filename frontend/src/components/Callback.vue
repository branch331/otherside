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


    <button v-on:click="getLocalEvents(selectedCities, selectedGenres)">Build Playlist</button>

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
      availableGenres: null
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
        sourceArray.splice(index, 1);
        targetArray.push(value);
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
      var cityString = '[';
      for (var i = 0; i < cities.length; i++) {
        cityString = cityString.concat(cities[i] + ',');
      }
      cityString = cityString.slice(0, -1).concat(']');

      var genreString = '[';
      for (var j = 0; j < genres.length; j++) {
        genreString = genreString.concat(genres[j] + ',');
      }
      genreString = genreString.slice(0, -1).concat(']');

      var localEventsURL = "http://localhost:8081/localevents?cities=" +
      cityString +
      "&genres=" +
      genreString;

      axios.get(localEventsURL)
        .then((response => {
          this.localEvents = response.data;
          this.buildPlaylist('MY special playlist!', 'My special playlists description');
        }));
    },
    buildPlaylist: function (name, desc) {
      var buildPlaylistURL = "http://localhost:8081/buildplaylist?name=" +
        name +
        "&desc=" +
        desc;

      axios.post("http://localhost:8081/toptracks", JSON.stringify(this.localEvents))
        .then((response => {
          axios.post(buildPlaylistURL, JSON.stringify(response.data))
            .then((response => {
              this.playlistStatus = response.data;
          }));
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
