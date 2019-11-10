<template>
  <div class="callback">

    <h2>CALLBACK</h2>

    <h2>Events Data:</h2>
    <p>THIS component's spotify state: {{ spotifyStateString }}</p>
    <div v-if="localEvents">
      <button v-on:click="buildPlaylist('MY special playlist!', 'My special playlists description')">Build Playlist</button>

      <p>{{localEvents}}</p>

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
      spotifyStateString: null
    }
  },
  mounted () {
    this.getLocalEvents('[78759,22032,70001]', '[rock,electronic,hip-hop]');
  },
  methods: {
    getLocalEvents: function (postalCodes, genres) {
      var localEventsURL = "http://localhost:8081/localevents?postcodes=" +
      postalCodes +
      "&genres=" +
      genres;

      axios.get(localEventsURL)
        .then((response => {
          this.localEvents = response.data;
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
