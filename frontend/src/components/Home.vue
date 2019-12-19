<template>
  <div class="home">
    <div v-if="apiAddress">
      <AppHeader :token="token" :apiAddress="apiAddress">
      </AppHeader>
    </div>

    <div class="mainBody">
    <div class="row">
      <div class="leftsidebar col-md-4">
            {{selectedCities}}

        <Selector 
          selectorName="cities"
          maxItems="6"
          v-bind:selectedItems="selectedCities"
          v-bind:availableItems="availableCities">
        </Selector>

    {{selectedGenres}}
        <Selector 
          selectorName="genres"
          maxItems="10"
          v-bind:selectedItems="selectedGenres"
          v-bind:availableItems="availableGenres">
        </Selector>
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
          <button class="btn" :disabled="playlistLoading" v-on:click="buildPlaylist('Spooky Title', 'Spooooky Description!')">Build Playlist</button>
        </div>

        <div v-if="token">
          {{token}}
        </div>
        <div v-if="localEvents">
        {{localEvents}}
        </div>

        <div class="playlistResult" v-if="playlistLoading === true">
          <h4>Building playlist...</h4>
        </div>

        <div class="playlistResult" v-if="playlistStatus === 200">
          <h4>{{topTracks.length}} tracks created!</h4>
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
import AppHeader from './AppHeader.vue';
import PerformerCarousel from './PerformerCarousel.vue';
import Selector from './Selector.vue';
import createPersistedState from "vuex-persistedstate";

Vue.use(Vuex);

const state = {
  selectedCities: null,
  selectedGenres: null,
  availableCities: null,
  availableGenres: null,
  stateString: null,
  token: null
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
  },
  UPDATE_TOKEN(state, newValue) {
    state.token = newValue;
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
    AppHeader,
    PerformerCarousel,
    Selector
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
      isStateStringCorrect: null,
      playlistLoading: false,
      apiAddress: null
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
    },
    token: {
      get: function() {
        return store.state.token;
      },
      set: function(newValue) {
        store.commit("UPDATE_TOKEN", newValue);
      }
    }
  },
  watch: {
    'selectedCities': function() {
      if (this.selectedCities.length != 0 && this.selectedGenres.length != 0) {
        store.commit("UPDATE_SELECTED_CITIES", this.selectedCities);
        store.commit("UPDATE_AVAILABLE_CITIES", this.availableCities);

        this.getLocalEvents(this.selectedCities, this.selectedGenres);
      } else {
        this.localEvents = null;
        this.localPerformers = null;
        this.artistImages = null;
      }
    },
    'selectedGenres': function() {
      if (this.selectedCities.length != 0 && this.selectedGenres.length != 0) {
        store.commit("UPDATE_SELECTED_GENRES", this.selectedGenres);
        store.commit("UPDATE_AVAILABLE_GENRES", this.availableGenres);

        this.getLocalEvents(this.selectedCities, this.selectedGenres);
      } else {
        this.localEvents = null;
        this.localPerformers = null;
        this.artistImages = null;
      }
    }
  },
  mounted () {
    this.apiAddress = "https://otherside-api.herokuapp.com";
    //this.apiAddress = "http://localhost:8081";
    this.initializeStore();

    if (this.selectedCities.length != 0 && this.selectedGenres.length != 0) {
      this.getLocalEvents(this.selectedCities, this.selectedGenres);
    } 
    if (this.$route.query.state == null) {
      this.setNewSpotifyAuthenticationUrl();
    } 
    this.isStateStringCorrect = this.$route.query.state == this.stateString;
    this.token = this.$route.query.token;
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
    getAvailableCities: function() {
      this.availableCities = new Array();

      var citiesURL = `${this.apiAddress}/cities`;

      axios.get(citiesURL)
        .then((response => {
          this.availableCities = response.data;
        }));
    },
    getAvailableGenres: function() {
      this.availableGenres = new Array();

      var genresURL = `${this.apiAddress}/genres`;

      axios.get(genresURL)
        .then((response => {
          this.availableGenres = response.data;
        }));
    },
    getLocalEvents: function (cities, genres) {
      var cityString = this.arrayToQueryString(cities);
      var genreString = this.arrayToQueryString(genres);

      var localEventsURL = `${this.apiAddress}/localevents?cities=` +
      cityString +
      "&genres=" +
      genreString;

      axios.get(localEventsURL)
        .then((response => {
          this.localEvents = response.data;
          if (this.token != null) { 
            this.getArtistIDs(this.localEvents);
          }
          this.setPerformersArray(response.data);
        }));
    },
    setNewSpotifyAuthenticationUrl: function() {
      this.stateString = this.getRandomStateString()

      var getAuthenticationRequestUrl = `${this.apiAddress}/authenticate?state=`
      +  this.stateString;

      axios.get(getAuthenticationRequestUrl)
        .then(response => {
          this.spotifyAuthenticationUrl = response.data;
        })
    },
    setNewSpotifyToken: function() {
      var getSpotifyTokenRequestURl = `${this.apiAddress}/token?state=` + this.stateString;

      axios.get(getSpotifyTokenRequestURl)
        .then(response => {
          this.token = response.data;
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
      var artistIDsURL = `${this.apiAddress}/artistids`;

      const auth = {
        headers: {
          'Authorization': `Bearer ${this.token}`
        }
      }

      axios.post(artistIDsURL, JSON.stringify(events), auth)
        .then((response => {
          this.artists = response.data;
          this.setArtistImagesArray(this.artists);

          this.$nextTick(() => {
            this.$nextTick(() => {
              this.setCarouselStartSlides();
            });
          });
        })).catch(() => {
          this.forceLogOff();
        });
    },
    buildPlaylist: function (name, desc) {
      var topTracksURL = `${this.apiAddress}/toptracks`;
      var buildPlaylistURL = `${this.apiAddress}/buildplaylist?name=` +
        name +
        "&desc=" +
        desc;

      var artistIDs = new Array();

      this.artists.forEach(function(value) {
        artistIDs.push(value.Id);
      });

      this.playlistLoading = true;
      this.playlistStatus = null;

      const auth = {
        headers: {
          'Authorization': `Bearer ${this.token}`
        }
      }

      axios.post(topTracksURL, JSON.stringify(artistIDs), auth)
        .then((response => {
          this.topTracks = response.data;

          axios.post(buildPlaylistURL, JSON.stringify(response.data), auth)
            .then((response => {
              this.playlistLoading = false;
              this.playlistStatus = response.status;
            })).catch(() => {
              this.forceLogOff();
            })
          })).catch(() => {
            this.forceLogOff();
          }
        );
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
          if (artist.Name != "" && artist.ImageURL != "") {
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
    },
    forceLogOff: function() {
      this.localEvents = null;
      this.localPerformers = null;
      this.playlistStatus = null;
      this.spotifyAuthenticationUrl = null;
      this.topTracks = null;
      this.artists = null;
      this.artistImages = null;
      this.isStateStringCorrect = null;
      this.playlistLoading = false;
      this.token = null;
      this.stateString = null;

      window.location = window.location.href.split("?")[0];
    },
    obtainApiAddress: function() {
      var hostname = window.location.host;
      var hostnameWithoutPort = hostname.substring(0, hostname.indexOf(":"));
      return window.location.protocol + "//" + hostnameWithoutPort;
    }
  }
}
</script>


