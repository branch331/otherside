<template>
  <div class="home">

    <h2>HOME</h2>

      <button v-on:click="redirectToURL">Log In</button>

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

    <div v-if="artistIDs">
      <button v-on:click="buildPlaylist('Spooky Title', 'Spooooky Description!')">Build Playlist</button>
    </div>

    <h2>Playlist Status:</h2>
    <h2 v-if="playlistStatus">{{playlistStatus}}</h2>

  </div>
</template>

<script>
import axios from 'axios';
import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";

Vue.use(Vuex);

const state = {
  selectedCities: null,
  selectedGenres: null,
  availableCities: null,
  availableGenres: null
};

const mutations = {
  UPDATE_SELECTED_CITIES(state, payload) {
    state.selectedCities = payload;
  },
  UPDATE_SELECTED_GENRES(state, payload) {
    state.selectedGenres = payload;
  },
  UPDATE_AVAILABLE_CITIES(state, payload) {
    state.availableCities = payload;
  },
  UPDATE_AVAILABLE_GENRES(state, payload) {
    state.availableGenres = payload;
  },
};

const store = new Vuex.Store({
  state,
  mutations,
  plugins: [createPersistedState()]
});

export default {
  name: 'home',
  store,
  data () {
    return {
      localEvents: null,
      playlistStatus: null,
      spotifyAuthenticationUrl: null,
      spotifyStateString: null,
      topTracks: null,
      artistIDs: null
    }
  },
  computed: {
    selectedCities: {
      get: function() {
        return store.state.selectedCities;
      },
      set: function(newValue) {
        store.commit("UPDATE_SELECTED_CITIES", newValue);
      }
    },
    availableCities: {
      get: function() {
        return store.state.availableCities;
      },
      set: function(newValue) {
        store.commit("UPDATE_AVAILABLE_CITIES", newValue);
      }
    },
    selectedGenres: {
      get: function() {
        return store.state.selectedGenres;
      },
      set: function(newValue) {
        store.commit("UPDATE_SELECTED_GENRES", newValue);
      }
    },
    availableGenres: {
      get: function() {
        return store.state.availableGenres; 
      },
      set: function(newValue) {
        store.commit("UPDATE_AVAILABLE_GENRES", newValue);
      }
    }
  },
  mounted () {
    this.setNewSpotifyAuthenticationUrl();
    if (this.selectedCities != null && this.selectedGenres != null) {
      this.getLocalEvents(this.selectedCities, this.selectedGenres);
    }
    this.initializeStore();
  },
  methods: {
    initializeStore: function() {
      if (this.selectedCities == null) {
        this.selectedCities = new Array();
      }
      if (this.selectedGenres == null) {
        this.selectedGenres = new Array();
      }
      if (this.availableCities == null) {
        this.getAvailableCities();
      }
      if (this.availableGenres == null) {
        this.getAvailableGenres();
      }
    },
    transferArrayValue: function(sourceArray, targetArray, value) {
      var index = sourceArray.indexOf(value);
      if (index > -1) {

        if (!(targetArray == this.selectedCities && targetArray.length == 6) && 
        !(targetArray == this.selectedGenres && targetArray.length == 10)) {
          if (sourceArray == this.selectedCities) {
            sourceArray.splice(index, 1);
            targetArray.push(value);
            this.selectedCities = sourceArray;
            this.availableCities = targetArray;
          } 
          else if (sourceArray == this.availableCities) {
            sourceArray.splice(index, 1);
            targetArray.push(value);
            this.availableCities = sourceArray;
            this.selectedCities = targetArray;
          }
          else if (sourceArray == this.selectedGenres) {
            sourceArray.splice(index, 1);
            targetArray.push(value);
            this.selectedGenres = sourceArray;
            this.availableGenres = targetArray;
          }
          else if (sourceArray == this.availableGenres) {
            sourceArray.splice(index, 1);
            targetArray.push(value);
            this.availableGenres = sourceArray;
            this.selectedGenres = targetArray;
          }
        }
      }
      if (this.selectedCities.length != 0 && this.selectedGenres.length != 0) {
        this.getLocalEvents(this.selectedCities, this.selectedGenres);
      } else {
        this.localEvents = null;
      }
    },
    getAvailableCities: function() {
      this.availableCities = new Array();

      var citiesURL = "http://localhost:8081/cities";

      axios.get(citiesURL)
        .then((response => {
          this.availableCities = response.data;
        }));
    },
    getAvailableGenres: function() {
      this.availableGenres = new Array();

      var genresURL = "http://localhost:8081/genres";

      axios.get(genresURL)
        .then((response => {
          this.availableGenres = response.data;
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
          this.getArtistIDs(this.localEvents);
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
    },
    arrayToQueryString: function (array) {
      var queryString = '[';

      for (var i = 0; i < array.length; i++) {
        queryString = queryString.concat(array[i] + ',');
      }

      queryString = queryString.slice(0, -1).concat(']');
      return queryString;
    },
    getArtistIDs: function(events) {
      var artistIDsURL = "http://localhost:8081/artistids";

      axios.post(artistIDsURL, JSON.stringify(events))
        .then((response => {
          this.artistIDs = response.data;
        }));
    },
    buildPlaylist: function (name, desc) {
      var topTracksURL = "http://localhost:8081/toptracks";
      var buildPlaylistURL = "http://localhost:8081/buildplaylist?name=" +
        name +
        "&desc=" +
        desc;

      axios.post(topTracksURL, JSON.stringify(this.artistIDs))
        .then((response => {
          this.topTracks = response.data;
          axios.post(buildPlaylistURL, JSON.stringify(response.data))
            .then((response => {
              this.playlistStatus = response.status;
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
