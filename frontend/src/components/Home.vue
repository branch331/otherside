<template>
  <div class="home">
    <div class="mainBody">

    <div class="row">
      <div class="leftsidebar col-md-4">
        <div class="selections">
          <h5>1. Pick some cities</h5>
          <ul v-if="selectedCities">
            <li v-for="city in selectedCities" v-bind:key="city" v-on:click="transferArrayValue(selectedCities, availableCities, city)">
              {{city}} <span class="selectedX">X</span>
            </li>
          </ul>

          <h5>Available Cities</h5>
          <ul v-if="availableCities">
            <li v-for="city in availableCities" v-bind:key="city" v-on:click="transferArrayValue(availableCities, selectedCities, city)">
              {{city}} <span class="selectedX">X</span>
            </li>
          </ul>
        </div>

        <div class="selections">
          <h5>2. Pick some genres</h5>
          <ul v-if="selectedGenres">
            <li v-for="genre in selectedGenres" v-bind:key="genre" v-on:click="transferArrayValue(selectedGenres, availableGenres, genre)">
              {{genre}} <span class="selectedX">X</span>
            </li>
          </ul>

          <h5>Available Genres</h5>
          <ul v-if="availableGenres">
            <li v-for="genre in availableGenres" v-bind:key="genre" v-on:click="transferArrayValue(availableGenres, selectedGenres, genre)">
              {{genre}} <span class="selectedX">X</span>
            </li>
          </ul>
        </div>
      </div>

      <div class="playlistcontent col-md-6">

        <div class="eventsHolder">
          <h2>Events Data:</h2>

          <div v-if="artistImages"> 
            <PerformerCarousel 
              carouselId="performerCarouselLeft"
              ref="carouselLeft"
              v-bind:artistImages="artistImages">
            </PerformerCarousel>

            <PerformerCarousel 
              carouselId="performerCarouselRight"
              ref="carouselRight"
              v-bind:artistImages="artistImages">
            </PerformerCarousel>

            <PerformerCarousel 
              carouselId="performerCarouselMain"
              ref="carouselMain"
              v-bind:artistImages="artistImages">
            </PerformerCarousel>
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

        <div v-if="artistImages">
          {{artistImages}}
          </div>

      </div>
    </div>

    </div>
  </div>

</template>

<script>
import axios from 'axios';
import Vue from "vue";
import Vuex from "vuex";
import PerformerCarousel from './PerformerCarousel.vue';
import createPersistedState from "vuex-persistedstate";

Vue.use(Vuex);

const state = {
  selectedCities: null,
  selectedGenres: null,
  availableCities: null,
  availableGenres: null,
  stateString: null
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
  components: {
    PerformerCarousel
  },
  data () {
    return {
      localEvents: null,
      localPerformers: null,
      playlistStatus: null,
      spotifyAuthenticationUrl: null,
      topTracks: null,
      artists: null,
      artistImages: null,
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
    this.initializeStore();
    
    if (this.selectedCities.length != 0 && this.selectedGenres.length != 0) {
      this.getLocalEvents(this.selectedCities, this.selectedGenres);
    } 
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
        this.localPerformers = null;
        this.artistImages = null;
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
          this.setPerformersArray(response.data);
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
          this.artists = response.data;
          this.setArtistImagesArray(this.artists);

          this.$nextTick(() => {
            this.$nextTick(() => {
              this.setCarouselStartSlides();
            });
          });
        }));
    },
    buildPlaylist: function (name, desc) {
      var topTracksURL = "http://localhost:8081/toptracks";
      var buildPlaylistURL = "http://localhost:8081/buildplaylist?name=" +
        name +
        "&desc=" +
        desc;

      var artistIDs = new Array();

      this.artists.forEach(function(value) {
        artistIDs.push(value.Id);
      });

      axios.post(topTracksURL, JSON.stringify(artistIDs))
        .then((response => {
          this.topTracks = response.data;
          axios.post(buildPlaylistURL, JSON.stringify(response.data))
            .then((response => {
              this.playlistStatus = response.status;
          }));
        }));
    },
    setPerformersArray: function(localEvents) {
      var tempPerformers = new Array();

      if (localEvents != null) {
        localEvents.forEach(function(event) {
          if (event.Performers != null) {
              event.Performers.forEach(function(performer) {
                tempPerformers.push(performer);
              });
          }
        }); 

        this.localPerformers = tempPerformers;
      }
    },
    setArtistImagesArray: function(artists) {
      var tempArtistImages = new Array();

      if (artists != null) {
        artists.forEach(function(artist) {
          if (artist.Name != null) {
            tempArtistImages.push({"Name": artist.Name, "ImageURL": artist.ImageURL});
          }
        }); 

        this.artistImages = tempArtistImages;
      }
    },
    setCarouselStartSlides: function() {
      var numberOfPerformers = -1; 
      
      if (this.artistImages != null) {
        numberOfPerformers = this.artistImages.length;
      }

      if (numberOfPerformers > 0) {
        this.$refs.carouselLeft.pauseWrapper();
        this.$refs.carouselRight.pauseWrapper();
        this.$refs.carouselMain.pauseWrapper();

        if (numberOfPerformers >= 3) {
          this.$refs.carouselLeft.setSlideWrapper(0);
          this.$refs.carouselMain.setSlideWrapper(1);
          this.$refs.carouselRight.setSlideWrapper(2);
        } else if (numberOfPerformers == 2) {
          this.$refs.carouselLeft.setSlideWrapper(0);
          this.$refs.carouselMain.setSlideWrapper(1);
          this.$refs.carouselRight.setSlideWrapper(0);
        } else {
          this.$refs.carouselLeft.setSlideWrapper(0);
          this.$refs.carouselMain.setSlideWrapper(0);
          this.$refs.carouselRight.setSlideWrapper(0);
        }

        this.$refs.carouselLeft.startWrapper();
        this.$refs.carouselRight.startWrapper();
        this.$refs.carouselMain.startWrapper();
      }
    }
  }
}
</script>


