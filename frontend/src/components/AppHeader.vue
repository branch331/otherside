<template>
    <b-navbar toggleable="lg" type="dark" variant="info">
      <b-navbar-brand href="#">Otherside</b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>

        <!-- Right aligned nav items -->

        <b-navbar-nav class="ml-auto">
          <b-nav-item-dropdown right>
            <!-- Using 'button-content' slot -->
            <template v-slot:button-content>
              <span v-if="currentUser"> 
                {{currentUser}}
              </span>
              <span v-else>
                My Account
              </span>
            </template>
            <b-dropdown-item href="#">Profile</b-dropdown-item>
            <b-dropdown-item href="#">Sign Out</b-dropdown-item>
          </b-nav-item-dropdown>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
</template>

<script>
import axios from 'axios';

export default {
  name: 'AppHeader',
  props: ['token', 'apiAddress'],
  data () {
    return {
      currentUser: null
    }
  },
  mounted () {
    if (this.token && this.apiAddress != null && this.currentUser == null) {
      this.getCurrentUser();
    }
  },
  methods: {
    getCurrentUser: function() {
      const auth = {
        headers: {
          'Authorization': `Bearer ${this.token}`
        }
      }

      var currentUserURL = `${this.apiAddress}/user`;
      
      axios.get(currentUserURL, auth)
        .then((response => {
          this.currentUser = response.data;
        }));
    }
  }
}
</script>
