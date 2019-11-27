<template>
  <div class="home">
    <div class="mainBody">
    <b-navbar toggleable="lg" type="dark" variant="info">
      <b-navbar-brand href="#">Otherside</b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>

        <!-- Right aligned nav items -->

        <b-navbar-nav class="ml-auto">
          <b-nav-item href="#">Link</b-nav-item>
          <b-nav-item href="#">Disabled</b-nav-item>

          <b-nav-item-dropdown right>
            <!-- Using 'button-content' slot -->
            <template v-slot:button-content>
              My Account
            </template>
            <b-dropdown-item href="#">Profile</b-dropdown-item>
            <b-dropdown-item href="#">Sign Out</b-dropdown-item>
          </b-nav-item-dropdown>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>

    <div class="row">
      <div class="leftsidebar col-md-4">
        <div class="selections">
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
        </div>

        <div class="selections">
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
        </div>
      </div>

      <div class="playlistcontent col-md-6">

        <div class="eventsHolder">
          <h2>Events Data:</h2>

          <div v-if="localEvents">
            <b-carousel
              id="performerCarousel"
              v-model="slide"
              :interval="1000"
              fade=true
              no-hover-pause=true 
              no-touch=true 
              background="#ababab"
              img-width="1024"
              img-height="480"
              style="text-shadow: 1px 1px 2px #333;">

              <template v-for="event in localEvents">
                  <b-carousel-slide class="performerSlide" v-for="(performer, i) in event.Performers" v-bind:key="`${i}-${performer}`" img-src="https://picsum.photos/1024/1024/?image=54">
                    <p>{{performer}}</p>
                  </b-carousel-slide>
              </template>
            </b-carousel>
          </div>
        </div>

        <div class="btnHolder" v-if="!isStateStringCorrect">
          <button class="btn" v-on:click="redirectToURL">Log In</button>
        </div>
        
        <div class="btnHolder" v-if="isStateStringCorrect">
          <button class="btn" v-on:click="buildPlaylist('Spooky Title', 'Spooooky Description!')">Build Playlist</button>
        </div>

        <div v-if="topTracks">
          {{topTracks.length}}
        </div>

      </div>
    </div>

    </div>
    <div class="footer">
      Powered by <a href="https://seatgeek.com/"><img src="../assets/seatgeek-logo.png"></a>
    </div>
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
  availableGenres: null,
  stateString: null,
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
  UPDATE_STATE_STRING(state, newValue) {
    state.stateString = newValue;
  }
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
      topTracks: null,
      artistIDs: null,
      isStateStringCorrect: null
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
    },
    stateString: {
      get: function() {
        return store.state.stateString;
      },
      set: function(newValue) {
        store.commit("UPDATE_STATE_STRING", newValue);
      }
    }
  },
  mounted () {
    if (this.selectedCities != null && this.selectedGenres != null) {
      this.getLocalEvents(this.selectedCities, this.selectedGenres);
    }
    this.initializeStore();
    if (this.$route.query.state == null) {
      this.setNewSpotifyAuthenticationUrl();
    } 
    this.isStateStringCorrect = this.$route.query.state == this.stateString;
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
          if (this.$route.query.state == this.stateString) { //User has logged in successfully
            this.getArtistIDs(this.localEvents);
          }
        }));
    },
    setNewSpotifyAuthenticationUrl: function() {
      this.stateString = this.getRandomStateString()

      var getAuthenticationRequestUrl = "http://localhost:8081/authenticate?state="
      +  this.stateString;

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


